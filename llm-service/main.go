package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	genie "capstone-llm-service/llm"

	"github.com/joho/godotenv"
)

type ChatRequest struct {
	Messages []genie.Message `json:"messages"`
	CunetId  string          `json:"cunet_id"`
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not loaded, using system env")
	}
}

func validateEnv() {
	required := []string{
		"GENIE_API_KEY",
		"GENIE_APP_ID",
		"GENIE_MODEL",
	}

	for _, key := range required {
		if os.Getenv(key) == "" {
			log.Fatalf("Missing required env: %s", key)
		}
	}
}

func main() {
	// Load .env file
	loadEnv()
	validateEnv()

	client := genie.NewClient()

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", 405)
			return
		}

		var messages []genie.Message
		var email, cunetId string
		contentType := r.Header.Get("Content-Type")

		if strings.Contains(contentType, "application/json") {
			var req ChatRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			messages = req.Messages
			email = req.CunetId + "@student.chula.ac.th"
			cunetId = req.CunetId
		} else if strings.Contains(contentType, "multipart/form-data") {
			if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB limit
				log.Printf("Error parsing multipart form: %v", err)
				http.Error(w, "file too large", 400)
				return
			}

			email = r.FormValue("cunet_id") + "@student.chula.ac.th"
			cunetId = r.FormValue("cunet_id")
			msgStr := r.FormValue("messages")
			if msgStr != "" {
				if err := json.Unmarshal([]byte(msgStr), &messages); err != nil {
					log.Printf("Error unmarshaling messages: %v", err)
					http.Error(w, "invalid messages json", 400)
					return
				}
			}

			appendPart := func(part genie.ContentPart) {
				if len(messages) == 0 {
					messages = append(messages, genie.Message{Role: "user", Content: []genie.ContentPart{part}})
					return
				}
				lastIdx := len(messages) - 1
				lastMsg := &messages[lastIdx]

				if lastMsg.Role != "user" {
					messages = append(messages, genie.Message{Role: "user", Content: []genie.ContentPart{part}})
					return
				}

				lastMsg.Content = append(lastMsg.Content, part)
			}

			// Process uploaded files
			for _, key := range []string{"image", "file"} {
				file, _, err := r.FormFile(key)
				if err != nil {
					if err != http.ErrMissingFile {
						log.Printf("Error retrieving file for key %s: %v", key, err)
					}
					continue
				}

				func() {
					log.Printf("Processing file upload for key: %s", key)
					defer file.Close()

					buf := new(bytes.Buffer)
					if _, err := io.Copy(buf, file); err != nil {
						log.Printf("Error reading file content: %v", err)
						return
					}

					fileType := http.DetectContentType(buf.Bytes())
					log.Printf("Detected file type: %s", fileType)

					// Encode to base64
					encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
					dataURL := fmt.Sprintf("data:%s;base64,%s", fileType, encoded)

					part := genie.ContentPart{
						Type: "image_url",
						ImageURL: &genie.ImageURL{
							URL: dataURL,
						},
					}
					appendPart(part)
				}()
			}
		} else {
			http.Error(w, "unsupported content type", 400)
			return
		}

		if cunetId == "" {
			http.Error(w, "missing required fields: cunet_id", 400)
			return
		}

		res, err := client.Chat(messages, email, cunetId)
		if err != nil {
			log.Printf("Chat error: %v", err)
			http.Error(w, err.Error(), 500)
			return
		}

		json.NewEncoder(w).Encode(res)
	})

	log.Println("🚀 server started :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
