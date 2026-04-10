package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GenieApiKey        string
	GenieAppId         string
	GenieModel         string
	GoogleGenAIApiKey  string
	GoogleGenAIModel   string
	S3Bucket           string
	S3Region           string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	Port               string
}

func LoadEnv() Config {
	// Only load .env when env vars aren't already set (e.g. local dev, not Docker)
	if os.Getenv("GENIE_API_KEY") == "" {
		godotenv.Load()
	}

	validateEnv()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("config: port=%s genie_app_id=%s genie_model=%s google_genai_model=%s s3_bucket=%s s3_region=%s aws_access_key_id=%s",
		port,
		os.Getenv("GENIE_APP_ID"),
		os.Getenv("GENIE_MODEL"),
		os.Getenv("GOOGLE_GENAI_MODEL"),
		os.Getenv("AWS_S3_BUCKET"),
		os.Getenv("AWS_S3_REGION"),
		os.Getenv("AWS_ACCESS_KEY_ID"),
	)

	config := Config{
		GenieApiKey:        os.Getenv("GENIE_API_KEY"),
		GenieAppId:         os.Getenv("GENIE_APP_ID"),
		GenieModel:         os.Getenv("GENIE_MODEL"),
		GoogleGenAIApiKey:  os.Getenv("GOOGLE_GENAI_API_KEY"),
		GoogleGenAIModel:   os.Getenv("GOOGLE_GENAI_MODEL"),
		S3Bucket:           os.Getenv("AWS_S3_BUCKET"),
		S3Region:           os.Getenv("AWS_S3_REGION"),
		AWSAccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Port:               port,
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
