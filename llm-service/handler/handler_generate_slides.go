package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"capstone-llm-service/imagegen"
	genie "capstone-llm-service/llm"

	"golang.org/x/sync/errgroup"
)

// TODO: Calculate cost of generation and log for monitoring. Genie charges per input and output tokens, and GenAI charges per image generated. This will help us understand expenses and optimize prompts.
// systemSlidesPrompt is a fixed system instruction that enforces the JSON schema
// and layout rules. It is always sent as the system message regardless of user input.
func systemSlidesPrompt() string {
	return "You are a slide content generator. Your ONLY output is a raw JSON array.\n" +
		"Return ONLY a JSON array with 3 to 5 slide objects in this exact format:\n" +
		"[{\"title\":\"...\",\"layout\":\"hero_image|two_column|three_column|text_only\"," +
		"\"key_points\":[\"...\",\"...\"]," +
		"\"images\":[{\"prompt\":\"...\",\"caption\":\"...\"}]," +
		"\"transcript\":\"...\"}]\n\n" +
		"Layout rules:\n" +
		"- \"hero_image\": Use for intro/overview slides. Exactly 1 image.\n" +
		"- \"two_column\": Use for comparisons or two related concepts. Exactly 2 images.\n" +
		"- \"three_column\": Use for categories or three examples. Exactly 3 images.\n" +
		"- \"text_only\": Use for definitions, formulas, or content with no useful visual. 0 images.\n\n" +
		"Content rules:\n" +
		"- key_points: 2-4 short phrases per slide. Keep text minimal.\n" +
		"- images[].prompt: Write a vivid, specific description suitable for an AI image generator. Focus on educational clarity.\n" +
		"- images[].caption: 2-5 word label for the image.\n" +
		"- transcript: Write as if lecturing to students. Explain concepts, give examples, provide context. " +
		"Should be 3-5 sentences and must NOT simply restate the key_points.\n\n" +
		"Do not output markdown or code fences. Output raw JSON only."
}

// defaultGuidancePrompt is used when the user does not supply a custom guidance prompt.
// It provides general content and style direction to the model.
func defaultGuidancePrompt() string {
	return "Act as an expert curriculum designer and instructional strategist.\n" +
		"Generate detailed slide content for ONLY the current slide topic.\n" +
		"Use outline_context and rag_data deeply. Prioritize educational clarity and engagement."
}

type sectionResult struct {
	sectionIndex int
	title        string
	slides       []EnrichedSlide
}

type imageKey struct {
	sectionIdx int
	slideIdx   int
	imageIdx   int
}

