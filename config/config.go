package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ApiHash  string `yaml:"app_hash"`
	ApiID    string `yaml:"api_id"`
	BotToken string `yaml:"bot_token"`
	DB       string `yaml:"db"`
}

func FromFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("can't open config file: %s", err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("can't read config file: %s", err)
	}
	cfg := &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal config file: %s", err)
	}
	return cfg, nil
}
