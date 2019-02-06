package main

import "github.com/kelseyhightower/envconfig"

type _Config struct {
	Debug                 bool
	MaxWorkers            int    `envconfig:"MAX_WORKERS" default:"32"`
	ClickHouseReadTimeout int    `envconfig:"CLICKHOUSE_READ_TIMEOUT" default:"10"`
	ShodanAPIKey          string `envconfig:"SHODAN_API_KEY" required:"true"`
}

func _ConfigFromEnv() (*_Config, error) {
	cfg := &_Config{}
	err := envconfig.Process("", cfg)
	return cfg, err
}
