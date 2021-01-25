package server

import (
	"encoding/json"
	"net/http"

	"github.com/garciacer87/weatherAPI/service"
	"github.com/gin-gonic/gin"
)

//HealthCheck handler used as a health
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "UP"})
}

//GetWeather handler used to get weather info
func GetWeather(srv service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		city := c.Query("city")
		country := c.Query("country")

		respCode, respBody := srv.GetWeather(city, country)

		var body interface{}
		json.Unmarshal(respBody, &body)

		c.JSON(respCode, body)
	}
}
