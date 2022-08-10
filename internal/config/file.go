package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

func ReadRawConfigFromPath(configPath string) (map[string]interface{}, error) {
	if _, err := os.Stat(configPath); err != nil {
		return nil, err
	}
	yamlBytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var config map[string]interface{}
	err = yaml.UnmarshalStrict(yamlBytes, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
