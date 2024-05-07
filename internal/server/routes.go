package server

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) RegisterRoutes() {
	s.App.Get("/", s.getIndex)
}

func (s *Server) getIndex(c *fiber.Ctx) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	jsonString, err := json.Marshal(s.CurrentForecast)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendString(string(jsonString))
}