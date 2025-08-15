package builder

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/onedarktheme/onedark/openpalette/internal/color"
	"github.com/onedarktheme/onedark/openpalette/internal/model"
)

type PaletteBuilder struct {
	palette model.Palette
}

func NewPaletteBuilder(version string) *PaletteBuilder {
	return &PaletteBuilder{
		palette: model.Palette{
			Version:  version,
			Variants: make(map[string]model.PaletteVariant),
		},
	}
}

func (pb *PaletteBuilder) AddVariant(name, emoji string, dark bool, order int) {
	pb.palette.Variants[name] = model.PaletteVariant{
		Name:       name,
		Emoji:      emoji,
		Order:      order,
		Dark:       dark,
		Colors:     make(map[string]model.ColorDefinition),
		AnsiColors: make(map[string]model.ANSI),
	}
}

func (pb *PaletteBuilder) AddColor(variantName, colorName, hex string, accent bool, order int) error {
	variant, exists := pb.palette.Variants[variantName]
	if !exists {
		return fmt.Errorf("variant %s not found", variantName)
	}

	r, g, b, err := color.HexToRGB(hex)
	if err != nil {
		return err
	}

	h, s, l := color.RGBToHSL(int(r), int(g), int(b))

	variant.Colors[colorName] = model.ColorDefinition{
		Name:   colorName,
		Hex:    hex,
		RGB:    model.RGB{R: r, G: g, B: b},
		HSL:    model.HSL{H: h, S: s, L: l},
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
