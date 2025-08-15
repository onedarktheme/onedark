package colors_test

import (
	"math"
	"testing"

	"github.com/onedarktheme/onedark/palettegen/internal/utils/colors"
)

func TestHexToRGB(t *testing.T) {
	tests := []struct {
		hex        string
		wantR      uint8
		wantG      uint8
		wantB      uint8
		shouldFail bool
	}{
		{"#ff0000", 255, 0, 0, false},
		{"00ff00", 0, 255, 0, false},
		{"#0000ff", 0, 0, 255, false},
		{"#zzz", 0, 0, 0, true}, // invalid hex
	}

	for _, tt := range tests {
		r, g, b, err := colors.HexToRGB(tt.hex)
		if tt.shouldFail {
			if err == nil {
				t.Errorf("HexToRGB(%q) expected error but got none", tt.hex)
			}
			continue
		}
		if err != nil {
			t.Errorf("HexToRGB(%q) unexpected error: %v", tt.hex, err)
			continue
		}
		if r != tt.wantR || g != tt.wantG || b != tt.wantB {
			t.Errorf("HexToRGB(%q) = (%d, %d, %d), want (%d, %d, %d)",
				tt.hex, r, g, b, tt.wantR, tt.wantG, tt.wantB)
		}
	}
}

func TestRGBToHSL(t *testing.T) {
	tests := []struct {
		r, g, b int
		wantH   float64
		wantS   float64
		wantL   float64
	}{
		{255, 0, 0, 0, 1, 0.5},        // pure red
		{0, 255, 0, 120, 1, 0.5},      // pure green
		{0, 0, 255, 240, 1, 0.5},      // pure blue
		{255, 255, 255, 0, 0, 1},      // white
		{0, 0, 0, 0, 0, 0},            // black
		{128, 128, 128, 0, 0, 0.5019}, // gray
	}

	for _, tt := range tests {
		h, s, l := colors.RGBToHSL(tt.r, tt.g, tt.b)

		// Allow tiny floating-point differences
		if math.Abs(h-tt.wantH) > 0.5 {
			t.Errorf("RGBToHSL(%d, %d, %d) hue got %.2f, want %.2f", tt.r, tt.g, tt.b, h, tt.wantH)
		}
		if math.Abs(s-tt.wantS) > 0.01 {
			t.Errorf("RGBToHSL(%d, %d, %d) saturation got %.2f, want %.2f", tt.r, tt.g, tt.b, s, tt.wantS)
		}
		if math.Abs(l-tt.wantL) > 0.01 {
			t.Errorf("RGBToHSL(%d, %d, %d) lightness got %.4f, want %.4f", tt.r, tt.g, tt.b, l, tt.wantL)
		}
	}
}
