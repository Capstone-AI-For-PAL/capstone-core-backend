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
	LessonId string         `json:"lesson_id"`
}

type LessonInput struct {
	Title       string `json:"title"`
	MainContent string `json:"main_content"`
	Objective   string `json:"objective"`
}

type LayoutType string

const (
	LayoutHeroImage   LayoutType = "hero_image"
	LayoutTwoColumn   LayoutType = "two_column"
	LayoutThreeColumn LayoutType = "three_column"
	LayoutTextOnly    LayoutType = "text_only"
)

type EnrichedSlide struct {
	Title          string        `json:"title"`
	Layout         LayoutType    `json:"layout"`
	KeyPoints      []string      `json:"key_points"`
	Images         []ImagePrompt `json:"images"`
	Transcript     string        `json:"transcript"`
	ThaiTranscript string        `json:"thai_transcript"`
}

type ImagePrompt struct {
	Prompt  string `json:"prompt"`
	Caption string `json:"caption"`
}

type RenderedSlide struct {
	Title          string          `json:"title"`
	Layout         LayoutType      `json:"layout"`
	KeyPoints      []string        `json:"key_points"`
	Images         []RenderedImage `json:"images"`
	Transcript     string          `json:"transcript"`
	ThaiTranscript string          `json:"thai_transcript"`
}

type RenderedImage struct {
	URL     string `json:"url"`
	Caption string `json:"caption"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
}

type TranscriptEntry struct {
	SlideIndex     int    `json:"slide_index"`
	SlideTitle     string `json:"slide_title"`
	Transcript     string `json:"transcript"`
	ThaiTranscript string `json:"thai_transcript"`
}

type SectionOutput struct {
	Title  string          `json:"title"`
	Slides []RenderedSlide `json:"slides"`
}

type SlideGenerationResponse struct {
	Sections    []SectionOutput   `json:"sections"`
	Markdown    string            `json:"markdown"`
	Transcripts []TranscriptEntry `json:"transcripts"`
}
