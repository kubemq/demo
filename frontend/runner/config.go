package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	ApiAddress string
	Interval   int
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	viper.AddConfigPath("./")
	viper.SetConfigName(".config")
	viper.BindEnv("ApiAddress", "API_ADDRESS")
	viper.BindEnv("Interval", "INTERVAL")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
