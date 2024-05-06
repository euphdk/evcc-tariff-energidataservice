package main

import (
	"log/slog"

	"github.com/euphdk/evcc-tariff-energidataservice/internal/config"
	"github.com/euphdk/evcc-tariff-energidataservice/internal/server"
)

func main() {

	c, err := config.GetConfig("config.toml")
	if err != nil {
		slog.Error("Couldn't get config", "error", err.Error())
	}

	slog.Info("Config",
		"server.listen", c.Server.Listen,
		"server.updateinterval", c.Server.UpdateInterval,
		"energidataservice.region", c.Energidataservice.Region,
		"energidataservice.gridcompany", c.Energidataservice.GridCompany,
	)

	done := make(chan error)
	s, err := server.GetServer(c)
	if err != nil {
		slog.Error("Couldn't get server", "error", err.Error())
	}

	go s.RunBackgroundJobs(done)
	go s.RunApp(done)
	err = <-done
	slog.Error(err.Error())

}
