package imagegen

import (
	"context"
	"net/http"
)

type ImageGenerator interface {
	Generate(ctx context.Context, req ImageRequest) (string, error)
}

type ImageRequest struct {
	ImagePath string `json:"imagePath"`
	Prompt    string `json:"prompt"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

type imageResponse struct {
	URL string `json:"url"`
}

type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

func (c *HTTPClient) Generate(ctx context.Context, req ImageRequest) (string, error) {
	// body, err := json.Marshal(req)
	// if err != nil {
	// 	return "", fmt.Errorf("marshal image request: %w", err)
	// }

	// httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(body))
	// if err != nil {
	// 	return "", fmt.Errorf("create image request: %w", err)
	// }
	// httpReq.Header.Set("Content-Type", "application/json")

	// resp, err := c.httpClient.Do(httpReq)
	// if err != nil {
	// 	return "", fmt.Errorf("image generation request: %w", err)
	// }
	// defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	respBody, _ := io.ReadAll(resp.Body)
	// 	return "", fmt.Errorf("image generation failed (%d): %s", resp.StatusCode, string(respBody))
	// }

	// var result imageResponse
	// if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
	// 	return "", fmt.Errorf("decode image response: %w", err)
	// }

	// return result.URL, nil
	return "https://picsum.photos/450/300?image=577", nil
}
