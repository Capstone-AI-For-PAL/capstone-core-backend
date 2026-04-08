package handler

import (
	"encoding/json"
	"fmt"
	"strings"
)

func decodeEnrichedSlides(raw string) ([]EnrichedSlide, error) {
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
	var slides []EnrichedSlide
	if err := json.Unmarshal([]byte(jsonArray), &slides); err != nil {
		return nil, fmt.Errorf("invalid enriched slide JSON: %w", err)
	}

	for i := range slides {
		switch slides[i].Layout {
		case LayoutHeroImage, LayoutTwoColumn, LayoutThreeColumn, LayoutTextOnly:
		default:
			switch len(slides[i].Images) {
			case 0:
				slides[i].Layout = LayoutTextOnly
			case 1:
				slides[i].Layout = LayoutHeroImage
			case 2:
				slides[i].Layout = LayoutTwoColumn
			default:
				slides[i].Layout = LayoutThreeColumn
			}
		}
	}

	return slides, nil
}
