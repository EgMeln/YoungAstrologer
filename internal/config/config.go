package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config represents config's variables
type Config struct {
	NASAAPIKey  string
	PostgresURL string
	ServerPort  string
}

// NewConfig loads and parses config file from given paths
func NewConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./internal/config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	apiKey := viper.GetString("NASA_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("please set the NASA_API_KEY in the config file")
	}

	dbURL := viper.GetString("POSTGRES_URL")
	if apiKey == "" {
		return nil, fmt.Errorf("please set the POSTGRES_URL in the config file")
	}

	port := viper.GetString("SERVER_PORT")
	if apiKey == "" {
		return nil, fmt.Errorf("please set the SERVER_PORT in the config file")
	}

	return &Config{
		NASAAPIKey:  apiKey,
		PostgresURL: dbURL,
		ServerPort:  port,
	}, nil
}
