package imagegen

import "context"

// ImageGenerator abstracts image generation so implementations can be swapped
// (e.g. Google GenAI, DALL-E, Stable Diffusion HTTP service).
type ImageGenerator interface {
	Generate(ctx context.Context, req ImageRequest) (string, error)
}

type ImageRequest struct {
	ImagePath string `json:"imagePath"`
	Prompt    string `json:"prompt"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}
