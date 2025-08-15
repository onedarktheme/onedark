package color

import (
	"math"
	"strconv"
	"strings"
)

func HexToRGB(hex string) (uint8, uint8, uint8, error) {
	hex = strings.TrimPrefix(hex, "#")
	parsedHex, err := strconv.ParseUint(string(hex), 16, 32)
	if err != nil {
		return 0, 0, 0, err
	}

	var r, g, b uint8

	r = uint8(parsedHex >> 16)
	g = uint8((parsedHex >> 8) & 0xFF)
	b = uint8(parsedHex & 0xFF)

	return r, g, b, nil
}

func RGBToHSL(r, g, b int) (float64, float64, float64) {
	rf := float64(r) / 255
	gf := float64(g) / 255
	bf := float64(b) / 255

	max := math.Max(rf, math.Max(gf, bf))
	min := math.Min(rf, math.Min(gf, bf))
	h, s, l := 0.0, 0.0, (max+min)/2

	if max != min {
		d := max - min
		s = d / (1 - math.Abs(2*l-1))
		switch max {
		case rf:
			h = math.Mod(((gf - bf) / d), 6)
		case gf:
			h = ((bf - rf) / d) + 2
		case bf:
			h = ((rf - gf) / d) + 4
		}
		h *= 60
		if h < 0 {
			h += 360
		}
	}
	return h, s, l
}
