package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	genie "capstone-llm-service/llm"
)

func defaultSlidesPrompt() string {
	return "Act as an expert curriculum designer and instructional strategist.\n" +
		"Generate detailed, easy-to-understand slide content for ONLY the current slide topic.\n" +
		"Return ONLY a JSON array with 3 to 5 slide objects in this format: " +
		"[{\"title\":\"...\",\"content\":[{\"bullet\":\"...\",\"sub_points\":[\"...\",\"...\",\"...\"]}]}].\n" +
		"Each generated slide must contain exactly 4 main bullets, and each main bullet must contain exactly 3 sub_points.\n" +
		"Include practical examples, concrete details, and explanations tailored for learners.\n" +
		"Use outline_context and rag_data deeply, do not output markdown or code fences. The outline_context is the overall course outline, and rag_data is additional relevant information that can be used to enrich the slide content."
}

func MakeGenerateSlidesHandler(client *genie.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
			http.Error(w, "content type must be application/json", http.StatusBadRequest)
			return
		}

		var req SlideGenerationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.CunetId == "" {
			http.Error(w, "missing required fields: cunet_id", http.StatusBadRequest)
			return
		}

		if len(req.Sections) == 0 {
			http.Error(w, "missing required fields: slides", http.StatusBadRequest)
			return
		}

		log.Println("Generating slide for ", req.CunetId, " with prompt: ", req.Prompt)
		email := req.CunetId + "@student.chula.ac.th"
		promptTemplate := req.Prompt
		if strings.TrimSpace(promptTemplate) == "" {
			promptTemplate = defaultSlidesPrompt()
		}

		results := make([]MarkdownInput, 0, len(req.Sections))

		for sectionNumber, slide := range req.Sections {
			inputPayload := map[string]interface{}{
				"slide":           slide,
				"outline":         req.Outline,
				"outline_context": req.Outline,
				"rag_data":        req.RagData,
			}

			payloadBytes, err := json.Marshal(inputPayload)
			if err != nil {
				log.Printf("Error marshaling slide input (slide %d): %v", sectionNumber, err)
				http.Error(w, "failed to prepare slide input", http.StatusInternalServerError)
				return
			}

			messageText := fmt.Sprintf("%s\n\nCURRENT INPUT JSON:\n%s", promptTemplate, string(payloadBytes))
			messages := []genie.Message{
				{
					Role: "user",
					Content: []genie.ContentPart{
						{Type: "text", Text: messageText},
					},
				},
			}

			res, err := client.Chat(messages, email, req.CunetId)
			log.Println("Generated slide: ", res)
			if err != nil {
				log.Printf("Slide generation error (section %d): %v", sectionNumber+1, err)
				continue
			}

			generatedSlides, err := decodeOrchestratedSlides(res)
			if err != nil {
				log.Printf("Invalid slide JSON from model (section %d): %v", sectionNumber+1, err)
				continue
			}
			log.Printf("Generated section %d: %s with %d slides", sectionNumber+1, slide.Title, len(generatedSlides))

			results = append(results, MarkdownInput{
				Title:  slide.Title,
				Slides: convertToMarkdownInput(generatedSlides),
			})
			log.Printf("Generated markdown for section %d with length %d", sectionNumber+1, len(results[len(results)-1].Slides))
		}

		log.Println("Generating markdown")
		markdown := generateMarkdown(results, req.Lesson)
		log.Println("Outputting as markdown with length: ", len(markdown))

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(SlideGenerationResponse{
			Data:     results,
			Markdown: markdown,
		})
	}
}

func convertToMarkdownInput(generatedSlides []OrchestratedSlide) []SlideMarkdownInput {
	markdownInputs := make([]SlideMarkdownInput, 0)
	for _, s := range generatedSlides {
		for _, b := range s.Content {
			markdownInputs = append(markdownInputs, SlideMarkdownInput{
				SlideHeader: b.Bullet,
				SubPoints:   b.SubPoints,
			})
		}
	}
	return markdownInputs
}

func generateMarkdown(results []MarkdownInput, lesson LessonInput) string {
	var sb strings.Builder
	sb.WriteString("---\nmarp: true\ntheme: default\nmath: mathjax\npaginate: true\nsize: 16:9\nstyle: |\n  section {\n    font-size: 25px;\n    padding: 40px;\n    justify-content: center; /* Keeps content centered vertically */\n  }\n  h1 {\n    font-size: 40px;\n    color: #0288d1;\n  }\n  h2 {\n    font-size: 35px;\n    color: #333;\n  }\n---\n")
	sb.WriteString(fmt.Sprintf("# Lesson: %s\n", lesson.Title))

	for _, s := range results {
		sb.WriteString("\n\n---\n")
		sb.WriteString(fmt.Sprintf("# %s\n", s.Title))
		for _, b := range s.Slides {
			sb.WriteString("\n---\n")
			bulletHeader := fmt.Sprintf("## %s\n\n", b.SlideHeader)
			sb.WriteString(bulletHeader)

			wordCount := 0

			for _, sp := range b.SubPoints {
				spWordCount := len(strings.Fields(sp))

				// If adding this subpoint exceeds 150 words AND the page already has some subpoints, otherwise leave it overflow
				if wordCount+spWordCount > 150 && wordCount > 0 {
					sb.WriteString("\n---\n")
					sb.WriteString(bulletHeader)
					wordCount = 0
				}

				sb.WriteString(fmt.Sprintf("* %s\n", sp))
				wordCount += spWordCount
			}
		}
	}

	return sb.String()
}
