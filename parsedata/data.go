package parsedata

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/gsakun/alerttransfer/datatype"
	"github.com/gsakun/alerttransfer/db"
	logf "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

var (
	wg     sync.WaitGroup
	logger = log.New(os.Stderr, "[srama]", log.LstdFlags)
)

func Parsedata(topic string, ips []string) {
	sarama.Logger = logger

	consumer, err := sarama.NewConsumer(ips, nil)
	if err != nil {
		logger.Println("Failed to start consumer: %s", err)
	}

	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		logger.Println("Failed to get the list of partitions: ", err)
	}

	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			logger.Printf("Failed to start consumer for partition %d: %s\n", partition, err)
		}
		defer pc.AsyncClose()

		wg.Add(1)

		go func(sarama.PartitionConsumer) {
			defer wg.Done()
			for msg := range pc.Messages() {
				logf.Printf("Partition:%d, Offset:%d, Key:%s, Value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				fmt.Println()
				data := HandleMessage(string(msg.Value))
				if data == nil {
					log.Println("nil")
				} else {
					db.Dataprocess(data)
				}
			}
		}(pc)
	}

	wg.Wait()

	logger.Println("Done consuming topic packetbeat")
	consumer.Close()
}

func HandleMessage(message string) (n *datatype.NetData) {
	n = &datatype.NetData{}
	protype := gjson.Get(message, "type").String()
	logf.Println(protype)

	//ndirection := gjson.Get(message, "direction").String()
	if protype == "flow" {
		stat := gjson.Get(message, "final").Bool()
		if stat {
			/*if direction == "" {
				logf.Warningln("this message don not need parse")
				return nil
			}*/
			log.Println("Start get serverjson data")
			n.ServerIp = gjson.Get(message, "dest.ip").String()
			n.ClientIp = gjson.Get(message, "source.ip").String()
			log.Println("FLOW SERVER IP")
			log.Println(n.ServerIp)
			log.Println("FLOW CLIENT IP")
			log.Println(n.ClientIp)
			if n.ClientIp == n.ServerIp || n.ClientIp == "127.0.0.1" || n.ServerIp == "127.0.0.1" || n.ClientIp == "" || n.ServerIp == "" || strings.Contains(n.ServerIp, "::") || strings.Contains(n.ClientIp, "::") || n.ServerPort == 0 || n.ClientPort == 0 {
				return nil
			}
			n.ServerPort = gjson.Get(message, "dest.port").Int()
			n.ClientPort = gjson.Get(message, "source.port").Int()
			n.SourceBytes = gjson.Get(message, "source.stats.net_bytes_total").Int()
			n.DestBytes = gjson.Get(message, "dest.stats.net_bytes_total").Int()
			n.Time = HandleDateTime(gjson.Get(message, "@timestamp").String())
			n.ProtoType = "flow"
			n.Stat = true
		} else {
			log.Println("FINAL FALSE")
			return nil
		}
	} else {
		log.Println("Start get serverjson data")
		n.ServerIp = gjson.Get(message, "ip").String()
		n.ServerPort = gjson.Get(message, "port").Int()
		n.ClientIp = gjson.Get(message, "client_ip").String()
		n.ClientPort = gjson.Get(message, "client_port").Int()
		if n.ClientIp == n.ServerIp || n.ClientIp == "127.0.0.1" || n.ServerIp == "127.0.0.1" || n.ClientIp == "" || n.ServerIp == "" || strings.Contains(n.ServerIp, "::") || strings.Contains(n.ClientIp, "::") || n.ServerPort == 0 || n.ClientPort == 0 {
			return nil
		}
		n.Time = HandleDateTime(gjson.Get(message, "@timestamp").String())
		n.ProtoType = protype
		stat := gjson.Get(message, "status").String()
		n.Stat = false
		if stat == "OK" {
			n.Stat = true
		}
		if protype == "in" {
			n.SourceBytes = gjson.Get(message, "bytes_in").Int()
			n.DestBytes = gjson.Get(message, "bytes_out").Int()
		} else {
			n.SourceBytes = gjson.Get(message, "bytes_out").Int()
			n.DestBytes = gjson.Get(message, "bytes_in").Int()
		}
	}
	logf.Println(n)
	return n
}
func HandleDateTime(date string) string {
	f1 := strings.Split(date, "T")
	f2 := strings.Split(f1[1], ".")[0]
	t := f1[0] + " " + f2
	return t
}
