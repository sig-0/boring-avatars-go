package avatars

import (
	"fmt"
	"strings"
)

const (
	ringSize    = 90
	ringShuffle = 5
)

// buildRingColors creates the 9-entry color list:
//   - first generates 5 shuffled colors based on successive offsets,
//   - then fans them out into the fixed 9-slot pattern used by the SVG paths
func buildRingColors(id int, palette Palette) Palette {
	if len(palette) == 0 {
		palette = DefaultPalette
	}

	shuffle := make([]string, ringShuffle)
	for i := 0; i < ringShuffle; i++ {
		shuffle[i] = palette[(id+i)%len(palette)]
	}

	return []string{
		shuffle[0], // 0
		shuffle[1], // 1
		shuffle[1], // 2
		shuffle[2], // 3
		shuffle[2], // 4
		shuffle[3], // 5
		shuffle[3], // 6
		shuffle[0], // 7
		shuffle[4], // 8
	}
}

// GenerateRing returns the ring-style avatar SVG
func GenerateRing(name string, palette Palette, size int, square bool) string {
	var (
		id     = NameToID(name)
		colors = buildRingColors(id, palette)
		maskID = fmt.Sprintf("mask_ring_%d", id)
	)

	// Start building out the SVG
	var b strings.Builder

	if size > 0 {
		// Custom size
		_, _ = fmt.Fprintf(
			&b,
			`<svg viewBox="0 0 %d %d" fill="none" role="img"`+
				` xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`,
			ringSize, ringSize,
			size, size,
		)
	} else {
		_, _ = fmt.Fprintf(
			&b,
			`<svg viewBox="0 0 %d %d" fill="none" role="img"`+
				` xmlns="http://www.w3.org/2000/svg">`,
			ringSize, ringSize,
		)
	}

	// Mask group
	_, _ = fmt.Fprintf(
		&b,
		`<mask id="%s" maskUnits="userSpaceOnUse" x="0" y="0" width="%d" height="%d">`,
		maskID,
		ringSize, ringSize,
	)

	if square {
		_, _ = fmt.Fprintf(
			&b,
			`<rect width="%d" height="%d" fill="#FFFFFF"/>`,
			ringSize, ringSize,
		)
	} else {
		_, _ = fmt.Fprintf(
			&b,
			`<rect width="%d" height="%d" rx="%d" fill="#FFFFFF"/>`,
			ringSize, ringSize,
			ringSize*2,
		)
	}

	b.WriteString(`</mask>`)

	// Masked group
	_, _ = fmt.Fprintf(&b, `<g mask="url(#%s)">`, maskID)

	// Two halves
	_, _ = fmt.Fprintf(&b, `<path d="M0 0h90v45H0z" fill="%s"/>`, colors[0])
	_, _ = fmt.Fprintf(&b, `<path d="M0 45h90v45H0z" fill="%s"/>`, colors[1])

	// Three concentric rings (2 paths each)
	_, _ = fmt.Fprintf(&b, `<path d="M83 45a38 38 0 00-76 0h76z" fill="%s"/>`, colors[2])
	_, _ = fmt.Fprintf(&b, `<path d="M83 45a38 38 0 01-76 0h76z" fill="%s"/>`, colors[3])

	_, _ = fmt.Fprintf(&b, `<path d="M77 45a32 32 0 10-64 0h64z" fill="%s"/>`, colors[4])
	_, _ = fmt.Fprintf(&b, `<path d="M77 45a32 32 0 11-64 0h64z" fill="%s"/>`, colors[5])

	_, _ = fmt.Fprintf(&b, `<path d="M71 45a26 26 0 00-52 0h52z" fill="%s"/>`, colors[6])
	_, _ = fmt.Fprintf(&b, `<path d="M71 45a26 26 0 01-52 0h52z" fill="%s"/>`, colors[7])

	// Center circle
	_, _ = fmt.Fprintf(&b, `<circle cx="45" cy="45" r="23" fill="%s"/>`, colors[8])

	// Group closure
	b.WriteString(`</g>`)

	// Final svg closure
	b.WriteString(`</svg>`)

	return b.String()
}
