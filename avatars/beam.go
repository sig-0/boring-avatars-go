package avatars

import (
	"fmt"
	"strings"
)

const beamSize = 36

type colors struct {
	wrapper    string // color of the rotated square / circle
	face       string // black or white for eyes and mouth (contrast to Wrapper)
	background string // full-canvas background
}

type wrapper struct {
	translateX, translateY float64 // translate
	rotate                 int     // degrees
	scale                  float64 // 1.0 - ~1.3
	circle                 bool    // true -> circle; false -> square
}

type face struct {
	translateX, translateY float64 // translate
	rotate                 int     // degrees
	eyeSpread              int     // -5..+5 px
	mouthSpread            int     // -3..+3 px
	mouthOpen              bool    // open (curved line) or closed (arc)
}

type params struct {
	colors  colors
	face    face
	wrapper wrapper
}

// buildBeamParams deterministically derives all beam avatar parameters from an ID
func buildBeamParams(id int, palette Palette) params {
	if len(palette) == 0 {
		palette = DefaultPalette
	}

	translatePoint := func(id int, place int) float64 {
		point := IDToPoint(id, 10, place)

		if point < 5 {
			point += beamSize / 9
		}

		return float64(point)
	}

	var (
		wrapperTranslateX = translatePoint(id, 1)
		wrapperTranslateY = translatePoint(id, 2)

		wrapperColor = palette[id%len(palette)]
	)

	// Prepare the params
	p := params{
		colors: colors{
			wrapper:    wrapperColor,
			face:       Contrast(wrapperColor),
			background: palette[(id+13)%len(palette)],
		},
		wrapper: wrapper{
			translateX: wrapperTranslateX,
			translateY: wrapperTranslateY,
			rotate:     IDToPoint(id, 360, 0),
			scale:      1 + float64(IDToPoint(id, beamSize/12, 0))/10,
			circle:     IDToBoolean(id, 1),
		},
	}

	// Prepare the face translation.
	// It's tied to the wrapper offset
	faceTranslateX := wrapperTranslateX / 2
	if wrapperTranslateX <= beamSize/6 {
		faceTranslateX = float64(IDToPoint(id, 8, 1))
	}

	faceTranslateY := wrapperTranslateY / 2
	if wrapperTranslateY <= beamSize/6 {
		faceTranslateY = float64(IDToPoint(id, 7, 2))
	}

	p.face = face{
		translateX:  faceTranslateX,
		translateY:  faceTranslateY,
		rotate:      IDToPoint(id, 10, 3),
		eyeSpread:   IDToPoint(id, 5, 0),
		mouthSpread: IDToPoint(id, 3, 0),
		mouthOpen:   IDToBoolean(id, 2),
	}

	return p
}

// GenerateBeam returns a beam-style avatar SVG
func GenerateBeam(name string, palette Palette, size int, square bool) string {
	var (
		id     = NameToID(name)
		p      = buildBeamParams(id, palette)
		maskID = fmt.Sprintf("mask_beam_%d", id)
	)

	// Start building out the SVG
	var b strings.Builder

	if size > 0 {
		// Custom size
		_, _ = fmt.Fprintf(
			&b,
			`<svg viewBox="0 0 %d %d" fill="none" role="img"`+
				` xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`,
			beamSize, beamSize,
			size, size,
		)
	} else {
		_, _ = fmt.Fprintf(
			&b,
			`<svg viewBox="0 0 %d %d" fill="none" role="img"`+
				` xmlns="http://www.w3.org/2000/svg">`,
			beamSize, beamSize,
		)
	}

	// Mask group
	_, _ = fmt.Fprintf(
		&b,
		`<mask id="%s" maskUnits="userSpaceOnUse" x="0" y="0" width="%d" height="%d">`,
		maskID,
		beamSize, beamSize,
	)

	if square {
		_, _ = fmt.Fprintf(
			&b,
			`<rect width="%d" height="%d" fill="#FFFFFF"/>`,
			beamSize, beamSize,
		)
	} else {
		_, _ = fmt.Fprintf(
			&b,
			`<rect width="%d" height="%d" rx="%d" fill="#FFFFFF"/>`,
			beamSize, beamSize,
			beamSize*2,
		)
	}

	b.WriteString(`</mask>`)

	// Masked group
	_, _ = fmt.Fprintf(&b, `<g mask="url(#%s)">`, maskID)

	// Background
	_, _ = fmt.Fprintf(
		&b,
		`<rect width="%d" height="%d" fill="%s"/>`,
		beamSize, beamSize,
		p.colors.background,
	)

	// Wrapper square / circle
	wrapperRX := beamSize / 6
	if p.wrapper.circle {
		wrapperRX = beamSize
	}

	_, _ = fmt.Fprintf(
		&b,
		`<rect x="0" y="0" width="%d" height="%d"`+
			` transform="translate(%.2f %.2f) rotate(%d %d %d) scale(%.2f)"`+
			` fill="%s" rx="%d"/>`,
		beamSize, beamSize,
		p.wrapper.translateX, p.wrapper.translateY,
		p.wrapper.rotate, beamSize/2, beamSize/2,
		p.wrapper.scale,
		p.colors.wrapper,
		wrapperRX,
	)

	// Face group
	_, _ = fmt.Fprintf(
		&b,
		`<g transform="translate(%.2f %.2f) rotate(%d %d %d)">`,
		p.face.translateX, p.face.translateY,
		p.face.rotate, beamSize/2, beamSize/2,
	)

	// Mouth
	if p.face.mouthOpen {
		_, _ = fmt.Fprintf(
			&b,
			`<path d="M15 %dc2 1 4 1 6 0" stroke="%s" fill="none" stroke-linecap="round"/>`,
			19+p.face.mouthSpread, p.colors.face,
		)
	} else {
		_, _ = fmt.Fprintf(
			&b,
			`<path d="M13,%d a1,0.75 0 0,0 10,0" fill="%s"/>`,
			19+p.face.mouthSpread, p.colors.face,
		)
	}

	// Eyes
	_, _ = fmt.Fprintf(
		&b,
		`<rect x="%d" y="14" width="1.5" height="2" rx="1" stroke="none" fill="%s"/>`,
		14-p.face.eyeSpread, p.colors.face,
	)
	_, _ = fmt.Fprintf(
		&b,
		`<rect x="%d" y="14" width="1.5" height="2" rx="1" stroke="none" fill="%s"/>`,
		20+p.face.eyeSpread, p.colors.face,
	)

	// Group closures
	b.WriteString(`</g></g>`)

	// Final svg closure
	b.WriteString(`</svg>`)

	return b.String()
}
