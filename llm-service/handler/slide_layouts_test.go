package handler

import (
	"fmt"
	"strings"
	"testing"
)

// Picsum placeholder URLs – seeded so results are deterministic.
const (
	picsumHero = "https://picsum.photos/seed/hero/800/400"
	picsumColA = "https://picsum.photos/seed/colA/450/300"
	picsumColB = "https://picsum.photos/seed/colB/450/300"
	picsumColC = "https://picsum.photos/seed/colC/300/200"
	picsumColD = "https://picsum.photos/seed/colD/300/200"
	picsumColE = "https://picsum.photos/seed/colE/300/200"
)

func marpDocument(body string) string {
	return marpFrontmatter() + "---\n\n" + body
}

// TestRenderHeroImage prints a complete Marp document for the hero_image layout.
func TestRenderHeroImage(t *testing.T) {
	slide := RenderedSlide{
		Title:  "\"It Works on My Machine!\" Syndrome",
		Layout: LayoutHeroImage,
		KeyPoints: []string{
			"Software runs locally, fails elsewhere.",
			"Common developer frustration.",
			"Highlights environment inconsistency.",
		},
		Images: []RenderedImage{
			{
				URL:     picsumHero,
				Caption: "Contrasting development and production environments with subtle differences in software versions of key components like Python or Node.js.",
				Width:   HeroImageWidth,
				Height:  HeroImageHeight,
			},
		},
	}

	output := marpDocument(renderHeroImage(slide))
	fmt.Println("=== TestRenderHeroImage ===")
	fmt.Println(output)
}

// TestRenderHeroImageNoCaption prints a hero_image slide where the image has no caption.
func TestRenderHeroImageNoCaption(t *testing.T) {
	slide := RenderedSlide{
		Title:  "The Need for Consistent Environments",
		Layout: LayoutHeroImage,
		KeyPoints: []string{
			"Demands a universal packaging solution.",
			"Ensures consistent runtime behavior.",
			"Foundation for reliable deployment.",
		},
		Images: []RenderedImage{
			{
				URL:     picsumHero,
				Caption: "",
				Width:   HeroImageWidth,
				Height:  HeroImageHeight,
			},
		},
	}

	output := marpDocument(renderHeroImage(slide))
	fmt.Println("=== TestRenderHeroImageNoCaption ===")
	fmt.Println(output)
}

// TestRenderTwoColumn prints a complete Marp document for the two_column layout.
func TestRenderTwoColumn(t *testing.T) {
	slide := RenderedSlide{
		Title:  "The Root Cause: Environmental Drift",
		Layout: LayoutTwoColumn,
		KeyPoints: []string{
			"Missing libraries or dependencies.",
			"Different operating system versions or configurations.",
			"Inconsistent tool versions (e.g., Python, Node.js).",
		},
		Images: []RenderedImage{
			{
				URL:     picsumColA,
				Caption: "Developer machine environment",
				Width:   TwoColWidth,
				Height:  TwoColHeight,
			},
			{
				URL:     picsumColB,
				Caption: "Production server environment",
				Width:   TwoColWidth,
				Height:  TwoColHeight,
			},
		},
	}

	output := marpDocument(renderTwoColumn(slide))
	fmt.Println("=== TestRenderTwoColumn ===")
	fmt.Println(output)
}

// TestRenderThreeColumn prints a complete Marp document for the three_column layout.
func TestRenderThreeColumn(t *testing.T) {
	slide := RenderedSlide{
		Title:  "The Dependency Management Mess",
		Layout: LayoutThreeColumn,
		KeyPoints: []string{
			"Tracking each version is crucial.",
			"Conflicts arise across service stacks.",
			"Updates can break other systems.",
		},
		Images: []RenderedImage{
			{
				URL:     picsumColC,
				Caption: "Service A dependencies",
				Width:   ThreeColWidth,
				Height:  ThreeColHeight,
			},
			{
				URL:     picsumColD,
				Caption: "Service B dependencies",
				Width:   ThreeColWidth,
				Height:  ThreeColHeight,
			},
			{
				URL:     picsumColE,
				Caption: "Service C dependencies",
				Width:   ThreeColWidth,
				Height:  ThreeColHeight,
			},
		},
	}

	output := marpDocument(renderThreeColumn(slide))
	fmt.Println("=== TestRenderThreeColumn ===")
	fmt.Println(output)
}

// TestRenderTextOnly prints a complete Marp document for the text_only layout.
func TestRenderTextOnly(t *testing.T) {
	slide := RenderedSlide{
		Title:  "The Cost of Inconsistency",
		Layout: LayoutTextOnly,
		KeyPoints: []string{
			"Delayed deployments and release cycles.",
			"Time-consuming debugging efforts.",
			"Increased operational overhead.",
			"Reduced team productivity and morale.",
		},
	}

	output := marpDocument(renderTextOnly(slide))
	fmt.Println("=== TestRenderTextOnly ===")
	fmt.Println(output)
}

// TestRenderAllLayouts prints a single multi-slide Marp document containing
// one slide for every layout type — useful for a quick full-deck preview.
func TestRenderAllLayouts(t *testing.T) {
	slides := []RenderedSlide{
		{
			Title:  "\"It Works on My Machine!\" Syndrome",
			Layout: LayoutHeroImage,
			KeyPoints: []string{
				"Software runs locally, fails elsewhere.",
				"Common developer frustration.",
				"Highlights environment inconsistency.",
			},
			Images: []RenderedImage{
				{URL: picsumHero, Caption: "A developer frustrated by environment differences.", Width: HeroImageWidth, Height: HeroImageHeight},
			},
		},
		{
			Title:  "The Root Cause: Environmental Drift",
			Layout: LayoutTwoColumn,
			KeyPoints: []string{
				"Missing libraries or dependencies.",
				"Inconsistent tool versions.",
			},
			Images: []RenderedImage{
				{URL: picsumColA, Caption: "Developer machine", Width: TwoColWidth, Height: TwoColHeight},
				{URL: picsumColB, Caption: "Production server", Width: TwoColWidth, Height: TwoColHeight},
			},
		},
		{
			Title:  "The Dependency Management Mess",
			Layout: LayoutThreeColumn,
			KeyPoints: []string{
				"Tracking each version is crucial.",
				"Conflicts arise across service stacks.",
			},
			Images: []RenderedImage{
				{URL: picsumColC, Caption: "Service A", Width: ThreeColWidth, Height: ThreeColHeight},
				{URL: picsumColD, Caption: "Service B", Width: ThreeColWidth, Height: ThreeColHeight},
				{URL: picsumColE, Caption: "Service C", Width: ThreeColWidth, Height: ThreeColHeight},
			},
		},
		{
			Title:  "The Cost of Inconsistency",
			Layout: LayoutTextOnly,
			KeyPoints: []string{
				"Delayed deployments and release cycles.",
				"Time-consuming debugging efforts.",
				"Reduced team productivity and morale.",
			},
		},
	}

	var body strings.Builder
	for i, slide := range slides {
		if i > 0 {
			body.WriteString("\n---\n\n")
		}
		body.WriteString(renderSlide(slide))
	}

	output := marpDocument(body.String())
	fmt.Println("=== TestRenderAllLayouts ===")
	fmt.Println(output)
}
