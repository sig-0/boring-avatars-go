package avatars

import (
	"fmt"
	"strings"
)

const (
	marbleSize     = 80
	marbleElements = 3
)

// marbleElement stores the params
type marbleElement struct {
	color                  string
	translateX, translateY float64
	scale                  float64
	rotate                 int
}

// buildMarbleElements deterministically derives all element values
func buildMarbleElements(id int, palette Palette) []marbleElement {
	if len(palette) == 0 {
		palette = DefaultPalette
	}

	elements := make([]marbleElement, marbleElements)

	for i := 0; i < marbleElements; i++ {
		n := id + i
		m := id * (i + 1)

		elements[i] = marbleElement{
			color:      palette[n%len(palette)],
			translateX: float64(IDToPoint(m, marbleSize/10, 1)),
			translateY: float64(IDToPoint(m, marbleSize/10, 2)),
			scale:      1.2 + float64(IDToPoint(m, marbleSize/20, 0))/10.0,
			rotate:     IDToPoint(m, 360, 1),
		}
	}

	return elements
}

// GenerateMarble returns a marble-style avatar SVG
func GenerateMarble(name string, palette Palette, size int, square bool) string {
	var (
		id       = NameToID(name)
		props    = buildMarbleElements(id, palette)
		maskID   = fmt.Sprintf("mask_marble_%d", id)
		filterID = "filter_" + maskID
		center   = marbleSize / 2
	)

	// Start building out the SVG
	var b strings.Builder

	if size > 0 {
		// Custom size
		_, _ = fmt.Fprintf(
			&b,
			`<svg viewBox="0 0 %d %d" fill="none" role="img"`+
				` xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`,
			marbleSize, marbleSize,
			size, size,
		)
	} else {
		_, _ = fmt.Fprintf(
			&b,
			`<svg viewBox="0 0 %d %d" fill="none" role="img"`+
				` xmlns="http://www.w3.org/2000/svg">`,
			marbleSize, marbleSize,
		)
	}

	// Mask group
	_, _ = fmt.Fprintf(
		&b,
		`<mask id="%s" maskUnits="userSpaceOnUse" x="0" y="0" width="%d" height="%d">`,
		maskID,
		marbleSize, marbleSize,
	)

	if square {
		_, _ = fmt.Fprintf(
			&b,
			`<rect width="%d" height="%d" fill="#FFFFFF"/>`,
			marbleSize, marbleSize,
		)
	} else {
		_, _ = fmt.Fprintf(
			&b,
			`<rect width="%d" height="%d" rx="%d" fill="#FFFFFF"/>`,
			marbleSize, marbleSize,
			marbleSize*2,
		)
	}

	b.WriteString(`</mask>`)

	// Masked group
	_, _ = fmt.Fprintf(&b, `<g mask="url(#%s)">`, maskID)

	// Background
	_, _ = fmt.Fprintf(
		&b,
		`<rect width="%d" height="%d" fill="%s"/>`,
		marbleSize, marbleSize,
		props[0].color,
	)

	// First path
	_, _ = fmt.Fprintf(
		&b,
		`<path filter="url(#%s)" d="M32.414 59.35L50.376 70.5H72.5v-71H33.728L26.5 13.381l19.057 27.08L32.414 59.35z"`+
			` fill="%s"`+
			` transform="translate(%.2f %.2f) rotate(%d %d %d) scale(%.2f)"/>`,
		filterID,
		props[1].color,
		props[1].translateX, props[1].translateY,
		props[1].rotate, center, center,
		props[2].scale,
	)

	// Second path
	_, _ = fmt.Fprintf(
		&b,
		`<path filter="url(#%s)" style="mix-blend-mode:overlay"`+
			` d="M22.216 24L0 46.75l14.108 38.129L78 86l-3.081-59.276-22.378 4.005 12.972 20.186-23.35 27.395L22.215 24z"`+
			` fill="%s"`+
			` transform="translate(%.2f %.2f) rotate(%d %d %d) scale(%.2f)"/>`,
		filterID,
		props[2].color,
		props[2].translateX, props[2].translateY,
		props[2].rotate, center, center,
		props[2].scale,
	)

	// Close group
	b.WriteString(`</g>`)

	// Blur filter
	_, _ = fmt.Fprintf(
		&b,
		`<defs>`+
			`<filter id="%s" filterUnits="userSpaceOnUse" color-interpolation-filters="sRGB">`+
			`<feFlood flood-opacity="0" result="BackgroundImageFix"/>`+
			`<feBlend in="SourceGraphic" in2="BackgroundImageFix" result="shape"/>`+
			`<feGaussianBlur stdDeviation="7" result="effect1_foregroundBlur"/>`+
			`</filter>`+
			`</defs>`,
		filterID,
	)

	// Final svg closure
	b.WriteString(`</svg>`)

	return b.String()
}
