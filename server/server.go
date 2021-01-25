package server

import (
	"os"
	"strconv"

	"github.com/garciacer87/weatherAPI/service"
	"github.com/gin-gonic/gin"
)

//Server impl
type Server struct {
	*gin.Engine
	service service.Service
}

//New returns new gin server
func New() Server {
	host := os.Getenv("OPENWEATHERMAP_HOST")
	apiKey := os.Getenv("OPENWEATHERMAP_APIKEY")

	unit := os.Getenv("OPENWEATHERMAP_UNIT")
	if unit == "" {
		unit = "metric"
	}

	cacheDuration := 2
	d := os.Getenv("CACHE_DURATION")
	if d != "" {
		cacheDuration, _ = strconv.Atoi(d)
	}

	service := service.New(host, apiKey, unit, cacheDuration)
	s := Server{gin.New(), service}

	registerRoutes(s)
	return s
}

func registerRoutes(s Server) {
	s.Group("").GET("/health", HealthCheck)

	s.Group("").
		Use(ValidateRequest()).
		GET("/weather", GetWeather(s.service))
}
