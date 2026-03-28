package main

import (
	"net/http"

	handler "capstone-llm-service/handler"
	genie "capstone-llm-service/llm"
)

func registerRoutes(mux *http.ServeMux, client *genie.Client) {
	mux.HandleFunc("/chat", handler.MakeChatHandler(client))
	mux.HandleFunc("/generate-slides", handler.MakeGenerateSlidesHandler(client))
}
