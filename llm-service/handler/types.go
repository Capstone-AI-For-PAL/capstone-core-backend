package handler

import genie "capstone-llm-service/llm"

type ChatRequest struct {
	Messages []genie.Message `json:"messages"`
	CunetId  string          `json:"cunet_id"`
}

type SectionInput struct {
	Title       string `json:"slide_title"`
	MainContent string `json:"main_content"`
	Objective   string `json:"objective"`
}

type SlideGenerationRequest struct {
	Prompt   string         `json:"prompt"`
	Sections []SectionInput `json:"sections"`
	Lesson   LessonInput    `json:"lesson"`
	Outline  string         `json:"outline"`
	RagData  string         `json:"rag_data"`
	CunetId  string         `json:"cunet_id"`
}

type LessonInput struct {
	Title       string `json:"title"`
	MainContent string `json:"main_content"`
	Objective   string `json:"objective"`
}

type OrchestratedSlideBullet struct {
	Bullet    string   `json:"bullet"`
	SubPoints []string `json:"sub_points"`
}

type OrchestratedSlide struct {
	Title   string                    `json:"title"`
	Content []OrchestratedSlideBullet `json:"content"`
}
type MarkdownInput struct {
	Title  string               `json:"title"`
	Slides []SlideMarkdownInput `json:"slides"`
}

type SlideMarkdownInput struct {
	SlideHeader string   `json:"slide_header"`
	SubPoints   []string `json:"sub_points"`
}

type SlideGenerationResponse struct {
	Data     []MarkdownInput `json:"data"`
	Markdown string          `json:"markdown"`
}
