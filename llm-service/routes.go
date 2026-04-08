package main

import (
	"net/http"

	handler "capstone-llm-service/handler"
	"capstone-llm-service/imagegen"
	genie "capstone-llm-service/llm"
)

func registerRoutes(mux *http.ServeMux, client *genie.Client, imgGen imagegen.ImageGenerator) {
	mux.HandleFunc("/chat", handler.MakeChatHandler(client))
	mux.HandleFunc("/generate-slides", handler.MakeGenerateSlidesHandler(client, imgGen))
}
