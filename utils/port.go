package utils

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

// func getPort returns port
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	return ":" + port
}
