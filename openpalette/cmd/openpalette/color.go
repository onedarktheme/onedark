package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Color represents a color that can work in different color spaces
type Color struct {
	hex   string     // Original hex value
	lch   [3]float64 // L, C, H values when in LCH mode
	isLCH bool       // Track if we're working in LCH space
}

// NewColor creates a color from hex string (matches new Color("#hex"))
func NewColor(hex string) *Color {
	return &Color{
		hex:   strings.TrimPrefix(hex, "#"),
		isLCH: false,
	}
}

// Clone creates a copy (matches new Color(normalColor))
func (c *Color) Clone() *Color {
	clone := &Color{
		hex:   c.hex,
		isLCH: c.isLCH,
	}
	if c.isLCH {
		clone.lch = c.lch
	}
	return clone
}

// ToString returns hex string (matches color.toString({ format: "hex" }))
func (c *Color) ToString() string {
	if c.isLCH {
		// Convert LCH back to hex
		return c.lchToHex()
	}
	return "#" + c.hex
}

// ToSRGBGamut returns RGB coords like color.to("srgb").toGamut().coords
func (c *Color) ToSRGBGamut() [3]float64 {
	var r, g, b float64

	if c.isLCH {
		// Convert LCH to sRGB
		r, g, b = c.lchToSRGB()
	} else {
		// Parse hex to RGB (0-1 range)
		r, g, b = c.hexToSRGB()
	}

	// Gamut clamp (equivalent to .toGamut())
	r = clampFloat(r, 0, 1)
	g = clampFloat(g, 0, 1)
	b = clampFloat(b, 0, 1)

	return [3]float64{r, g, b}
}

// GetLCH returns LCH representation for manipulation
func (c *Color) GetLCH() *LCHColor {
	if !c.isLCH {
		// Convert hex to LCH
		c.lch = c.hexToLCH()
		c.isLCH = true
	}

	return &LCHColor{color: c}
}

// LCHColor provides access to L, C, H components
type LCHColor struct {
	color *Color
}

// L gets/sets lightness
func (lch *LCHColor) L() float64 {
	return lch.color.lch[0]
}

func (lch *LCHColor) SetL(value float64) {
	lch.color.lch[0] = value
}

// C gets/sets chroma
func (lch *LCHColor) C() float64 {
	return lch.color.lch[1]
}

func (lch *LCHColor) SetC(value float64) {
	lch.color.lch[1] = value
}

// H gets/sets hue
func (lch *LCHColor) H() float64 {
	return lch.color.lch[2]
}

func (lch *LCHColor) SetH(value float64) {
	lch.color.lch[2] = value
}

// hexToSRGB converts hex to sRGB (0-1 range)
func (c *Color) hexToSRGB() (float64, float64, float64) {
	r, _ := strconv.ParseInt(c.hex[0:2], 16, 0)
	g, _ := strconv.ParseInt(c.hex[2:4], 16, 0)
	b, _ := strconv.ParseInt(c.hex[4:6], 16, 0)

	return float64(r) / 255.0, float64(g) / 255.0, float64(b) / 255.0
}

// TinycolorHSL implements tinycolor2's exact toHsl() function
func TinycolorHSL(hex string) HSL {
	// Parse hex to RGB (0-255)
	hex = strings.TrimPrefix(hex, "#")
	r, _ := strconv.ParseInt(hex[0:2], 16, 0)
	g, _ := strconv.ParseInt(hex[2:4], 16, 0)
	b, _ := strconv.ParseInt(hex[4:6], 16, 0)

	// Convert to 0-1 range using tinycolor2's bound01 logic
	rNorm := bound01(float64(r), 255)
	gNorm := bound01(float64(g), 255)
	bNorm := bound01(float64(b), 255)

	max := math.Max(math.Max(rNorm, gNorm), bNorm)
	min := math.Min(math.Min(rNorm, gNorm), bNorm)

	var h, s, l float64
	l = (max + min) / 2

	if max == min {
		h = 0 // achromatic
		s = 0
	} else {
		d := max - min
		if l > 0.5 {
			s = d / (2 - max - min)
		} else {
			s = d / (max + min)
		}

		switch max {
		case rNorm:
			h = (gNorm - bNorm) / d
			if gNorm < bNorm {
				h += 6
			}
		case gNorm:
			h = (bNorm-rNorm)/d + 2
		case bNorm:
			h = (rNorm-gNorm)/d + 4
		}
		h /= 6
	}

	return HSL{
		H: h * 360, // Convert to degrees
		S: s,       // Keep as ratio
		L: l,       // Keep as ratio
	}
}

// bound01 implements tinycolor2's bound01 function exactly
func bound01(value, max float64) float64 {
	if max == 360 {
		// Special handling for hue
		n := value
		if math.Abs(n-max) < 0.000001 {
			return 1.0
		}
		if n < 0 {
			return (math.Mod(n, max) + max) / max
		}
		return math.Mod(n, max) / max
	}

	n := math.Min(max, math.Max(0, value))
	if math.Abs(n-max) < 0.000001 {
		return 1.0
	}
	return math.Mod(n, max) / max
}

