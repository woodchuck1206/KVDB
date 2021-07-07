package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	BASIC_CONFIG_FILE_ADDR string = "./src/config/config.yaml"
)

type Config struct {
	EConfig EngineConfig `yaml:"engine"`
	SConfig ServerConfig `yaml:"server"`
}

type EngineConfig struct {
	RValue       int `yaml:"rValue"`
	MemtableSize int `yaml:"memtableSize"`
	SstableSize  int `yaml:"sstableBaseSize"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

func ParseConfig(configFileAddr string) Config {
	if len(configFileAddr) == 0 {
		configFileAddr = BASIC_CONFIG_FILE_ADDR
	}
	bytesRead, err := ioutil.ReadFile(configFileAddr)
	if err != nil {
		fmt.Printf("READ CONFIG FILE ERROR\n%v\n", err)
		os.Exit(1) // no config, no nothing
	}
	config := Config{}
	err = yaml.Unmarshal(bytesRead, &config)
	if err != nil {
		fmt.Printf("READ CONFIG FILE ERROR\n%v\n", err)
		os.Exit(1)
	}
	return config
}
