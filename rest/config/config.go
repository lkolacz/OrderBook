package config

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type httpSettings struct {
	Addr                 string   `yaml:"addr"`
	CORSAllowOrigins     []string `yaml:"cors_allow_origins"`
	ProxyForwardedHeader string   `yaml:"proxy_forwarded_header"`
	LogAllRequests       bool     `yaml:"log_all_requests"`
}

// Config is the service configuration
type Config struct {
	Base    Base         `yaml:"base"`
	HTTP    httpSettings `yaml:"http"`
	Logging Logging      `yaml:"logging"`
}

type Base struct {
	Domain   string `yaml:"domain"`
	BasePath string `yaml:"base_path"`
}

type Logging struct {
	Format string `yaml:"format"`
	Level  string `yaml:"level"`
}

// Default creates configuration struct with default values
func (c *Config) Default() {
	c.HTTP = httpSettings{
		Addr:             "127.0.0.1:8007",
		CORSAllowOrigins: []string{"*"},
	}
	c.Base.Domain = "api.order-book.network"
	c.Base.BasePath = "/api/v1"
	c.Logging.Level = "info"
	c.Logging.Format = "json"
}

// Load & Parse ConfigFile (yaml config)
func (c *Config) Load(fileName string) error {
	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		return errors.Wrap(err, "Can not read the config file.")
	}

	c.Default()
	if err := yaml.Unmarshal(f, c); err != nil {
		return errors.Wrap(err, "Configuration Parser Error.")
	}

	return nil
}

// Save Configuration (to yaml file)
func (c *Config) Save(fileName string) error {
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, b, 0600)
}
