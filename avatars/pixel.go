package avatars

import (
	"fmt"
	"strings"
)

const (
	PixelSize     = 80
	pixelElements = 64 // 8x8 grid
	pixelStep     = 10 // each "pixel" is 10 px
)

// buildPixelColors returns the 64 palette entries used in the grid
func buildPixelColors(id int, palette Palette) []string {
	if len(palette) == 0 {
		palette = DefaultPalette
	}

	out := make([]string, pixelElements)

	for i := 0; i < pixelElements; i++ {
		n := id % (i + 1)
		out[i] = palette[n%len(palette)]
	}

	return out
}

// GeneratePixel returns an 8x8 pixel-art avatar SVG
func GeneratePixel(name string, palette Palette, size int, square bool) string {
	var (
		id     = NameToID(name)
		colors = buildPixelColors(id, palette)
		maskID = fmt.Sprintf("mask_pixel_%d", id)
	)

	// Start building out the SVG
	var b strings.Builder

	if size > 0 {
		// Custom size
		_, _ = fmt.Fprintf(
			&b,
			`<svg viewBox="0 0 %d %d" fill="none" role="img"`+
				` xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`,
			PixelSize, PixelSize,
			size, size,
		)
	} else {
		_, _ = fmt.Fprintf(
			&b,
			`<svg viewBox="0 0 %d %d" fill="none" role="img"`+
				` xmlns="http://www.w3.org/2000/svg">`,
			PixelSize, PixelSize,
		)
	}

	// Mask group
	_, _ = fmt.Fprintf(
		&b,
		`<mask id="%s" mask-type="alpha" maskUnits="userSpaceOnUse" x="0" y="0" width="%d" height="%d">`,
		maskID,
		PixelSize, PixelSize,
	)

	if square {
		_, _ = fmt.Fprintf(
			&b,
			`<rect width="%d" height="%d" fill="#FFFFFF"/>`,
			PixelSize, PixelSize,
		)
	} else {
		_, _ = fmt.Fprintf(
			&b,
			`<rect width="%d" height="%d" rx="%d" fill="#FFFFFF"/>`,
			PixelSize, PixelSize,
			PixelSize*2,
		)
	}

	b.WriteString(`</mask>`)

	// Masked grid
	_, _ = fmt.Fprintf(&b, `<g mask="url(#%s)">`, maskID)

	var (
		cols = []int{0, 20, 40, 60, 10, 30, 50, 70} // even columns first, then odd
		idx  = 0
	)

	writePixel := func(x, y int, fill string) {
		switch {
		case x == 0 && y == 0:
			_, _ = fmt.Fprintf(&b, `<rect width="10" height="10" fill="%s"/>`, fill)
		case x == 0:
			_, _ = fmt.Fprintf(&b, `<rect y="%d" width="10" height="10" fill="%s"/>`, y, fill)
		case y == 0:
			_, _ = fmt.Fprintf(&b, `<rect x="%d" width="10" height="10" fill="%s"/>`, x, fill)
		default:
			_, _ = fmt.Fprintf(&b, `<rect x="%d" y="%d" width="10" height="10" fill="%s"/>`, x, y, fill)
		}
	}

	// Row 0 (y = 0)
	for _, x := range cols {
		writePixel(x, 0, colors[idx])

		idx++
	}

	// The remaining rows, column by column (y = 10..70)
	for _, x := range cols {
		for y := 10; y < PixelSize; y += 10 {
			writePixel(x, y, colors[idx])

			idx++
		}
	}

	// Group closure
	b.WriteString(`</g>`)

	// Final svg closure
	b.WriteString(`</svg>`)

	return b.String()
}
