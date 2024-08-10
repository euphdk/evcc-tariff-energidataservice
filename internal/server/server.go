package server

import (
	"log/slog"
	"sync"
	"time"

	"github.com/euphdk/evcc-tariff-energidataservice/internal/config"
	"github.com/euphdk/evcc-tariff-energidataservice/internal/energidataservice"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	mu              *sync.Mutex
	Config          config.Config
	App             *fiber.App
	CurrentForecast []energidataservice.EvccAPIRate
}

func GetServer(conf config.Config) *Server {
	server := &Server{
		mu:              &sync.Mutex{},
		Config:          conf,
		App:             fiber.New(),
		CurrentForecast: []energidataservice.EvccAPIRate{},
	}
	server.RegisterRoutes()
	return server
}

func (s *Server) RunBackgroundJobs(done chan error) {

	tick := time.NewTicker(time.Duration(s.Config.Server.UpdateInterval) * time.Minute)
	for ; true; <-tick.C {
		currentForecast, err := energidataservice.GetEvccAPIRates(s.Config.Energidataservice.GridCompany, s.Config.Energidataservice.Region, s.Config.Energidataservice.TAX, s.Config.Energidataservice.VAT)
		if err != nil {
			slog.Error(err.Error())
		} else {
			s.mu.Lock()
			s.CurrentForecast = currentForecast
			s.mu.Unlock()
		}
	}
}

func (s *Server) RunApp(done chan error) {
	done <- s.App.Listen(s.Config.Server.Listen)
}
