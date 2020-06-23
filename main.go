package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/parsekafka/cli"
	"io"
	"os"
)

func main() {
	logFilename := "/tmp/parsekafka.log"
	logFile, _ := os.OpenFile(logFilename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	defer logFile.Close()

	writers := []io.Writer{
		logFile,
		os.Stdout,
	}
	fileAndStdoutWriter := io.MultiWriter(writers...)

	log.SetOutput(fileAndStdoutWriter)
	log.SetLevel(log.DebugLevel)
	log.Infoln("main.main():Start ParseKafka Main")
	cli.Run()
}
