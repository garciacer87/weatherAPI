package main

import (
	"fmt"
	"log"
	"os"

	"github.com/garciacer87/weatherAPI/server"
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
	s := server.New()

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		log.Println("No port was found on SERVER_PORT env. Using default port 8080")
		port = "8080"
	}

	err := s.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Error trying to serve application: %v", err)
	}
}
