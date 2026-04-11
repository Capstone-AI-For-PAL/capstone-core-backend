package main

import (
	"capstone-llm-service/config"
	"capstone-llm-service/imagegen"
	genie "capstone-llm-service/llm"
	"capstone-llm-service/storage"
	"context"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadEnv()

	genieClient := genie.NewClient(cfg)

	var imgGen imagegen.ImageGenerator
	if cfg.GoogleGenAIApiKey != "" && cfg.S3Bucket != "" && cfg.EnableImageGen {
		store, err := storage.NewS3Storage(cfg.S3Bucket, cfg.S3Region, cfg.AWSAccessKeyID, cfg.AWSSecretAccessKey)
		if err != nil {
			log.Fatalf("Failed to create S3 storage: %v", err)
		}

		genaiClient, err := imagegen.NewGenAIClient(context.Background(), cfg.GoogleGenAIApiKey, cfg.GoogleGenAIModel, store)
		if err != nil {
			log.Fatalf("Failed to create GenAI image client: %v", err)
		}

		imgGen = genaiClient
		log.Println("Image generator configured: Google GenAI (Imagen) + S3")
	} else {
		log.Println("GOOGLE_GENAI_API_KEY or AWS_S3_BUCKET not set, or image generation flag disabled, image generation disabled")
	}

	registerRoutes(http.DefaultServeMux, genieClient, imgGen)

	port := ":" + cfg.Port
	log.Println("server started " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
