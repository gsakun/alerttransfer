package parsedata

import (
	"regexp"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/gsakun/alerttransfer/datatype"
	"github.com/gsakun/alerttransfer/db"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

var (
	wg sync.WaitGroup
)

func Parsedata(topic string, ips []string) {
	consumer, err := sarama.NewConsumer(ips, nil)
	if err != nil {
		log.Errorf("Failed to start consumer: %s", err)
	}

	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		log.Errorf("Failed to get the list of partitions: ", err)
	}

	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			log.Errorf("Failed to start consumer for partition %d: %s\n", partition, err)
		}
		defer pc.AsyncClose()

		wg.Add(1)

		go func(sarama.PartitionConsumer) {
			defer wg.Done()
			for msg := range pc.Messages() {
				log.Debugf("Partition:%d, Offset:%d, Key:%s, Value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				handler(string(msg.Value))
			}
		}(pc)
	}

	wg.Wait()

	log.Infoln("Done consuming topic packetbeat")
	consumer.Close()
}

func handler(message string) {
	var alert datatype.AlarmData
	logfield := gjson.Get(message, "log").String()
	//re := regexp.MustCompile(`(^\d{2}:\d{2}:\d{2}.\d{9}): ([a-zA-Z]+) cluster=(\w+) ; content=(.*) ; condition=(.*) ; desc=(.*) ; tags=(.*) ; (.*)`)
	re := regexp.MustCompile(`(^\d{2}:\d{2}:\d{2}.\d{9}): ([a-zA-Z]+) cluster=(\w+) content=(.*) condition=(.*) desc=(.*) tags=(.*) pod_name=(.*) (.*)`)
	output := gjson.Get(logfield, "output").String()
	if !re.MatchString(output) {
		log.Errorf("Regexp failed info %s", output)
		return
		//return fmt.Errorf("Regexp failed info %s", output)
	}
	params := re.FindStringSubmatch(output)
	if len(params) != 10 {
		log.Errorf("Param num incorrectness")
		return
		//return fmt.Errorf("Regexp failed info %s", output)
	}
	outputfields := gjson.Get(logfield, "output_fields").Map()
	alert.PodName = outputfields["k8s.pod.name"].String()
	alert.PodNamespace = outputfields["k8s.ns.name"].String()
	alert.PodAlarmRule = gjson.Get(logfield, "rule").String()
	cstSh, _ := time.LoadLocation("Asia/Shanghai")
	alert.PodAlarmTime = gjson.Get(logfield, "time").Time().Local().In(cstSh).Format("2006-01-02 15:04:05")
	alert.PodAlarmPrority = params[2]
	alert.PodCluster = params[3]
	alert.PodAlarmContent = params[4]
	alert.PodAlarmCondition = params[5]
	alert.PodAlarmDesc = params[6]
	err := alert.Handler(db.DB)

	if err != nil {
		log.Errorf("Handler failed errinfo %v", err)
		//return err
	}
	return
}
