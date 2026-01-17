package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Forks           int    `yaml:"forks"`
	Timeout         int    `yaml:"timeout"`
	PrivateKeyFile  string `yaml:"private_key_file"`
	RemoteUser      string `yaml:"remote_user"`
	BecomeMethod    string `yaml:"become_method"`
	BecomeUser      string `yaml:"become_user"`
	HostKeyChecking bool   `yaml:"host_key_checking"`
}

func DefaultConfig() *Config {
	return &Config{
		Forks:           5,
		Timeout:         10,
		RemoteUser:      "root",
		BecomeMethod:    "sudo",
		BecomeUser:      "root",
		HostKeyChecking: false,
	}
}

func LoadConfig() (*Config, error) {
	cfg := DefaultConfig()

	paths := []string{
		"./goansible.yml",
		os.ExpandEnv("$HOME/.goansible/config.yml"),
		"/etc/goansible/config.yml",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return nil, err
			}
			if err := yaml.Unmarshal(data, cfg); err != nil {
				return nil, err
			}
			break
		}
	}

	return cfg, nil
}
