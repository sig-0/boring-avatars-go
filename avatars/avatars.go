package avatars

import (
	"errors"
)

var ErrUnknownStyle = errors.New("unknown style")

type (
	// Style represents a distinct avatar style
	Style string

	// Palette is the hex color palette
	Palette []string
)

var DefaultPalette = Palette{
	"#FFB703", "#219EBC", "#8ECAE6", "#023047", "#FB8500",
}

const (
	Beam    Style = "beam"
	Bauhaus Style = "bauhaus"
	Marble  Style = "marble"
	Pixel   Style = "pixel"
	Ring    Style = "ring"
	Sunset  Style = "sunset"
)

// Generate generates an avatar based on the requested style and params
func Generate(
	style Style,
	name string,
	palette Palette,
	size int,
	square bool,
) string {
	switch style {
	case Beam:
		return GenerateBeam(name, palette, size, square)
	case Bauhaus:
		return GenerateBauhaus(name, palette, size, square)
	case Pixel:
		return GeneratePixel(name, palette, size, square)
	case Ring:
		return GenerateRing(name, palette, size, square)
	case Sunset:
		return GenerateSunset(name, palette, size, square)
	default:
		return GenerateMarble(name, palette, size, square)
	}
}

// ValidStyle checks if the style is a valid boring avatar variant
func ValidStyle(style Style) bool {
	switch style {
	case Beam, Bauhaus, Marble, Pixel, Ring, Sunset:
		return true
	default:
		return false
	}
}
