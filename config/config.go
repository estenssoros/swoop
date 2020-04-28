package config

import (
	"encoding/json"
	"os"

	"github.com/estenssoros/swoop/resource"
	"github.com/estenssoros/swoop/secret"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Config config for swoop process
type Config struct {
	Source         *resource.Resource
	Destination    *resource.Resource
	SecretProvider *secret.Provider `yaml:"secretProvider"`
	Tables         []*Table
	Truncate       *bool
	WriteLimit     *int `yaml:"writeLimit"`
}

func (c Config) String() string {
	ju, _ := json.MarshalIndent(c, "", " ")
	return string(ju)
}

// NewFromFile create new config from file
func NewFromFile(fileName string) (*Config, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, errors.Wrap(err, "open file")
	}
	defer f.Close()
	c := &Config{}
	return c, yaml.NewDecoder(f).Decode(c)
}

// NewFromString create new config from string
func NewFromString(yml string) (*Config, error) {
	c := &Config{}
	return c, yaml.Unmarshal([]byte(yml), c)
}

// Connect connect resources in a config
func (c *Config) Connect() error {
	if err := c.Source.Connect(c.SecretProvider); err != nil {
		return errors.Wrap(err, "connect source")
	}
	if err := c.Destination.Connect(c.SecretProvider); err != nil {
		return errors.Wrap(err, "connect destination")
	}
	return nil
}
