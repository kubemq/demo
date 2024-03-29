package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	KubeMQHost   string
	KubeMQPort   int
	Channel      string
	Group        string
	SlackToken   string
	SlackChannel string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	viper.AddConfigPath("./")
	viper.SetConfigName(".config")
	viper.BindEnv("KubeMQHost", "KUBEMQ_HOST")
	viper.BindEnv("KubeMQPort", "KUBEMQ_POST")
	viper.BindEnv("Channel", "CHANNEL")
	viper.BindEnv("Group", "GROUP")
	viper.BindEnv("SlackToken", "SLACK_TOKEN")
	viper.BindEnv("SlackChannel", "SLACK_CHANNEL")

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
