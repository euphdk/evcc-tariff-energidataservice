package server

import (
	"fmt"
	"sync"
	"time"

	"github.com/euphdk/evcc-tariff-energidataservice/internal/config"
	"github.com/gofiber/fiber/v2"
)

type EvccRate struct {
	Start string `json:"start"`
	End   string `json:"end"`
	Price string `json:"price"`
}

type Server struct {
	mu              *sync.Mutex
	Config          config.Config
	App             *fiber.App
	CurrentForecast []*EvccRate
}

func GetServer(conf config.Config) (*Server, error) {
	server := &Server{
		mu:              &sync.Mutex{},
		Config:          conf,
		App:             fiber.New(),
		CurrentForecast: []*EvccRate{},
	}
	server.RegisterRoutes()
	return server, nil
}

func (s *Server) RunBackgroundJobs(done chan error) {
	fmt.Println("Running...")

	for range time.Tick(time.Duration(s.Config.Server.UpdateInterval) * time.Minute) {
		s.mu.Lock()
		fmt.Println("Running... again...")
		s.mu.Unlock()
		done <- fmt.Errorf("blah")
	}
}

func (s *Server) RunApp(done chan error) {
	done <- s.App.Listen(s.Config.Server.Listen)
}
