package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

func LoadEnvironmentVariables() {
	// Load the .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(".env file not found. Using default environment variables.")
	}
}
