package config

import (
	"log"
	"log/slog"
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
	EnableImageGen     bool
}

func LoadEnv() Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("config: unable to load .env file: %v; continuing with existing environment", err)
	}

	validateEnv()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

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
		EnableImageGen:     os.Getenv("ENABLE_IMAGE_GEN") == "true",
	}

	log.Printf("config: port=%s genie_app_id=%s genie_model=%s google_genai_api_key=%s google_genai_model=%s enable_image_gen=%t s3_bucket=%s s3_region=%s aws_access_key_id=%s",
		port,
		redactConfig(config.GenieAppId),
		config.GenieModel,
		redactConfig(config.GoogleGenAIApiKey),
		config.GoogleGenAIModel,
		config.EnableImageGen,
		config.S3Bucket,
		config.S3Region,
		redactConfig(config.AWSAccessKeyID),
	)
	return config
}

func redactConfig(text string) slog.Value {
	if text == "" {
		return slog.StringValue("")
	}
	if len(text) > 2 {
		redacted := text[:1] + "..." + text[len(text)-1:]
		return slog.StringValue(redacted)
	}
	return slog.StringValue("REDACTED")
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
