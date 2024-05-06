package main

import (
	"fmt"
	"log/slog"

	"github.com/euphdk/evcc-tariff-energidataservice/internal/config"
)

func main() {

	config, err := config.GetConfig("config.toml")
	if err != nil {
		slog.Error("Couldn't get config", "error", err.Error())
	}

	fmt.Printf("%#v", config)
	
}