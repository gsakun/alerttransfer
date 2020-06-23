package main

import (
	"io"
	"os"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gsakun/alerttransfer/cli"
	log "github.com/sirupsen/logrus"
)

func main() {
	logFilename := "/tmp/alerttransfer.log"
	logFile, _ := os.OpenFile(logFilename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	defer logFile.Close()

	writers := []io.Writer{
		logFile,
		os.Stdout,
	}
	fileAndStdoutWriter := io.MultiWriter(writers...)

	log.SetOutput(fileAndStdoutWriter)
	log.SetLevel(log.DebugLevel)
	log.Infoln("main.main():Start alerttransfer Main")
	cli.Run()
}
