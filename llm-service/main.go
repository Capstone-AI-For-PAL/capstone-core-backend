package main

import (
	"capstone-llm-service/config"
	"capstone-llm-service/imagegen"
	genie "capstone-llm-service/llm"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadEnv()

	genieClient := genie.NewClient(cfg)

	var imgGen imagegen.ImageGenerator
	if cfg.ImageApiBaseURL != "" {
		imgGen = imagegen.NewHTTPClient(cfg.ImageApiBaseURL)
		log.Printf("Image generator configured: %s", cfg.ImageApiBaseURL)
	} else {
		log.Println("IMAGE_API_BASE_URL not set, image generation disabled")
	}

	registerRoutes(http.DefaultServeMux, genieClient, imgGen)

	port := ":" + cfg.Port
	log.Println("server started " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
