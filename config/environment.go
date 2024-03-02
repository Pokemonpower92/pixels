package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnvironmentVariables() {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Find the root directory of the project
	rootDir := filepath.Dir(filepath.Dir(wd))

	// Load the .env file
	envFile := filepath.Join(rootDir, ".env")
	err = godotenv.Load(envFile)
	if err != nil {
		fmt.Println(".env file not found")
	}
}
