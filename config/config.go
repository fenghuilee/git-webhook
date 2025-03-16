package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port   int    `yaml:"port"`
		Secret string `yaml:"secret"`
	} `yaml:"server"`
	Projects []Project `yaml:"projects"`
}

type Project struct {
	Name     string   `yaml:"name"`
	Path     string   `yaml:"path"`
	Branch   string   `yaml:"branch"`
	Secret   string   `yaml:"secret"`
	Commands []string `yaml:"commands"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
