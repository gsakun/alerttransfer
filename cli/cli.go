package cli

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gsakun/alerttransfer/config"
	"github.com/gsakun/alerttransfer/db"
	"github.com/gsakun/alerttransfer/parsedata"
	log "github.com/sirupsen/logrus"
)

var kafkainfo config.KafkaConfig

func Run() {
	if len(os.Args) == 1 {
		help()
		return
	}

	var err error
	command := os.Args[1]
	log.Debugf("cli.Run(): cli args:%+v\n", os.Args)
	if command == "-c" {
		if len(os.Args) != 3 {
			importErr := "The command is wrong. See help"
			log.Errorln(importErr)
			return
		}
		path := os.Args[2]
		kafkainfo, err := config.UnmarshalConfig(path)
		if err != nil {
			log.Errorln(err)
		}
		db.Init(kafkainfo.Database, kafkainfo.MaxIdle, kafkainfo.MaxOpen)
		parsedata.Parsedata(kafkainfo.Topic, kafkainfo.Ips)
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigs
			fmt.Println("Bye")
			db.DB.Close()
			os.Exit(0)
		}()
	}
	if command == "help" {
		help()
	}
	if err != nil {
		log.Errorf("cli.Run():%+v\n", err)
		return
	}
}
func help() {
	var helpString = `Usage: alerttransfer -c config.yaml`
	log.Println(helpString)
}
