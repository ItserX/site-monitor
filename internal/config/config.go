package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig(path string, out interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, out)
	if err != nil {
		return err
	}
	return nil
}
