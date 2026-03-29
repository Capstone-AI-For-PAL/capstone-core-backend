package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GenieApiKey string
	GenieAppId  string
	GenieModel  string
	Port        string
}

func LoadEnv() Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not loaded, using system env")
	}

	validateEnv()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	config := Config{
		GenieApiKey: os.Getenv("GENIE_API_KEY"),
		GenieAppId:  os.Getenv("GENIE_APP_ID"),
		GenieModel:  os.Getenv("GENIE_MODEL"),
		Port:        port,
	}
	return config
}

func validateEnv() {
	required := []string{
		"GENIE_API_KEY",
		"GENIE_APP_ID",
		"GENIE_MODEL",
	}

	for _, key := range required {
		if os.Getenv(key) == "" {
			log.Fatalf("Missing required env: %s", key)
		}
	}
}
