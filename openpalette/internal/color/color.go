package color

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func HexToRGB(hex string) (int, int, int, error) {
	h, err := normalizeHex(hex)
	if err != nil {
		return 0, 0, 0, err
	}

	parsedHex, err := strconv.ParseUint(h, 16, 32)
	if err != nil {
		return 0, 0, 0, err
	}

	r := int((parsedHex >> 16) & 0xFF)
	g := int((parsedHex >> 8) & 0xFF)
	b := int(parsedHex & 0xFF)

	return r, g, b, nil
}

func RGBToHSL(r, g, b int) (float64, float64, float64) {
	rf := float64(r) / 255
	gf := float64(g) / 255
	bf := float64(b) / 255

	max := math.Max(rf, math.Max(gf, bf))
	min := math.Min(rf, math.Min(gf, bf))
	l := (max + min) / 2
	h, s := 0.0, 0.0

	if max != min {
		d := max - min
		if l > 0 && l < 1 {
			s = d / (1 - math.Abs(2*l-1))
		}
		switch max {
		case rf:
			h = (gf - bf) / d
			if gf < bf {
				h += 6
			}
		case gf:
			h = (bf-rf)/d + 2
		case bf:
			h = (rf-gf)/d + 4
		}
		h *= 60
	}

	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}
	s = clamp01(s)
	l = clamp01(l)

	return h, s, l
}

func HexToHSL(hex string) (float64, float64, float64, error) {
	r, g, b, err := HexToRGB(hex)
	if err != nil {
		return 0, 0, 0, err
	}
	h, s, l := RGBToHSL(r, g, b)

	return h, s, l, nil
}

func AdjustBrightness(hex string, dark bool) (string, error) {
	r, g, b, err := HexToRGB(hex)
	if err != nil {
		return "", err
	}

	l, a, bVal := rgbToLAB(r, g, b)
	lchL, c, h := labToLCH(l, a, bVal)

	// Apply brightness adjustments to match TypeScript behavior
	if dark {
		lchL *= 0.94
		c += 8
	} else {
		lchL *= 1.09
		// Don't add to chroma for light themes
	}

	// Add 2 to hue for both dark and light (this matches the TS code)
	h += 2
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}

	// Clamp L value
	if lchL > 100 {
		lchL = 100
	} else if lchL < 0 {
		lchL = 0
	}

	l, a, bVal = lchToLAB(lchL, c, h)
	rf, gf, bf := labToRGB(l, a, bVal)

	// Use rounding instead of truncation to match TypeScript
	return ToHex(int(math.Round(rf)), int(math.Round(gf)), int(math.Round(bf))), nil
}

func rgbToLAB(r, g, b int) (float64, float64, float64) {
	rf := gammaCorrectInverse(float64(r) / 255.0)
	gf := gammaCorrectInverse(float64(g) / 255.0)
	bf := gammaCorrectInverse(float64(b) / 255.0)

	// sRGB to XYZ conversion matrix (more precise values to match colorjs)
	x := rf*0.41239079926595934 + gf*0.357584339383878 + bf*0.1804807884018343
	y := rf*0.21263900587151027 + gf*0.715168678767756 + bf*0.07219231536073371
	z := rf*0.01933081871559182 + gf*0.11919477979462598 + bf*0.9505321522496607

	// D65 reference white normalization (more precise)
	x /= 0.9504559270516716
	y /= 1.0
	z /= 1.0890577507598784

	x = pivotXYZ(x)
	y = pivotXYZ(y)
	z = pivotXYZ(z)

	l := (116 * y) - 16
	a := 500 * (x - y)
	bVal := 200 * (y - z)

	return l, a, bVal
}

func labToLCH(l, a, b float64) (float64, float64, float64) {
	c := math.Sqrt(a*a + b*b)
	h := math.Atan2(b, a) * (180 / math.Pi)
	if h < 0 {
		h += 360
	}
	return l, c, h
}

func lchToLAB(l, c, h float64) (float64, float64, float64) {
	hRad := h * math.Pi / 180
	a := c * math.Cos(hRad)
	b := c * math.Sin(hRad)
	return l, a, b
}

func labToRGB(l, a, b float64) (float64, float64, float64) {
	y := (l + 16) / 116
	x := a/500 + y
	z := y - b/200

	x = inversePivotXYZ(x) * 0.9504559270516716
	y = inversePivotXYZ(y) * 1.0
	z = inversePivotXYZ(z) * 1.0890577507598784

	// XYZ to sRGB conversion matrix (more precise to match colorjs)
	rf := x*3.2409699419045226 + y*-1.537383177570094 + z*-0.4986107602930034
	gf := x*-0.9692436362808796 + y*1.8759675015077202 + z*0.04155505740717559
	bf := x*0.05563007969699366 + y*-0.20397695888897652 + z*1.0569715142428786

	rf = gammaCorrect(rf)
	gf = gammaCorrect(gf)
	bf = gammaCorrect(bf)

	return clamp255(rf * 255), clamp255(gf * 255), clamp255(bf * 255)
}

func normalizeHex(hex string) (string, error) {
	h := strings.TrimPrefix(strings.TrimSpace(hex), "#")
	if len(h) == 3 {
		h = string([]byte{h[0], h[0], h[1], h[1], h[2], h[2]})
	}
	if len(h) != 6 {
		return "", errors.New("hex must be 3 or 6 digits")
	}
	return h, nil
}

func ToHex(r, g, b int) string {
	// Clamp values to valid range
	r = int(clamp255(float64(r)))
	g = int(clamp255(float64(g)))
	b = int(clamp255(float64(b)))
	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}

func gammaCorrect(value float64) float64 {
	if value <= 0.0031308 {
		return 12.92 * value
	}
	return 1.055*math.Pow(value, 1/2.4) - 0.055
}

func gammaCorrectInverse(value float64) float64 {
	if value <= 0.04045 {
		return value / 12.92
	}
	return math.Pow((value+0.055)/1.055, 2.4)
}

func pivotXYZ(value float64) float64 {
	if value > 0.008856451679035631 {
		return math.Pow(value, 1.0/3.0)
	}
	return (value*903.2962962962961 + 16.0) / 116.0
}

func inversePivotXYZ(value float64) float64 {
	cube := math.Pow(value, 3)
	if cube > 0.008856451679035631 {
		return cube
	}
	return (value*116.0 - 16.0) / 903.2962962962961
}

func clamp255(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return v
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
