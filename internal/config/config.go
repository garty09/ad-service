package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"ad-service/internal/ad"
)

const defaultServerPort         = 8080

type Config struct {
	ServerPort int `yaml:"server_port"`
	// the data source name (DSN) for connecting to the database
	DSN string `yaml:"dsn"`
}

func (c Config) Validate() error {
	if c.DSN == "" {
		return ad.ErrIsEmpty
	}
	return nil
}

func Load(file string) (*Config, error) {
	// default config
	c := Config{
		ServerPort:    defaultServerPort,
	}

	// load from YAML config file
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(bytes, &c); err != nil {
		return nil, err
	}

	// validation
	if err = c.Validate(); err != nil {
		return nil, err
	}

	return &c, err
}
