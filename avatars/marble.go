package avatars

import (
	"fmt"
	"strings"
)

const (
	MarbleSize     = 80
	marbleElements = 3
)

// MarbleElement stores the params
type MarbleElement struct {
	Color                  string
	TranslateX, TranslateY float64
	Scale                  float64
	Rotate                 int
}

// buildMarbleElements deterministically derives all element values
func buildMarbleElements(id int, palette Palette) []MarbleElement {
	if len(palette) == 0 {
		palette = DefaultPalette
	}

	elements := make([]MarbleElement, marbleElements)

	for i := 0; i < marbleElements; i++ {
		n := id + i
		m := id * (i + 1)

		elements[i] = MarbleElement{
			Color:      palette[n%len(palette)],
			TranslateX: float64(IDToPoint(m, MarbleSize/10, 1)),
			TranslateY: float64(IDToPoint(m, MarbleSize/10, 2)),
			Scale:      1.2 + float64(IDToPoint(m, MarbleSize/20, 0))/10.0,
			Rotate:     IDToPoint(m, 360, 1),
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
		center   = MarbleSize / 2
	)

	// Start building out the SVG
	var b strings.Builder

	if size > 0 {
		// Custom size
		_, _ = fmt.Fprintf(
			&b,
			`<svg viewBox="0 0 %d %d" fill="none" role="img"`+
				` xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`,
			MarbleSize, MarbleSize,
			size, size,
		)
	} else {
		_, _ = fmt.Fprintf(
			&b,
			`<svg viewBox="0 0 %d %d" fill="none" role="img"`+
				` xmlns="http://www.w3.org/2000/svg">`,
			MarbleSize, MarbleSize,
		)
	}

	// Mask group
	_, _ = fmt.Fprintf(
		&b,
		`<mask id="%s" maskUnits="userSpaceOnUse" x="0" y="0" width="%d" height="%d">`,
		maskID,
		MarbleSize, MarbleSize,
	)

	if square {
		_, _ = fmt.Fprintf(
			&b,
			`<rect width="%d" height="%d" fill="#FFFFFF"/>`,
			MarbleSize, MarbleSize,
		)
	} else {
		_, _ = fmt.Fprintf(
			&b,
			`<rect width="%d" height="%d" rx="%d" fill="#FFFFFF"/>`,
			MarbleSize, MarbleSize,
			MarbleSize*2,
		)
	}

	b.WriteString(`</mask>`)

	// Masked group
	_, _ = fmt.Fprintf(&b, `<g mask="url(#%s)">`, maskID)

	// Background
	_, _ = fmt.Fprintf(
		&b,
		`<rect width="%d" height="%d" fill="%s"/>`,
		MarbleSize, MarbleSize,
		props[0].Color,
	)

	// First path
	_, _ = fmt.Fprintf(
		&b,
		`<path filter="url(#%s)" d="M32.414 59.35L50.376 70.5H72.5v-71H33.728L26.5 13.381l19.057 27.08L32.414 59.35z"`+
			` fill="%s"`+
			` transform="translate(%.2f %.2f) rotate(%d %d %d) scale(%.2f)"/>`,
		filterID,
		props[1].Color,
		props[1].TranslateX, props[1].TranslateY,
		props[1].Rotate, center, center,
		props[2].Scale,
	)

	// Second path
	_, _ = fmt.Fprintf(
		&b,
		`<path filter="url(#%s)" style="mix-blend-mode:overlay"`+
			` d="M22.216 24L0 46.75l14.108 38.129L78 86l-3.081-59.276-22.378 4.005 12.972 20.186-23.35 27.395L22.215 24z"`+
			` fill="%s"`+
			` transform="translate(%.2f %.2f) rotate(%d %d %d) scale(%.2f)"/>`,
		filterID,
		props[2].Color,
		props[2].TranslateX, props[2].TranslateY,
		props[2].Rotate, center, center,
		props[2].Scale,
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
