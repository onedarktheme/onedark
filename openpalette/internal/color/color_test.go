package color

import (
	"math"
	"testing"
)

func TestHexToRGB(t *testing.T) {
	tests := []struct {
		hex      string
		expected [3]float64
		wantErr  bool
	}{
		{"#FFFFFF", [3]float64{255, 255, 255}, false},
		{"#000000", [3]float64{0, 0, 0}, false},
		{"#FF0000", [3]float64{255, 0, 0}, false},
		{"#0F0", [3]float64{0, 255, 0}, false},
		{"123456", [3]float64{18, 52, 86}, false},
		{"#XYZ", [3]float64{0, 0, 0}, true}, // invalid
	}

	for _, tt := range tests {
		r, g, b, err := HexToRGB(tt.hex)
		if (err != nil) != tt.wantErr {
			t.Errorf("HexToRGB(%s) error = %v, wantErr %v", tt.hex, err, tt.wantErr)
			continue
		}
		t.Logf("HexToRGB(%s) -> R: %.2f, G: %.2f, B: %.2f", tt.hex, r, g, b)
		if !tt.wantErr {
			if r != tt.expected[0] || g != tt.expected[1] || b != tt.expected[2] {
				t.Errorf("HexToRGB(%s) = %v,%v,%v, want %v", tt.hex, r, g, b, tt.expected)
			}
		}
	}
}

func TestRGBToHSL(t *testing.T) {
	tests := []struct {
		r, g, b  float64
		expected [3]float64
	}{
		{255, 0, 0, [3]float64{0, 1, 0.5}},   // red
		{0, 255, 0, [3]float64{120, 1, 0.5}}, // green
		{0, 0, 255, [3]float64{240, 1, 0.5}}, // blue
		{255, 255, 255, [3]float64{0, 0, 1}}, // white
		{0, 0, 0, [3]float64{0, 0, 0}},       // black
	}

	for _, tt := range tests {
		h, s, l := RGBToHSL(tt.r, tt.g, tt.b)
		t.Logf("RGBToHSL(%.0f, %.0f, %.0f) -> H: %.2f, S: %.2f, L: %.2f", tt.r, tt.g, tt.b, h, s, l)
		if math.Abs(h-tt.expected[0]) > 0.01 || math.Abs(s-tt.expected[1]) > 0.01 || math.Abs(l-tt.expected[2]) > 0.01 {
			t.Errorf("RGBToHSL(%v,%v,%v) = %v,%v,%v, want %v", tt.r, tt.g, tt.b, h, s, l, tt.expected)
		}
	}
}

func TestAdjustBrightness(t *testing.T) {
	tests := []struct {
		hex     string
		factor  float64
		wantErr bool
	}{
		{"#000000", 2, false},   // should remain black (0)
		{"#FFFFFF", 0.5, false}, // white reduces
		{"#FF0000", 1.5, false}, // red increases brightness
		{"#00FF00", 0, false},   // green becomes black
		{"#XYZ", 1, true},       // invalid hex
	}

	for _, tt := range tests {
		newHex, err := AdjustBrightness(tt.hex, tt.factor)
		if (err != nil) != tt.wantErr {
			t.Errorf("AdjustBrightness(%s, %v) error = %v, wantErr %v", tt.hex, tt.factor, err, tt.wantErr)
			continue
		}
		t.Logf("AdjustBrightness(%s, %.2f) -> %s", tt.hex, tt.factor, newHex)
		if !tt.wantErr && len(newHex) != 7 {
			t.Errorf("AdjustBrightness(%s, %v) returned invalid hex %s", tt.hex, tt.factor, newHex)
		}
	}
}
