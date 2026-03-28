package main

import (
	"capstone-llm-service/config"
	genie "capstone-llm-service/llm"
	"log"
	"net/http"
)

func main() {
	config := config.LoadEnv()

	genieClient := genie.NewClient(config)
	registerRoutes(http.DefaultServeMux, genieClient)

	port := config.Port
	log.Println("🚀 server started " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
