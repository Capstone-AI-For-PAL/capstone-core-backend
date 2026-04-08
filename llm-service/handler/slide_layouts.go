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
	return "---\n" +
		"marp: true\n" +
		"theme: default\n" +
		"paginate: true\n" +
		"backgroundColor: #fff\n" +
		"math: mathjax\n" +
		"size: 16:9\n" +
		"style: |\n" +
		"  section {\n" +
		"    font-size: 25px;\n" +
		"    padding: 40px;\n" +
		"    justify-content: center;\n" +
		"  }\n" +
		"  h1 {\n" +
		"    font-size: 40px;\n" +
		"    color: #0288d1;\n" +
		"  }\n" +
		"  h2 {\n" +
		"    font-size: 35px;\n" +
		"    color: #333;\n" +
		"  }\n" +
		"  .grid-2 {\n" +
		"    display: grid;\n" +
		"    grid-template-columns: 1fr 1fr;\n" +
		"    gap: 30px;\n" +
		"    text-align: center;\n" +
		"  }\n" +
		"  .grid-3 {\n" +
		"    display: grid;\n" +
		"    grid-template-columns: 1fr 1fr 1fr;\n" +
		"    gap: 15px;\n" +
		"    text-align: center;\n" +
		"    font-size: 0.7em;\n" +
		"  }\n" +
		"  img {\n" +
		"    border-radius: 8px;\n" +
		"    object-fit: cover;\n" +
		"  }\n" +
		"---\n"
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
		sb.WriteString("\n<div style=\"text-align: center; margin-top: 20px;\">\n\n")
		sb.WriteString(fmt.Sprintf("![%s](%s)\n\n", img.Caption, img.URL))
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
		sb.WriteString("  <div>\n")
		sb.WriteString(fmt.Sprintf("    <img src=\"%s\" alt=\"%s\" />\n", img.URL, img.Caption))
		sb.WriteString(fmt.Sprintf("    <p><strong>%s</strong></p>\n", img.Caption))
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
		sb.WriteString("  <div>\n")
		sb.WriteString(fmt.Sprintf("    <img src=\"%s\" />\n", img.URL))
		sb.WriteString(fmt.Sprintf("    <br><strong>%s</strong>\n", img.Caption))
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
