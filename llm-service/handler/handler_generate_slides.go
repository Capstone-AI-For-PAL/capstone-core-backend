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

		if len(req.Slides) == 0 {
			http.Error(w, "missing required fields: slides", http.StatusBadRequest)
			return
		}

		log.Println("Generating slide for ", req.CunetId, " with prompt: ", req.Prompt)
		email := req.CunetId + "@student.chula.ac.th"
		promptTemplate := req.Prompt
		if strings.TrimSpace(promptTemplate) == "" {
			promptTemplate = defaultSlidesPrompt()
		}

		results := make([]OrchestratedSlide, 0)

		for slideNumber, slide := range req.Slides {
			inputPayload := map[string]interface{}{
				"slide":           slide,
				"outline":         req.OutlineContext,
				"outline_context": req.OutlineContext,
				"rag_data":        req.RagData,
			}

			payloadBytes, err := json.Marshal(inputPayload)
			if err != nil {
				log.Printf("Error marshaling slide input (slide %d): %v", slideNumber, err)
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
			if err != nil {
				log.Printf("Slide generation error (slide %d): %v", slideNumber, err)
				continue
			}

			generatedSlides, err := decodeOrchestratedSlides(res)
			if err != nil {
				log.Printf("Invalid slide JSON from model (slide %d): %v", slideNumber, err)
				continue
			}

			for slideIdx, genSlide := range generatedSlides {
				results = append(results, OrchestratedSlide{
					Title:   genSlide.Title,
					Content: genSlide.Content,
				})
				log.Printf("Generated slide %d.%d: %s with %d bullets", slideNumber, slideIdx+1, genSlide.Title, len(genSlide.Content))
			}
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(SlideGenerationResponse{Results: results})
	}
}
