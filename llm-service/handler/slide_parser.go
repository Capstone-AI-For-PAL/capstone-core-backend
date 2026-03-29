package handler

import (
	"encoding/json"
	"fmt"
	"strings"
)

func decodeOrchestratedSlides(raw string) ([]OrchestratedSlide, error) {
	cleaned := strings.TrimSpace(raw)

	if strings.HasPrefix(cleaned, "```") {
		lines := strings.Split(cleaned, "\n")
		if len(lines) >= 3 {
			lines = lines[1:]
			if strings.TrimSpace(lines[len(lines)-1]) == "```" {
				lines = lines[:len(lines)-1]
			}
			cleaned = strings.TrimSpace(strings.Join(lines, "\n"))
		}
	}

	start := strings.Index(cleaned, "[")
	end := strings.LastIndex(cleaned, "]")
	if start == -1 || end == -1 || end < start {
		return nil, fmt.Errorf("model output does not contain a JSON array")
	}

	jsonArray := cleaned[start : end+1]
	var slides []OrchestratedSlide
	if err := json.Unmarshal([]byte(jsonArray), &slides); err != nil {
		return nil, fmt.Errorf("invalid slide JSON: %w", err)
	}

	return slides, nil
}
