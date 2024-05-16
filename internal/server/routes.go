package server

import (
	"github.com/gofiber/fiber/v2"
)

func (s *Server) RegisterRoutes() {
	s.App.Get("/", s.getIndex)
}

func (s *Server) getIndex(c *fiber.Ctx) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return c.JSON(s.CurrentForecast)
}