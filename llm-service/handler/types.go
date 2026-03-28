package handler

import genie "capstone-llm-service/llm"

type ChatRequest struct {
	Messages []genie.Message `json:"messages"`
	CunetId  string          `json:"cunet_id"`
}

type SlideInput struct {
	Title       string `json:"slide_title"`
	MainContent string `json:"main_content"`
	Objective   string `json:"objective"`
}

type SlideGenerationRequest struct {
	Prompt         string       `json:"prompt"`
	Slides         []SlideInput `json:"slides"`
	OutlineContext string       `json:"outline_context"`
	RagData        string       `json:"rag_data"`
	CunetId        string       `json:"cunet_id"`
}

type OrchestratedSlideBullet struct {
	Bullet    string   `json:"bullet"`
	SubPoints []string `json:"sub_points"`
}

type OrchestratedSlide struct {
	Title   string                    `json:"title"`
	Content []OrchestratedSlideBullet `json:"content"`
}

type SlideGenerationResponse struct {
	Results []OrchestratedSlide `json:"results"`
}
