package main

import (
	"log"
	"os"

	"github.com/garciacer87/weatherAPI/src/openweather"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()

	if os.Getenv("OPENWEATHERMAP_HOST") == "" {
		log.Fatal("Cannot init API. Missing environment var: OPENWEATHERMAP_HOST")
	}
	if os.Getenv("OPENWEATHERMAP_APIKEY") == "" {
		log.Fatal("Cannot init API. Missing environment var: OPENWEATHERMAP_APIKEY")
	}
}

func main() {
	host := os.Getenv("OPENWEATHERMAP_HOST")
	apiKey := os.Getenv("OPENWEATHERMAP_APIKEY")

	weatherClient := openweather.NewClient(host, apiKey)
	statusCode, body := weatherClient.GetWeather("", "")

	log.Printf("Body: %s", body)
	log.Printf("Status code: %v", statusCode)
}
