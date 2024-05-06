package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/euphdk/evcc-tariff-energidataservice/internal/config"
	"github.com/euphdk/evcc-tariff-energidataservice/internal/server"
)

func main() {

	c, err := config.GetConfig("config.toml")
	if err != nil {
		slog.Error("Couldn't get config", "error", err.Error())
	}

	fmt.Printf("%#v\n", c)

	done := make(chan error)
	s, err := server.GetServer(c)
	if err != nil {
		slog.Error("Couldn't get server", "error", err.Error())
	}

	go s.RunBackground(done)

	time.Sleep(120 * time.Second)

	
}