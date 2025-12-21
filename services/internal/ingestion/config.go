package ingestion

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const ConfigPrefix = "INGESTION"

type Config struct {
	BindAddress string `envconfig:"BIND_ADDRESS" required:"true" default:"0.0.0.0:8080"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	if err := envconfig.Process(ConfigPrefix, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
