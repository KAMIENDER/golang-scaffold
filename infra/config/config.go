package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const ENV_CONFIG_PATH = "config_path"

type Config struct {
	RedisConf RedisConfig `yaml:"redis_config"`
	DBConf    DBConfig    `yaml:"db_config"`
	PayConf   PayConfig   `yaml:"pay_config"`
}

type PayConfig struct {
	AppID         string `yaml:"appid"`
	AppPrivateKey string `yaml:"app_private_key"`
	AppPublicKey  string `yaml:"app_public_key"`
	AliPublicKey  string `yaml:"ali_public_key"`
	IsOnline      bool   `yaml:"is_online"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type DBConfig struct {
	Addr     string `yaml:"addr"`
	Port     string `yaml:"port"`
	UserName string `yaml:"user_name"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
}

func NewConfig() (*Config, error) {
	// read yaml file
	filePath := os.Getenv(ENV_CONFIG_PATH)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// file content to struct
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
