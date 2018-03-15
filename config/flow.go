package config

import (
	"io/ioutil"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type FlowConfig struct {
	Label     string          `yaml:"label"`
	When      []time.Duration `yaml:"when"` // time.ParseDuration
	Collector string          `yaml:"collector"`
	ClientCfg ClientConfig    `yaml:"client"`
	ServerCfg ServerConfig    `yaml:"server"`
}

func NewFlowConfigFromYaml(buf []byte) (*FlowConfig, error) {
	var t FlowConfig

	err := yaml.Unmarshal(buf, &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func NewFlowConfigFromFile(path string) (*FlowConfig, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return NewFlowConfigFromYaml(buf)
}
