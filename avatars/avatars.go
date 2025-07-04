package avatars

import (
	"errors"
	"fmt"
)

var ErrUnknownStyle = errors.New("unknown style")

type (
	// Style represents a distinct avatar style
	Style   string
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
) (string, error) {
	switch style {
	case Beam:
		return GenerateBeam(name, palette, size, square), nil
	case Bauhaus:
		return GenerateBauhaus(name, palette, size, square), nil
	case Marble:
		return GenerateMarble(name, palette, size, square), nil
	case Pixel:
		return GeneratePixel(name, palette, size, square), nil
	case Ring:
		return GenerateRing(name, palette, size, square), nil
	case Sunset:
		return GenerateSunset(name, palette, size, square), nil
	default:
		return "", fmt.Errorf("%w: %q", ErrUnknownStyle, style)
	}
}