func MakeGenerateSlidesHandler(client *genie.Client, imgGen imagegen.ImageGenerator) http.HandlerFunc {
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
			http.Error(w, "missing required fields: sections", http.StatusBadRequest)
			return
		}

		log.Printf("Generating slides for %s (%d sections)", req.CunetId, len(req.Sections))
		email := req.CunetId + "@student.chula.ac.th"
		guidancePrompt := req.Prompt
		if strings.TrimSpace(guidancePrompt) == "" {
			guidancePrompt = defaultGuidancePrompt()
		}

		// --- Phase 1: Content Generation ---
		sectionResults := make([]sectionResult, 0, len(req.Sections))

		for sectionIdx, section := range req.Sections {
			inputPayload := map[string]interface{}{
				"slide":           section,
				"outline":         req.Outline,
				"outline_context": req.Outline,
				"rag_data":        req.RagData,
			}

			payloadBytes, err := json.Marshal(inputPayload)
			if err != nil {
				log.Printf("Error marshaling slide input (section %d): %v", sectionIdx, err)
				http.Error(w, "failed to prepare slide input", http.StatusInternalServerError)
				return
			}

			messageText := fmt.Sprintf("%s\n\nCURRENT INPUT JSON:\n%s", guidancePrompt, string(payloadBytes))
			messages := []genie.Message{
				{
					Role: "system",
					Content: []genie.ContentPart{
						{Type: "text", Text: systemSlidesPrompt()},
					},
				},
				{
					Role: "user",
					Content: []genie.ContentPart{
						{Type: "text", Text: messageText},
					},
				},
			}

			res, err := client.Chat(messages, email, req.CunetId)
			if err != nil {
				log.Printf("Slide generation error (section %d): %v", sectionIdx+1, err)
				continue
			}

			enrichedSlides, err := decodeEnrichedSlides(res)
			if err != nil {
				log.Printf("Invalid enriched slide JSON (section %d): %v", sectionIdx+1, err)
				continue
			}
			log.Printf("Generated section %d: %s with %d slides", sectionIdx+1, section.Title, len(enrichedSlides))

			sectionResults = append(sectionResults, sectionResult{
				sectionIndex: sectionIdx,
				title:        section.Title,
				slides:       enrichedSlides,
			})
		}

		// --- Phase 2: Image Generation ---
		type imageTask struct {
			key imageKey
			req imagegen.ImageRequest
		}

		var tasks []imageTask
		for _, sr := range sectionResults {
			for slideIdx, slide := range sr.slides {
				imgWidth, imgHeight := ImageDimensionsForLayout(slide.Layout)
				for imgIdx, img := range slide.Images {
					path := fmt.Sprintf("slides/%s/section_%d/slide_%d/img_%d.png",
						req.CunetId, sr.sectionIndex, slideIdx, imgIdx)
					tasks = append(tasks, imageTask{
						key: imageKey{sr.sectionIndex, slideIdx, imgIdx},
						req: imagegen.ImageRequest{
							ImagePath: path,
							Prompt:    img.Prompt,
							Width:     imgWidth,
							Height:    imgHeight,
						},
					})
				}
			}
		}

		imageURLs := make(map[imageKey]string)
		if len(tasks) > 0 && imgGen != nil {
			var mu sync.Mutex
			g, ctx := errgroup.WithContext(r.Context())
			for _, t := range tasks {
				t := t
				g.Go(func() error {
					url, err := imgGen.Generate(ctx, t.req)
					if err != nil {
						log.Printf("Image generation failed (%s): %v", t.req.ImagePath, err)
						return nil
					}
					mu.Lock()
					imageURLs[t.key] = url
					mu.Unlock()
					return nil
				})
			}
			_ = g.Wait()
		}
		log.Printf("Generated %d/%d images", len(imageURLs), len(tasks))

		// --- Phase 3: Assembly ---
		var sections []SectionOutput
		var transcripts []TranscriptEntry
		slideCounter := 0

		transcripts = append(transcripts, TranscriptEntry{
			SlideIndex: slideCounter,
			SlideTitle: req.Lesson.Title,
			Text:       fmt.Sprintf("Welcome to today's lesson on %s.", req.Lesson.Title),
		})
		slideCounter++

		for _, sr := range sectionResults {
			var renderedSlides []RenderedSlide

			transcripts = append(transcripts, TranscriptEntry{
				SlideIndex: slideCounter,
				SlideTitle: sr.title,
				Text:       fmt.Sprintf("In this section, we will cover %s.", sr.title),
			})
			slideCounter++

			for slideIdx, slide := range sr.slides {
				imgWidth, imgHeight := ImageDimensionsForLayout(slide.Layout)
				var renderedImages []RenderedImage
				for imgIdx, img := range slide.Images {
					url := imageURLs[imageKey{sr.sectionIndex, slideIdx, imgIdx}]
					if url == "" {
						continue
					}
					renderedImages = append(renderedImages, RenderedImage{
						URL:     url,
						Caption: img.Caption,
						Width:   imgWidth,
						Height:  imgHeight,
					})
				}

				rendered := RenderedSlide{
					Title:      slide.Title,
					Layout:     slide.Layout,
					KeyPoints:  slide.KeyPoints,
					Images:     renderedImages,
					Transcript: slide.Transcript,
				}
				renderedSlides = append(renderedSlides, rendered)

				transcripts = append(transcripts, TranscriptEntry{
					SlideIndex: slideCounter,
					SlideTitle: slide.Title,
					Text:       slide.Transcript,
				})
				slideCounter++
			}

			sections = append(sections, SectionOutput{
				Title:  sr.title,
				Slides: renderedSlides,
			})
		}

		markdown := generateMarkdown(sections, req.Lesson)
		log.Printf("Generated markdown (%d bytes), %d transcripts", len(markdown), len(transcripts))

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(SlideGenerationResponse{
			Sections:    sections,
			Markdown:    markdown,
			Transcripts: transcripts,
		})
	}
}

func generateMarkdown(sections []SectionOutput, lesson LessonInput) string {
	var sb strings.Builder
	sb.WriteString(marpFrontmatter())
	sb.WriteString(fmt.Sprintf("\n# Lesson: %s\n", lesson.Title))

	for _, section := range sections {
		sb.WriteString("\n\n---\n")
		sb.WriteString(fmt.Sprintf("# %s\n", section.Title))
		for _, slide := range section.Slides {
			sb.WriteString("\n---\n")
			sb.WriteString(renderSlide(slide))
		}
	}

	return sb.String()
}
