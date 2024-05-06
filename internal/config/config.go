package config

import (
	"io"
	"os"

	"github.com/pelletier/go-toml"
)

type Config struct {
	Server struct {
		Listen         string `toml:"listen"`
		UpdateInterval int    `toml:"updateinterval"`
	} `toml:"server"`
	Energidataservice struct {
		Region      string
		GridCompany string
	} `toml:"energidataservice"`
}

func GetConfig(configfile string) (Config, error) {
	var config Config

	file, err := os.Open(configfile)

	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return Config{}, err
	}

	err = toml.Unmarshal(fileContent, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
