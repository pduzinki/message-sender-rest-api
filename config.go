package main

import (
	"encoding/json"
	"os"
)

// Config holds variables neccessary to properly configure the app
type Config struct {
	// Keyspace string
	Mailgun MailgunConfig `json:"mailgun"`
}

// MailgunConfig holds variables neccessary to enable Mailgun client
type MailgunConfig struct {
	APIKey string `json:"api_key"`
	Domain string `json:"domain"`
}

// LoadConfig loads config values from .config file,
// if not found, then loads default config values
func LoadConfig() Config {
	file, err := os.Open(".config")
	if err != nil {
		return LoadDefaultConfig()
	}

	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}

	return config
}

// LoadDefaultConfig load default config values
func LoadDefaultConfig() Config {
	return Config{
		Mailgun: MailgunConfig{
			APIKey: "fake-api-key",
			Domain: "fake-domain.com",
		},
	}
}
