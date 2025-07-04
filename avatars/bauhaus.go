package avatars

import (
	"fmt"
	"strings"
)

const (
	BauhausSize     = 80
	bauhausElements = 4
)

// BauhausElement holds the per-shape parameters
type BauhausElement struct {
	Color                  string
	TranslateX, TranslateY float64
	Rotate                 int
	Square                 bool
}

// buildBauhausElements deterministically derives all element values
func buildBauhausElements(id int, palette Palette) []BauhausElement {
	if len(palette) == 0 {
		palette = DefaultPalette
	}

	elements := make([]BauhausElement, bauhausElements)

	for i := 0; i < bauhausElements; i++ {
		var (
			n    = id + i
			mult = id * (i + 1)
			rng  = BauhausSize/2 - (i + 17)
		)

		elements[i] = BauhausElement{
			Color:      palette[n%len(palette)],
			TranslateX: float64(IDToPoint(mult, rng, 1)),
			TranslateY: float64(IDToPoint(mult, rng, 2)),
			Rotate:     IDToPoint(mult, 360, 0),
			Square:     IDToBoolean(id, 2),
		}
	}

	return elements
}

// GenerateBauhaus returns a bauhaus-style avatar SVG
func GenerateBauhaus(name string, palette Palette, size int, square bool) string {
	var (
		id     = NameToID(name)
		props  = buildBauhausElements(id, palette)
		maskID = fmt.Sprintf("mask_bauhaus_%d", id)
		center = BauhausSize / 2
	)

	// Start building out the SVG
	var b strings.Builder

	if size > 0 {
		// Custom size
		_, _ = fmt.Fprintf(
			&b,
			`<svg viewBox="0 0 %d %d" fill="none" role="img"`+
				` xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`,
			BauhausSize, BauhausSize,
			size, size,
		)
	} else {
		_, _ = fmt.Fprintf(
			&b,
			`<svg viewBox="0 0 %d %d" fill="none" role="img"`+
				` xmlns="http://www.w3.org/2000/svg">`,
			BauhausSize, BauhausSize,
		)
	}

	// Mask group
	_, _ = fmt.Fprintf(
		&b,
		`<mask id="%s" maskUnits="userSpaceOnUse" x="0" y="0" width="%d" height="%d">`,
		maskID,
		BauhausSize, BauhausSize,
	)

	if square {
		_, _ = fmt.Fprintf(
			&b,
			`<rect width="%d" height="%d" fill="#FFFFFF"/>`,
			BauhausSize, BauhausSize,
		)
	} else {
		_, _ = fmt.Fprintf(
			&b,
			`<rect width="%d" height="%d" rx="%d" fill="#FFFFFF"/>`,
			BauhausSize, BauhausSize,
			BauhausSize*2,
		)
	}

	b.WriteString(`</mask>`)

	// Masked group
	_, _ = fmt.Fprintf(&b, `<g mask="url(#%s)">`, maskID)

	// Background
	_, _ = fmt.Fprintf(
		&b,
		`<rect width="%d" height="%d" fill="%s"/>`,
		BauhausSize, BauhausSize,
		props[0].Color,
	)

	// Rotated / translated rectangle
	height := BauhausSize
	if !props[1].Square {
		height = BauhausSize / 8 // 10px
	}

	_, _ = fmt.Fprintf(
		&b,
		`<rect x="%d" y="%d" width="%d" height="%d" fill="%s"`+
			` transform="translate(%.2f %.2f) rotate(%d %d %d)"/>`,
		(BauhausSize-60)/2, (BauhausSize-20)/2, // 10, 30
		BauhausSize, height,
		props[1].Color,
		props[1].TranslateX, props[1].TranslateY,
		props[1].Rotate, center, center,
	)

	// Translated circle
	_, _ = fmt.Fprintf(
		&b,
		`<circle cx="%d" cy="%d" r="%d" fill="%s"`+
			` transform="translate(%.2f %.2f)"/>`,
		center, center, BauhausSize/5, // r = 16
		props[2].Color,
		props[2].TranslateX, props[2].TranslateY,
	)

	// Translated / rotated line
	_, _ = fmt.Fprintf(
		&b,
		`<line x1="0" y1="%d" x2="%d" y2="%d" stroke-width="2" stroke="%s"`+
			` transform="translate(%.2f %.2f) rotate(%d %d %d)"/>`,
		center, BauhausSize, center,
		props[3].Color,
		props[3].TranslateX, props[3].TranslateY,
		props[3].Rotate, center, center,
	)

	// Group closure
	b.WriteString(`</g>`)

	// Final svg closure
	b.WriteString(`</svg>`)

	return b.String()
}
