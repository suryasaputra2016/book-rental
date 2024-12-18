package utils

import "os"

// func getPort return port
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	return ":" + port
}
