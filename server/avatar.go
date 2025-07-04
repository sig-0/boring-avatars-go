package server

import (
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sig-0/boring-avatars-go/avatars"
)

const (
	defaultVariant = avatars.Marble
	defaultSize    = 80 // px

	nameParam    = "name"
	variantParam = "variant"
	sizeParam    = "size"
	squareParam  = "square"
	colorsParam  = "colors"
)

// avatarHandler serves
// GET /?name&variant&size&colors&square
func avatarHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	// Fetch the name
	name := q.Get(nameParam)
	if name == "" {
		// No name provided, generate a random avatar
		name = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	// Fetch the variant
	variant := avatars.Style(strings.ToLower(q.Get(variantParam)))
	if variant == "" {
		variant = defaultVariant
	}

	if !avatars.ValidStyle(variant) {
		http.Error(w, "invalid variant", http.StatusBadRequest)

		return
	}

	// Fetch the size
	size := defaultSize

	if sz := q.Get(sizeParam); sz != "" {
		n, err := strconv.Atoi(sz)
		if err != nil || n <= 0 || n > 512 { // TODO expose this in the config
			http.Error(w, "invalid size (1-512)", http.StatusBadRequest)

			return
		}

		size = n
	}

	// Fetch the square flag
	square := false
	if v := q.Get(squareParam); v == "true" {
		square = true
	}

	// Fetch the color palette
	var palette avatars.Palette

	if cs := q.Get(colorsParam); cs != "" {
		for _, c := range strings.Split(cs, ",") {
			c = strings.TrimSpace(c)

			if !strings.HasPrefix(c, "#") {
				c = "#" + c
			}

			_, err := hex.DecodeString(c[1:])
			if len(c) != 7 || err != nil {
				http.Error(w, "colors must be 6-digit hex, comma-separated", http.StatusBadRequest)

				return
			}

			palette = append(palette, c)
		}
	}

	// Generate the SVG
	svg := avatars.Generate(variant, name, palette, size, square)

	w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")

	_, _ = io.WriteString(w, svg)
}
