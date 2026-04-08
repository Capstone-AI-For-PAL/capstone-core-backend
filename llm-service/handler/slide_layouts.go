package handler

import (
	"fmt"
	"strings"
)

const (
	HeroImageWidth  = 800
	HeroImageHeight = 400
	TwoColWidth     = 450
	TwoColHeight    = 300
	ThreeColWidth   = 300
	ThreeColHeight  = 200
)

func ImageDimensionsForLayout(layout LayoutType) (width, height int) {
	switch layout {
	case LayoutHeroImage:
		return HeroImageWidth, HeroImageHeight
	case LayoutTwoColumn:
		return TwoColWidth, TwoColHeight
	case LayoutThreeColumn:
		return ThreeColWidth, ThreeColHeight
	default:
		return 0, 0
	}
}

func marpFrontmatter() string {
	return `---
marp: true
theme: default
paginate: true
backgroundColor: #f8f9fa
size: 16:9
style: |
  section {
    font-family: 'Segoe UI', system-ui, sans-serif;
    font-size: 24px;
    padding: 50px 70px;
    display: flex;
    flex-direction: column;
    box-sizing: border-box;
  }
  h1 {
    font-size: 38px;
    color: #1a73e8;
    margin: 0 0 16px 0;
    border-bottom: 2px solid #e8eaed;
    padding-bottom: 10px;
  }
  ul {
    margin: 0 0 20px 0;
    color: #3c4043;
    padding-left: 1.5em;
  }
  li {
    margin-bottom: 8px;
  }
  /* Fixed bounds to prevent overflow */
  .hero-container {
    display: flex;
    justify-content: center;
    align-items: center;
    width: 100%;
    height: 360px; /* Fixed hero height */
    margin-top: auto;
    margin-bottom: auto;
  }
  .hero-container img {
    max-width: 100%;
    max-height: 100%;
    object-fit: contain;
    border-radius: 8px;
  }
  .grid-2 {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 40px;
    margin-top: 10px;
  }
  .grid-3 {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 30px;
    margin-top: 10px;
  }
  .card {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
  }
  .card-img-2 {
    width: 100%;
    height: 280px; /* Fixed 2-col image height */
    margin-bottom: 12px;
  }
  .card-img-3 {
    width: 100%;
    height: 200px; /* Fixed 3-col image height */
    margin-bottom: 12px;
  }
  .card img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    border-radius: 8px;
  }
  .card strong {
    font-size: 18px;
    color: #5f6368;
    font-weight: 600;
  }
---
`
}

func renderSlide(slide RenderedSlide) string {
	switch slide.Layout {
	case LayoutHeroImage:
		return renderHeroImage(slide)
	case LayoutTwoColumn:
		return renderTwoColumn(slide)
	case LayoutThreeColumn:
		return renderThreeColumn(slide)
	default:
		return renderTextOnly(slide)
	}
}

func renderKeyPoints(keyPoints []string) string {
	if len(keyPoints) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("\n")
	for _, kp := range keyPoints {
		sb.WriteString(fmt.Sprintf("* %s\n", kp))
	}
	return sb.String()
}

func renderHeroImage(slide RenderedSlide) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n", slide.Title))
	sb.WriteString(renderKeyPoints(slide.KeyPoints))
	if len(slide.Images) > 0 {
		img := slide.Images[0]
		sb.WriteString("\n<div class=\"hero-container\">\n")
		sb.WriteString(fmt.Sprintf("  <img src=\"%s\" alt=\"%s\" />\n", img.URL, img.Caption))
		sb.WriteString("</div>\n")
	}
	return sb.String()
}

func renderTwoColumn(slide RenderedSlide) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n", slide.Title))
	sb.WriteString(renderKeyPoints(slide.KeyPoints))
	sb.WriteString("\n<div class=\"grid-2\">\n")
	for _, img := range slide.Images {
		sb.WriteString("  <div class=\"card\">\n")
		sb.WriteString("    <div class=\"card-img-2\">\n")
		sb.WriteString(fmt.Sprintf("      <img src=\"%s\" alt=\"%s\" />\n", img.URL, img.Caption))
		sb.WriteString("    </div>\n")
		sb.WriteString(fmt.Sprintf("    <span><strong>%s</strong></span>\n", img.Caption))
		sb.WriteString("  </div>\n")
	}
	sb.WriteString("</div>\n")
	return sb.String()
}

func renderThreeColumn(slide RenderedSlide) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n", slide.Title))
	sb.WriteString(renderKeyPoints(slide.KeyPoints))
	sb.WriteString("\n<div class=\"grid-3\">\n")
	for _, img := range slide.Images {
		sb.WriteString("  <div class=\"card\">\n")
		sb.WriteString("    <div class=\"card-img-3\">\n")
		sb.WriteString(fmt.Sprintf("      <img src=\"%s\" alt=\"%s\" />\n", img.URL, img.Caption))
		sb.WriteString("    </div>\n")
		sb.WriteString(fmt.Sprintf("    <span><strong>%s</strong></span>\n", img.Caption))
		sb.WriteString("  </div>\n")
	}
	sb.WriteString("</div>\n")
	return sb.String()
}

func renderTextOnly(slide RenderedSlide) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n", slide.Title))
	sb.WriteString(renderKeyPoints(slide.KeyPoints))
	return sb.String()
}