func clampFloat(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// hexToLCH converts hex to LCH coordinates
func (c *Color) hexToLCH() [3]float64 {
	// Hex → sRGB → Linear RGB → XYZ → Lab → LCH
	r, g, b := c.hexToSRGB()

	// sRGB to Linear RGB (gamma correction removal)
	r = srgbToLinear(r)
	g = srgbToLinear(g)
	b = srgbToLinear(b)

	// Linear RGB to XYZ (using sRGB matrix)
	x, y, z := linearRGBToXYZ(r, g, b)

	// XYZ to Lab
	l, a, labB := xyzToLab(x, y, z)

	// Lab to LCH
	lch_l, lch_c, lch_h := labToLCH(l, a, labB)

	return [3]float64{lch_l, lch_c, lch_h}
}

// lchToSRGB converts LCH coordinates to sRGB
func (c *Color) lchToSRGB() (float64, float64, float64) {
	// LCH → Lab → XYZ → Linear RGB → sRGB
	l := c.lch[0]
	ch := c.lch[1]
	h := c.lch[2]

	// LCH to Lab
	lab_l, lab_a, lab_b := lchToLab(l, ch, h)

	// Lab to XYZ
	x, y, z := labToXYZ(lab_l, lab_a, lab_b)

	// XYZ to Linear RGB
	r, g, b := xyzToLinearRGB(x, y, z)

	// Linear RGB to sRGB (gamma correction)
	r = linearToSRGB(r)
	g = linearToSRGB(g)
	b = linearToSRGB(b)

	return r, g, b
}

// lchToHex converts LCH back to hex string
func (c *Color) lchToHex() string {
	r, g, b := c.lchToSRGB()

	// Convert to 0-255 range and clamp
	rInt := int(math.Round(clampFloat(r, 0, 1) * 255))
	gInt := int(math.Round(clampFloat(g, 0, 1) * 255))
	bInt := int(math.Round(clampFloat(b, 0, 1) * 255))

	return fmt.Sprintf("#%02x%02x%02x", rInt, gInt, bInt)
}

// sRGB gamma correction functions
func srgbToLinear(val float64) float64 {
	if val <= 0.04045 {
		return val / 12.92
	}
	return math.Pow((val+0.055)/1.055, 2.4)
}

func linearToSRGB(val float64) float64 {
	if val <= 0.0031308 {
		return 12.92 * val
	}
	return 1.055*math.Pow(val, 1.0/2.4) - 0.055
}

// Linear RGB to XYZ conversion (D50 white point matrix from Color.js)
func linearRGBToXYZ(r, g, b float64) (float64, float64, float64) {
	// Matrix for sRGB to XYZ-D50 (from Color.js source)
	x := 0.4360747*r + 0.3850649*g + 0.1430804*b
	y := 0.2225045*r + 0.7168786*g + 0.0606169*b
	z := 0.0139322*r + 0.0971045*g + 0.7141733*b

	return x, y, z
}

// XYZ to Linear RGB conversion
func xyzToLinearRGB(x, y, z float64) (float64, float64, float64) {
	// Inverse matrix for XYZ-D50 to sRGB (from Color.js source)
	r := 3.1338561*x - 1.6168667*y - 0.4906146*z
	g := -0.9787684*x + 1.9161415*y + 0.0334540*z
	b := 0.0719453*x - 0.2289914*y + 1.4052427*z

	return r, g, b
}

// XYZ to Lab conversion (D50 white point)
func xyzToLab(x, y, z float64) (float64, float64, float64) {
	// D50 white point values
	const xn = 0.9642956 // D50 white point X
	const yn = 1.0       // D50 white point Y
	const zn = 0.8251046 // D50 white point Z

	// Constants from CIE Lab definition
	const epsilon = 216.0 / 24389.0 // 6^3/29^3
	const kappa = 24389.0 / 27.0    // 29^3/3^3

	// Normalize by white point
	fx := x / xn
	fy := y / yn
	fz := z / zn

	// Apply Lab transform
	if fx > epsilon {
		fx = math.Cbrt(fx)
	} else {
		fx = (kappa*fx + 16) / 116
	}

	if fy > epsilon {
		fy = math.Cbrt(fy)
	} else {
		fy = (kappa*fy + 16) / 116
	}

	if fz > epsilon {
		fz = math.Cbrt(fz)
	} else {
		fz = (kappa*fz + 16) / 116
	}

	l := 116*fy - 16
	a := 500 * (fx - fy)
	b := 200 * (fy - fz)

	return l, a, b
}

// Lab to XYZ conversion
func labToXYZ(l, a, b float64) (float64, float64, float64) {
	// D50 white point values
	const xn = 0.9642956
	const yn = 1.0
	const zn = 0.8251046

	// Constants
	const epsilon = 216.0 / 24389.0
	const kappa = 24389.0 / 27.0
	const epsilon3 = 24.0 / 116.0 // epsilon^(1/3)

	// Compute f values
	fy := (l + 16) / 116
	fx := a/500 + fy
	fz := fy - b/200

	// Convert f values back to XYZ
	var x, y, z float64

	if fx*fx*fx > epsilon {
		x = fx * fx * fx
	} else {
		x = (116*fx - 16) / kappa
	}

	if l > 8 {
		y = math.Pow((l+16)/116, 3)
	} else {
		y = l / kappa
	}

	if fz*fz*fz > epsilon {
		z = fz * fz * fz
	} else {
		z = (116*fz - 16) / kappa
	}

	// Scale by white point
	return x * xn, y * yn, z * zn
}

// Lab to LCH conversion
func labToLCH(l, a, b float64) (float64, float64, float64) {
	c := math.Sqrt(a*a + b*b)
	h := math.Atan2(b, a) * 180 / math.Pi

	// Normalize hue to 0-360
	if h < 0 {
		h += 360
	}

	return l, c, h
}

// LCH to Lab conversion
func lchToLab(l, c, h float64) (float64, float64, float64) {
	hRad := h * math.Pi / 180
	a := c * math.Cos(hRad)
	b := c * math.Sin(hRad)

	return l, a, b
}
