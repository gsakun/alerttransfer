package main

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gsakun/alerttransfer/cli"
	colorable "github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

func init() {
	loglevel := os.Getenv("LOG_LEVEL")
	var logLevel log.Level
	log.Infof("loglevel env is %s", loglevel)
	if loglevel == "debug" {
		log.SetLevel(log.DebugLevel)
		logLevel = log.DebugLevel
		log.Infof("log level is %s", loglevel)
		log.SetReportCaller(true)
	} else {
		log.SetLevel(log.InfoLevel)
		logLevel = log.InfoLevel
		log.Infoln("log level is normal")
	}
	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   "../logs/alerttransfer.log",
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
		Level:      logLevel,
		Formatter: &log.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		},
	})
	log.SetOutput(colorable.NewColorableStdout())
	if err != nil {
		log.Fatalf("Failed to initialize file rotate hook: %v", err)
	}
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})
	log.SetReportCaller(true)
	log.AddHook(rotateFileHook)
}
func main() {
	cli.Run()
}
