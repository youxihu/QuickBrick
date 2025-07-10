package config

import (
	"gopkg.in/yaml.v3"
	"os"

	"QuickBrick/internal/domain"
)

var Cfg domain.Config

func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &Cfg)
}
