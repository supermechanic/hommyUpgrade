package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

//Config 配置实例
var Config config

type service struct {
	Name     string `yaml:"name"`
	Username string `yaml:"user"`
	Password string `yaml:"pass"`
	Address  string `yaml:"addr"`
	Port     string `yaml:"port"`
	Version  string `yaml:"ver"`
	Filepath string `yaml:"path"`
}

type mysql struct {
	Base         service `yaml:"base"`
	DatabaseName string  `yaml:"database"`
}

type mredis struct {
	Base        service `yaml:"base"`
	MaxIdle     int     `yaml:"max-idle"`
	MaxActive   int     `yaml:"max-active"`
	IdleTimeout int     `yaml:"idle-timeout"`
	Timeout     int     `yaml:"timeout"`
}

type config struct {
	Mysql mysql  `yaml:"mysql"`
	Redis mredis `yaml:"redis"`
}

func init() {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	log.Println(Config)
}
