package avatars

import (
	"strconv"
	"strings"
)

const (
	black = "#000000"
	white = "#FFFFFF"
)

// NameToID converts the given name
// to a unique and deterministic ID representation
func NameToID(name string) int {
	var id int32

	for _, r := range name {
		id = ((id << 5) - id) + r
	}

	if id < 0 {
		id = -id
	}

	return int(id)
}

// IDToDigit returns the digit at 10^place in id
func IDToDigit(id int, place int) int {
	for i := 0; i < place; i++ {
		id /= 10
	}

	return id % 10
}

// IDToBoolean returns whether that digit is even
func IDToBoolean(id int, place int) bool {
	return IDToDigit(id, place)%2 == 0
}

// IDToPoint returns id%mod, negated when that digit is even and place > 0
func IDToPoint(id int, mod, place int) int {
	v := id % mod
	if place > 0 && IDToDigit(id, place)%2 == 0 {
		return -v
	}

	return v
}

// Contrast returns a cheap YIQ-based contrast
func Contrast(hex string) string {
	// Trim leading #
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return black
	}

	// Parse RGB
	r, _ := strconv.ParseInt(hex[0:2], 16, 64)
	g, _ := strconv.ParseInt(hex[2:4], 16, 64)
	b, _ := strconv.ParseInt(hex[4:6], 16, 64)

	// YIQ formula
	yiq := (float64(r)*299 + float64(g)*587 + float64(b)*114) / 1000
	if yiq >= 128 {
		return black
	}

	return white
}
