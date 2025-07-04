package avatars

import (
	"fmt"
	"strings"
)

const (
	sunsetSize     = 80
	sunsetElements = 4
)

// buildSunsetColors generates the sunset color palette
func buildSunsetColors(id int, palette Palette) Palette {
	if len(palette) == 0 {
		palette = DefaultPalette
	}

	out := make([]string, sunsetElements)
	for i := 0; i < sunsetElements; i++ {
		out[i] = palette[(id+i)%len(palette)]
	}

	return out
}

// GenerateSunset returns the sunset-style avatar SVG
func GenerateSunset(name string, palette Palette, size int, square bool) string {
	var (
		id     = NameToID(name)
		colors = buildSunsetColors(id, palette)
		maskID = fmt.Sprintf("mask_sunset_%d", id)
	)

	// Start building out the SVG
	var b strings.Builder

	if size > 0 {
		// Custom size
		_, _ = fmt.Fprintf(
			&b,
			`<svg viewBox="0 0 %d %d" fill="none" role="img"`+
				` xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`,
			sunsetSize, sunsetSize,
			size, size,
		)
	} else {
		_, _ = fmt.Fprintf(
			&b,
			`<svg viewBox="0 0 %d %d" fill="none" role="img"`+
				` xmlns="http://www.w3.org/2000/svg">`,
			sunsetSize, sunsetSize,
		)
	}

	// Mask group
	_, _ = fmt.Fprintf(&b,
		`<mask id="%s" maskUnits="userSpaceOnUse" x="0" y="0" width="%d" height="%d">`,
		maskID,
		sunsetSize, sunsetSize,
	)

	if square {
		_, _ = fmt.Fprintf(
			&b,
			`<rect width="%d" height="%d" fill="#FFFFFF"/>`,
			sunsetSize, sunsetSize,
		)
	} else {
		_, _ = fmt.Fprintf(
			&b,
			`<rect width="%d" height="%d" rx="%d" fill="#FFFFFF"/>`,
			sunsetSize, sunsetSize,
			sunsetSize*2,
		)
	}

	b.WriteString(`</mask>`)

	// Masked group
	_, _ = fmt.Fprintf(&b, `<g mask="url(#%s)">`, maskID)

	// Two half-rectangles filled by vertical gradients
	_, _ = fmt.Fprintf(
		&b,
		`<path fill="url(#gradient_paint0_linear_%d)" d="M0 0h80v40H0z"/>`,
		id,
	)

	_, _ = fmt.Fprintf(
		&b,
		`<path fill="url(#gradient_paint1_linear_%d)" d="M0 40h80v40H0z"/>`,
		id,
	)

	b.WriteString(`</g>`)

	// Two linear gradients
	_, _ = fmt.Fprintf(&b,
		`<defs>`+
			`<linearGradient id="gradient_paint0_linear_%d" x1="%d" y1="0" x2="%d" y2="%d" gradientUnits="userSpaceOnUse">`+
			`<stop stop-color="%s"/>`+
			`<stop offset="1" stop-color="%s"/>`+
			`</linearGradient>`+
			`<linearGradient id="gradient_paint1_linear_%d" x1="%d" y1="%d" x2="%d" y2="%d" gradientUnits="userSpaceOnUse">`+
			`<stop stop-color="%s"/>`+
			`<stop offset="1" stop-color="%s"/>`+
			`</linearGradient>`+
			`</defs>`,
		id, sunsetSize/2, sunsetSize/2, sunsetSize/2, // first gradient
		colors[0], colors[1],
		id, sunsetSize/2, sunsetSize/2, sunsetSize/2, sunsetSize, // second gradient
		colors[2], colors[3],
	)

	// Final svg closure
	b.WriteString(`</svg>`)

	return b.String()
}
