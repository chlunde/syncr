package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type SyncrConfig struct {
	Source, Destination string
	Hosts               []string
}

func LoadConfig() ([]SyncrConfig, error) {
	var config []SyncrConfig
	data, err := ioutil.ReadFile("syncr.yaml")
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal([]byte(data), &config); err != nil {
		return nil, err
	}
	return config, nil
}
