package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Port    string `json:"port"`
	AirqKey string `json:"airq_key"`
	FBKey   string `json:"fb_key"`

	DB struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Database string `json:"db"`
	} `json:"db"`
}

var (
	config Config
)

func loadConfig(fileName string) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}

	if len(config.Port) < 2 {
		panic("port in config.json should be like :8080")
	}

	if config.DB.User == "" {
		panic("No DB Username")
	}
}
