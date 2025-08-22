package main

import (
	"fmt"
	"math"
	"testing"
)

// Test cases from the original palette.json
type ColorTestCase struct {
	name        string
	hex         string
	expectedRGB RGB
	expectedHSL HSL
}

func getColorTestCases() []ColorTestCase {
	return []ColorTestCase{
		{
			name:        "Latte Rosewater",
			hex:         "#dc8a78",
			expectedRGB: RGB{R: 220, G: 138, B: 120},
			expectedHSL: HSL{H: 10.799999999999995, S: 0.5882352941176472, L: 0.6666666666666667},
		},
		{
			name:        "Latte Red",
			hex:         "#d20f39",
			expectedRGB: RGB{R: 210, G: 15, B: 57},
			expectedHSL: HSL{H: 347.0769230769231, S: 0.8666666666666666, L: 0.4411764705882353},
		},
		{
			name:        "Latte Blue",
			hex:         "#1e66f5",
			expectedRGB: RGB{R: 30, G: 102, B: 245},
			expectedHSL: HSL{H: 219.90697674418607, S: 0.9148936170212768, L: 0.5392156862745098},
		},
		{
			name:        "Mocha Text",
			hex:         "#cdd6f4",
			expectedRGB: RGB{R: 205, G: 214, B: 244},
			expectedHSL: HSL{H: 226.15384615384616, S: 0.6393442622950825, L: 0.8803921568627451},
		},
	}
}

// ANSI bright color test cases from palette.json
type ANSITestCase struct {
	name      string
	normalHex string
	brightHex string
	isDark    bool
}

func getANSITestCases() []ANSITestCase {
	return []ANSITestCase{
		{
			name:      "Latte Red (light theme)",
			normalHex: "#d20f39",
			brightHex: "#de293e", // From palette.json latte.ansiColors.red.bright.hex
			isDark:    false,
		},
		{
			name:      "Mocha Red (dark theme)",
			normalHex: "#f38ba8",
			brightHex: "#f37799", // From palette.json mocha.ansiColors.red.bright.hex
			isDark:    true,
		},
		{
			name:      "Latte Blue (light theme)",
			normalHex: "#1e66f5",
			brightHex: "#456eff", // From palette.json latte.ansiColors.blue.bright.hex
			isDark:    false,
		},
		{
			name:      "Mocha Blue (dark theme)",
			normalHex: "#89b4fa",
			brightHex: "#74a8fc", // From palette.json mocha.ansiColors.blue.bright.hex
			isDark:    true,
		},
	}
}

func TestColorConversions(t *testing.T) {
	testCases := getColorTestCases()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			color := NewColor(tc.hex)

			// Test RGB conversion
			coords := color.ToSRGBGamut()
			actualRGB := RGB{
				R: int(math.Round(coords[0] * 255)),
				G: int(math.Round(coords[1] * 255)),
				B: int(math.Round(coords[2] * 255)),
			}

			if actualRGB != tc.expectedRGB {
				t.Errorf("RGB mismatch for %s:\nExpected: %+v\nActual: %+v",
					tc.name, tc.expectedRGB, actualRGB)
			}

			// Test HSL conversion
			actualHSL := TinycolorHSL(tc.hex)

			// Allow small floating point differences
			if !floatEqual(actualHSL.H, tc.expectedHSL.H, 0.0001) ||
				!floatEqual(actualHSL.S, tc.expectedHSL.S, 0.0001) ||
				!floatEqual(actualHSL.L, tc.expectedHSL.L, 0.0001) {
				t.Errorf("HSL mismatch for %s:\nExpected: H=%.6f S=%.6f L=%.6f\nActual: H=%.6f S=%.6f L=%.6f",
					tc.name, tc.expectedHSL.H, tc.expectedHSL.S, tc.expectedHSL.L,
					actualHSL.H, actualHSL.S, actualHSL.L)
			}
		})
	}
}

func TestANSIBrightColors(t *testing.T) {
	testCases := getANSITestCases()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create the normal color
			normalColor := NewColor(tc.normalHex)

			// Apply the bright color transformation
			brightColor := normalColor.Clone()
			lch := brightColor.GetLCH()

			// Apply the same transformation as gen_palette.ts
			if tc.isDark {
				lch.SetL(lch.L() * 0.94) // Darken for dark themes
				lch.SetC(lch.C() + 8)    // Increase chroma
			} else {
				lch.SetL(lch.L() * 1.09) // Brighten for light themes
				lch.SetC(lch.C() + 0)    // No chroma change for light themes
			}
			lch.SetH(lch.H() + 2) // Always shift hue by 2 degrees

			actualBrightHex := brightColor.ToString()

			if actualBrightHex != tc.brightHex {
				t.Errorf("ANSI bright color mismatch for %s:\nExpected: %s\nActual: %s",
					tc.name, tc.brightHex, actualBrightHex)
			}
		})
	}
}

// Helper function to compare floats with tolerance
func floatEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

// You can also run individual tests manually
func TestSingleColor(t *testing.T) {
	hex := "#dc8a78"
	color := NewColor(hex)

	fmt.Printf("Testing %s:\n", hex)
	fmt.Printf("ToString: %s\n", color.ToString())

	coords := color.ToSRGBGamut()
	actualRGB := RGB{
		R: int(math.Round(coords[0] * 255)),
		G: int(math.Round(coords[1] * 255)),
		B: int(math.Round(coords[2] * 255)),
	}
	fmt.Printf("RGB: %+v\n", actualRGB)

	actualHSL := TinycolorHSL(hex)
	fmt.Printf("HSL: H=%.6f S=%.6f L=%.6f\n", actualHSL.H, actualHSL.S, actualHSL.L)
}
