package avatars

import (
	"fmt"
	"strings"
)

const BeamSize = 36

var DefaultPalette = Palette{
	"#FFB703", "#219EBC", "#8ECAE6", "#023047", "#FB8500",
}

// BuildBeamParams deterministically derives all beam avatar parameters from an ID
func BuildBeamParams(id int, palette Palette) Params {
	if len(palette) == 0 {
		palette = DefaultPalette
	}

	translatePoint := func(id int, place int) float64 {
		point := IDToPoint(id, 10, place)

		if point < 5 {
			point += BeamSize / 9
		}

		return float64(point)
	}

	var (
		wrapperTranslateX = translatePoint(id, 1)
		wrapperTranslateY = translatePoint(id, 2)

		wrapperColor = palette[id%len(palette)]
	)

	// Prepare the params
	params := Params{
		Colors: Colors{
			Wrapper:    wrapperColor,
			Face:       Contrast(wrapperColor),
			Background: palette[(id+13)%len(palette)],
		},
		Wrapper: Wrapper{
			TranslateX: wrapperTranslateX,
			TranslateY: wrapperTranslateY,
			Rotate:     IDToPoint(id, 360, 0),
			Scale:      1 + float64(IDToPoint(id, BeamSize/12, 0))/10,
			Circle:     IDToBoolean(id, 1),
		},
	}

	// Prepare the face translation.
	// It's tied to the wrapper offset
	faceTranslateX := wrapperTranslateX / 2
	if wrapperTranslateX <= BeamSize/6 {
		faceTranslateX = float64(IDToPoint(id, 8, 1))
	}

	faceTranslateY := wrapperTranslateY / 2
	if wrapperTranslateY <= BeamSize/6 {
		faceTranslateY = float64(IDToPoint(id, 7, 2))
	}

	params.Face = Face{
		TranslateX:  faceTranslateX,
		TranslateY:  faceTranslateY,
		Rotate:      IDToPoint(id, 10, 3),
		EyeSpread:   IDToPoint(id, 5, 0),
		MouthSpread: IDToPoint(id, 3, 0),
		MouthOpen:   IDToBoolean(id, 2),
	}

	return params
}

// GenerateBeam returns a beam-style avatar SVG
func GenerateBeam(name string, palette Palette, size int, square bool) string {
	var (
		id     = NameToID(name)
		params = BuildBeamParams(id, palette)
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
			BeamSize, BeamSize,
			size, size,
		)
	} else {
		_, _ = fmt.Fprintf(
			&b,
			`<svg viewBox="0 0 %d %d" fill="none" role="img"`+
				` xmlns="http://www.w3.org/2000/svg">`,
			BeamSize, BeamSize,
		)
	}

	// Mask group
	_, _ = fmt.Fprintf(
		&b,
		`<mask id="%s" maskUnits="userSpaceOnUse" x="0" y="0" width="%d" height="%d">`,
		maskID,
		BeamSize, BeamSize,
	)

	if square {
		_, _ = fmt.Fprintf(
			&b,
			`<rect width="%d" height="%d" fill="#FFFFFF"/>`,
			BeamSize, BeamSize,
		)
	} else {
		_, _ = fmt.Fprintf(
			&b,
			`<rect width="%d" height="%d" rx="%d" fill="#FFFFFF"/>`,
			BeamSize, BeamSize,
			BeamSize*2,
		)
	}

	b.WriteString(`</mask>`)

	// Masked group
	_, _ = fmt.Fprintf(&b, `<g mask="url(#%s)">`, maskID)

	// Background
	_, _ = fmt.Fprintf(
		&b,
		`<rect width="%d" height="%d" fill="%s"/>`,
		BeamSize, BeamSize,
		params.Colors.Background,
	)

	// Wrapper square / circle
	wrapperRX := BeamSize / 6
	if params.Wrapper.Circle {
		wrapperRX = BeamSize
	}

	_, _ = fmt.Fprintf(
		&b,
		`<rect x="0" y="0" width="%d" height="%d"`+
			` transform="translate(%.2f %.2f) rotate(%d %d %d) scale(%.2f)"`+
			` fill="%s" rx="%d"/>`,
		BeamSize, BeamSize,
		params.Wrapper.TranslateX, params.Wrapper.TranslateY,
		params.Wrapper.Rotate, BeamSize/2, BeamSize/2,
		params.Wrapper.Scale,
		params.Colors.Wrapper,
		wrapperRX,
	)

	// Face group
	_, _ = fmt.Fprintf(
		&b,
		`<g transform="translate(%.2f %.2f) rotate(%d %d %d)">`,
		params.Face.TranslateX, params.Face.TranslateY,
		params.Face.Rotate, BeamSize/2, BeamSize/2,
	)

	// Mouth
	if params.Face.MouthOpen {
		_, _ = fmt.Fprintf(
			&b,
			`<path d="M15 %dc2 1 4 1 6 0" stroke="%s" fill="none" stroke-linecap="round"/>`,
			19+params.Face.MouthSpread, params.Colors.Face,
		)
	} else {
		_, _ = fmt.Fprintf(
			&b,
			`<path d="M13,%d a1,0.75 0 0,0 10,0" fill="%s"/>`,
			19+params.Face.MouthSpread, params.Colors.Face,
		)
	}

	// Eyes
	_, _ = fmt.Fprintf(
		&b,
		`<rect x="%d" y="14" width="1.5" height="2" rx="1" stroke="none" fill="%s"/>`,
		14-params.Face.EyeSpread, params.Colors.Face,
	)
	_, _ = fmt.Fprintf(
		&b,
		`<rect x="%d" y="14" width="1.5" height="2" rx="1" stroke="none" fill="%s"/>`,
		20+params.Face.EyeSpread, params.Colors.Face,
	)

	// Group closures
	b.WriteString(`</g></g>`)

	// Final svg closure
	b.WriteString(`</svg>`)

	return b.String()
}
