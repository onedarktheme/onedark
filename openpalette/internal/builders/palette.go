package builders

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/onedarktheme/onedark/openpalette/internal/utils/colors"
)

func NewPaletteBuilder(version string) *PaletteBuilder {
	return &PaletteBuilder{
		palette: Palette{
			Version:  version,
			Variants: make(map[string]PaletteVariant),
		},
	}
}

func (pb *PaletteBuilder) AddVariant(name, emoji string, dark bool, order int) {
	pb.palette.Variants[name] = PaletteVariant{
		Name:       name,
		Emoji:      emoji,
		Order:      order,
		Dark:       dark,
		Colors:     make(map[string]ColorDefinition),
		AnsiColors: make(map[string]ANSI),
	}
}

func (pb *PaletteBuilder) AddColor(variantName, colorName, hex string, accent bool, order int) error {
	variant, exists := pb.palette.Variants[variantName]
	if !exists {
		return fmt.Errorf("variant %s not found", variantName)
	}

	r, g, b, err := colors.HexToRGB(hex)
	if err != nil {
		return err
	}

	h, s, l := colors.RGBToHSL(int(r), int(g), int(b))

	variant.Colors[colorName] = ColorDefinition{
		Name:   colorName,
		Hex:    hex,
		RGB:    RGB{R: r, G: g, B: b},
		HSL:    HSL{H: h, S: s, L: l},
		Accent: accent,
		Order:  order,
	}

	pb.palette.Variants[variantName] = variant
	return nil
}

func (pb *PaletteBuilder) ExportJSON(filename string) error {
	data, err := json.MarshalIndent(pb.palette, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
