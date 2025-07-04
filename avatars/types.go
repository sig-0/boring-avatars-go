package avatars

type Palette []string

type Colors struct {
	Wrapper    string // color of the rotated square / circle
	Face       string // black or white for eyes and mouth (contrast to Wrapper)
	Background string // full-canvas background
}

type Wrapper struct {
	TranslateX, TranslateY float64 // translate
	Rotate                 int     // degrees
	Scale                  float64 // 1.0 - ~1.3
	Circle                 bool    // true -> circle; false -> square
}

type Face struct {
	TranslateX, TranslateY float64 // translate
	Rotate                 int     // degrees
	EyeSpread              int     // -5..+5 px
	MouthSpread            int     // -3..+3 px
	MouthOpen              bool    // open (curved line) or closed (arc)
}

type Params struct {
	Colors  Colors
	Face    Face
	Wrapper Wrapper
}
