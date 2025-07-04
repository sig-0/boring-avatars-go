package avatars

import (
	"fmt"
	"strings"
)

const (
	bauhausSize     = 80
	bauhausElements = 4
)

// bauhausElement holds the per-shape parameters
type bauhausElement struct {
	color                  string
	translateX, translateY float64
	rotate                 int
	square                 bool
}

// buildBauhausElements deterministically derives all element values
func buildBauhausElements(id int, palette Palette) []bauhausElement {
	if len(palette) == 0 {
		palette = DefaultPalette
	}

	elements := make([]bauhausElement, bauhausElements)

	for i := 0; i < bauhausElements; i++ {
		var (
			n    = id + i
			mult = id * (i + 1)
			rng  = bauhausSize/2 - (i + 17)
		)

		elements[i] = bauhausElement{
			color:      palette[n%len(palette)],
			translateX: float64(IDToPoint(mult, rng, 1)),
			translateY: float64(IDToPoint(mult, rng, 2)),
			rotate:     IDToPoint(mult, 360, 0),
			square:     IDToBoolean(id, 2),
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
		center = bauhausSize / 2
	)

	// Start building out the SVG
	var b strings.Builder

	if size > 0 {
		// Custom size
		_, _ = fmt.Fprintf(
			&b,
			`<svg viewBox="0 0 %d %d" fill="none" role="img"`+
				` xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`,
			bauhausSize, bauhausSize,
			size, size,
		)
	} else {
		_, _ = fmt.Fprintf(
			&b,
			`<svg viewBox="0 0 %d %d" fill="none" role="img"`+
				` xmlns="http://www.w3.org/2000/svg">`,
			bauhausSize, bauhausSize,
		)
	}

	// Mask group
	_, _ = fmt.Fprintf(
		&b,
		`<mask id="%s" maskUnits="userSpaceOnUse" x="0" y="0" width="%d" height="%d">`,
		maskID,
		bauhausSize, bauhausSize,
	)

	if square {
		_, _ = fmt.Fprintf(
			&b,
			`<rect width="%d" height="%d" fill="#FFFFFF"/>`,
			bauhausSize, bauhausSize,
		)
	} else {
		_, _ = fmt.Fprintf(
			&b,
			`<rect width="%d" height="%d" rx="%d" fill="#FFFFFF"/>`,
			bauhausSize, bauhausSize,
			bauhausSize*2,
		)
	}

	b.WriteString(`</mask>`)

	// Masked group
	_, _ = fmt.Fprintf(&b, `<g mask="url(#%s)">`, maskID)

	// Background
	_, _ = fmt.Fprintf(
		&b,
		`<rect width="%d" height="%d" fill="%s"/>`,
		bauhausSize, bauhausSize,
		props[0].color,
	)

	// Rotated / translated rectangle
	height := bauhausSize
	if !props[1].square {
		height = bauhausSize / 8 // 10px
	}

	_, _ = fmt.Fprintf(
		&b,
		`<rect x="%d" y="%d" width="%d" height="%d" fill="%s"`+
			` transform="translate(%.2f %.2f) rotate(%d %d %d)"/>`,
		(bauhausSize-60)/2, (bauhausSize-20)/2, // 10, 30
		bauhausSize, height,
		props[1].color,
		props[1].translateX, props[1].translateY,
		props[1].rotate, center, center,
	)

	// Translated circle
	_, _ = fmt.Fprintf(
		&b,
		`<circle cx="%d" cy="%d" r="%d" fill="%s"`+
			` transform="translate(%.2f %.2f)"/>`,
		center, center, bauhausSize/5, // r = 16
		props[2].color,
		props[2].translateX, props[2].translateY,
	)

	// Translated / rotated line
	_, _ = fmt.Fprintf(
		&b,
		`<line x1="0" y1="%d" x2="%d" y2="%d" stroke-width="2" stroke="%s"`+
			` transform="translate(%.2f %.2f) rotate(%d %d %d)"/>`,
		center, bauhausSize, center,
		props[3].color,
		props[3].translateX, props[3].translateY,
		props[3].rotate, center, center,
	)

	// Group closure
	b.WriteString(`</g>`)

	// Final svg closure
	b.WriteString(`</svg>`)

	return b.String()
}
