package cmd

import (
	"io/ioutil"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Owner       string    `yaml:"owner"`
	Repo        string    `yaml:"repo"`
	LastUpdated time.Time `yaml:"last_updated"`
}

func (c *Config) Save() {
	c.LastUpdated = time.Now()
	b, _ := yaml.Marshal(c)
	ioutil.WriteFile("./.ghi.yml", b, 0755)
}

func (c *Config) SetFromArgs(args []string) {
	sp := strings.Split(args[0], "/")
	c.Owner = sp[0]
	c.Repo = sp[1]
	c.Save()
}

func LoadConfig() *Config {
	c := &Config{}
	b, err := ioutil.ReadFile("./.ghi.yml")
	if err != nil {
		return c
	}

	yaml.Unmarshal(b, c)
	return c
}
