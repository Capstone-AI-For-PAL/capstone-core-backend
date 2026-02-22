package genie

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

type AuthRequest struct {
	AppID   string `json:"app_id"`
	Email   string `json:"email"`
	CUNetID string `json:"cunet_id"`
}

type AuthResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
}

type Message struct {
	Role    string        `json:"role"`
	Content []ContentPart `json:"content"`
}

func (m *Message) UnmarshalJSON(data []byte) error {
	type Alias Message
	aux := &struct {
		Content json.RawMessage `json:"content"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if len(aux.Content) > 0 {
		if aux.Content[0] == '"' {
			var str string
			if err := json.Unmarshal(aux.Content, &str); err != nil {
				return err
			}
			m.Content = []ContentPart{{Type: "text", Text: str}}
		} else {
			if err := json.Unmarshal(aux.Content, &m.Content); err != nil {
				return err
			}
		}
	}
	return nil
}

type ContentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

type ImageURL struct {
	URL string `json:"url"`
}

type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

type ChatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

const (
	authURL       = "https://genie.chula.ac.th/api/v1/external/auth"
	completionURL = "https://genie.chula.ac.th/api/v1/external/completions"
)

type tokenData struct {
	token     string
	expiresAt time.Time
}

type Client struct {
	apiKey string
	appID  string
	model  string

	mu     sync.RWMutex
	tokens map[string]tokenData
	sf     singleflight.Group
}

func NewClient() *Client {
	return &Client{
		apiKey: os.Getenv("GENIE_API_KEY"),
		appID:  os.Getenv("GENIE_APP_ID"),
		model:  os.Getenv("GENIE_MODEL"),
		tokens: make(map[string]tokenData),
	}
}

func (c *Client) getToken(email, cunetId string) (string, error) {
	key := fmt.Sprintf("%s:%s", email, cunetId)

	c.mu.RLock()
	if data, ok := c.tokens[key]; ok {
		if time.Now().Before(data.expiresAt) {
			c.mu.RUnlock()
			return data.token, nil
		}
	}
	c.mu.RUnlock()

	v, err, _ := c.sf.Do(key, func() (interface{}, error) {
		payload, err := json.Marshal(AuthRequest{
			AppID:   c.appID,
			Email:   email,
			CUNetID: cunetId,
		})
		if err != nil {
			return "", fmt.Errorf("failed to marshal auth request: %v", err)
		}

		req, err := http.NewRequest("POST", authURL, bytes.NewBuffer(payload))
		if err != nil {
			return "", fmt.Errorf("failed to create auth request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("api-key", c.apiKey)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", errors.New("GENIE auth failed")
		}

		var res AuthResponse
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return "", fmt.Errorf("failed to decode auth response: %v", err)
		}

		c.mu.Lock()
		c.tokens[key] = tokenData{
			token:     res.Token,
			expiresAt: time.Now().Add(4 * time.Hour),
		}
		c.mu.Unlock()

		return res.Token, nil
	})

	if err != nil {
		return "", err
	}

	return v.(string), nil
}

// https://docs.litellm.ai/docs/completion/input

func (c *Client) Chat(messages []Message, email, cunetId string) (string, error) {
	token, err := c.getToken(email, cunetId)
	if err != nil {
		return "", err
	}

	reqBody := ChatCompletionRequest{
		Model:    c.model,
		Messages: messages,
		Stream:   false,
	}

	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", completionURL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("GENIE completion failed: %s - %s", resp.Status, string(bodyBytes))
	}

	var res ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	if len(res.Choices) == 0 {
		return "", errors.New("empty response")
	}

	return res.Choices[0].Message.Content, nil
}
