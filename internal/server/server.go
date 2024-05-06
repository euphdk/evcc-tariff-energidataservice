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
	Mutex           *sync.Mutex
	App             *fiber.App
	CurrentForecast []*EvccRate
}

func GetServer(conf config.Config) (Server, error) {
	server := Server{
		Mutex:           &sync.Mutex{},
		CurrentForecast: []*EvccRate{},
	}
	return server, nil
}

func (s *Server) RunBackground(done chan error) {
	fmt.Println("Running...")

	for range time.Tick(60 * time.Second) {
		fmt.Println("Running...")
		done <- fmt.Errorf("blah")
	}

}
