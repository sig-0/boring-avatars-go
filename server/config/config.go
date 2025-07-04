package config

import (
	"errors"
	"os"
	"regexp"

	"github.com/pelletier/go-toml"
)

const DefaultListenAddress = "0.0.0.0:8545"

var ErrInvalidListenAddress = errors.New("invalid listen address")

var listenAddressRegex = regexp.MustCompile(`^\d{1,3}(\.\d{1,3}){3}:\d+$`)

// Config defines the base-level server configuration
type Config struct {
	// The associated CORS config, if any
	CORSConfig *CORS `toml:"cors_config"`

	// The address at which the server will be served.
	// Format should be: <IP>:<PORT>
	ListenAddress string `toml:"listen_address"`
}

// DefaultConfig returns the default server configuration
func DefaultConfig() *Config {
	return &Config{
		ListenAddress: DefaultListenAddress,
		CORSConfig:    DefaultCORSConfig(),
	}
}

// ValidateConfig validates the server configuration
func ValidateConfig(config *Config) error {
	// Validate the listen address
	if !listenAddressRegex.MatchString(config.ListenAddress) {
		return ErrInvalidListenAddress
	}

	return nil
}

// Read reads the configuration from the given path
func Read(path string) (*Config, error) {
	// Read the config file
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse it
	var cfg Config

	if err := toml.Unmarshal(content, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
