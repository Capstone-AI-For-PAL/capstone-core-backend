package imagegen

import (
	"context"
	"fmt"
	"log"
	"strings"

	"capstone-llm-service/storage"

	"google.golang.org/genai"
)

const defaultImageModel = "gemini-2.5-flash-image"

type GenAIClient struct {
	client    *genai.Client
	model     string
	store     storage.Storage
	imageSize string
}

func NewGenAIClient(ctx context.Context, apiKey string, imageSize string, model string, store storage.Storage) (*GenAIClient, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("create genai client: %w", err)
	}

	if model == "" {
		model = defaultImageModel
	}

	log.Printf("GenAI client created with model %q", model)

	return &GenAIClient{
		client:    client,
		model:     model,
		store:     store,
		imageSize: imageSize,
	}, nil
}

func (g *GenAIClient) Generate(ctx context.Context, req ImageRequest) (string, error) {
	return g.generateGemini(ctx, req)
}

// generateGemini uses GenerateContent with IMAGE response modality.
// Works with Gemini models that support image output (e.g. gemini-2.0-flash-exp).
func (g *GenAIClient) generateGemini(ctx context.Context, req ImageRequest) (string, error) {
	contents := []*genai.Content{
		genai.NewContentFromText(req.Prompt, "user"),
	}

	config := &genai.GenerateContentConfig{
		ResponseModalities: []string{"IMAGE"},
		ImageConfig: &genai.ImageConfig{
			AspectRatio: aspectRatioFromDimensions(req.Width, req.Height),
			ImageSize:   g.imageSize,
		},
	}

	resp, err := g.client.Models.GenerateContent(ctx, g.model, contents, config)
	if err != nil {
		return "", fmt.Errorf("genai generate content (image): %w", err)
	}

	log.Printf("GenAI usage metadata: %v", resp.UsageMetadata.CandidatesTokenCount)

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return "", fmt.Errorf("genai returned no candidates")
	}

	for _, part := range resp.Candidates[0].Content.Parts {
		if part.InlineData != nil && strings.HasPrefix(part.InlineData.MIMEType, "image/") {
			url, err := g.store.Upload(ctx, req.ImagePath, part.InlineData.Data, part.InlineData.MIMEType)
			if err != nil {
				return "", fmt.Errorf("upload generated image: %w", err)
			}
			return url, nil
		}
	}

	return "", fmt.Errorf("genai response contained no image data")
}

// aspectRatioFromDimensions maps pixel dimensions to the closest Imagen-supported
// aspect ratio: "1:1", "3:4", "4:3", "9:16", "16:9".
func aspectRatioFromDimensions(w, h int) string {
	if w <= 0 || h <= 0 {
		return "1:1"
	}

	ratio := float64(w) / float64(h)

	type ar struct {
		label string
		value float64
	}
	supported := []ar{
		{"9:16", 9.0 / 16.0},
		{"3:4", 3.0 / 4.0},
		{"1:1", 1.0},
		{"4:3", 4.0 / 3.0},
		{"16:9", 16.0 / 9.0},
	}

	best := supported[0]
	bestDiff := abs(ratio - best.value)
	for _, s := range supported[1:] {
		diff := abs(ratio - s.value)
		if diff < bestDiff {
			best = s
			bestDiff = diff
		}
	}
	return best.label
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
