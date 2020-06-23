package config

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type KafkaConfig struct {
	Topic    string   `topic`
	Ips      []string `ips`
	Database string   `database`
	MaxIdle  int      `maxidle`
	MaxOpen  int      `maxopen`
}

func UnmarshalConfig(path string) (*KafkaConfig, error) {
	//import.yaml  path
	in, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := new(KafkaConfig)        //
	err = yaml.Unmarshal(in, &config) //
	if err != nil {
		return nil, err
	}
	log.Infoln("SUCCESS MARSHAL CONFIG")
	return config, nil
}
