package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

type RGB struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
}

type HSL struct {
	H float64 `json:"h"`
	S float64 `json:"s"`
	L float64 `json:"l"`
}

type PaletteColor struct {
	Name   string `json:"name"`
	Order  int    `json:"order"`
	Hex    string `json:"hex"`
	RGB    RGB    `json:"rgb"`
	HSL    HSL    `json:"hsl"`
	Accent bool   `json:"accent"`
}

type ANSIVariant struct {
	Name string `json:"name"`
	Hex  string `json:"hex"`
	RGB  RGB    `json:"rgb"`
	HSL  HSL    `json:"hsl"`
	Code int    `json:"code"`
}

type ANSIColor struct {
	Name   string      `json:"name"`
	Order  int         `json:"order"`
	Normal ANSIVariant `json:"normal"`
	Bright ANSIVariant `json:"bright"`
}

type PaletteVariant struct {
	Name              string                  `json:"name"`
	Emoji             string                  `json:"emoji"`
	Order             int                     `json:"order"`
	Dark              bool                    `json:"dark"`
	PaletteColors     map[string]PaletteColor `json:"colors"`
	AnsiPaletteColors map[string]ANSIColor    `json:"ansiColors"`
}

type PaletteResult struct {
	Version  string                    `json:"version"`
	Variants map[string]PaletteVariant `json:"-"`
}

type RawPaletteColor struct {
	ID     string
	Name   string
	Hex    string
	Accent bool
}

type RawVariant struct {
	ID            string
	Name          string
	Emoji         string
	Dark          bool
	PaletteColors []RawPaletteColor
}

func getRawVariants() []RawVariant {
	return []RawVariant{
		{
			ID:    "latte",
			Name:  "Latte",
			Emoji: "🌻",
			Dark:  false,
			PaletteColors: []RawPaletteColor{
				{ID: "rosewater", Name: "Rosewater", Hex: "#dc8a78", Accent: true},
				{ID: "flamingo", Name: "Flamingo", Hex: "#dd7878", Accent: true},
				{ID: "pink", Name: "Pink", Hex: "#ea76cb", Accent: true},
				{ID: "mauve", Name: "Mauve", Hex: "#8839ef", Accent: true},
				{ID: "red", Name: "Red", Hex: "#d20f39", Accent: true},
				{ID: "maroon", Name: "Maroon", Hex: "#e64553", Accent: true},
				{ID: "peach", Name: "Peach", Hex: "#fe640b", Accent: true},
				{ID: "yellow", Name: "Yellow", Hex: "#df8e1d", Accent: true},
				{ID: "green", Name: "Green", Hex: "#40a02b", Accent: true},
				{ID: "teal", Name: "Teal", Hex: "#179299", Accent: true},
				{ID: "sky", Name: "Sky", Hex: "#04a5e5", Accent: true},
				{ID: "sapphire", Name: "Sapphire", Hex: "#209fb5", Accent: true},
				{ID: "blue", Name: "Blue", Hex: "#1e66f5", Accent: true},
				{ID: "lavender", Name: "Lavender", Hex: "#7287fd", Accent: true},
				{ID: "text", Name: "Text", Hex: "#4c4f69", Accent: false},
				{ID: "subtext1", Name: "Subtext 1", Hex: "#5c5f77", Accent: false},
				{ID: "subtext0", Name: "Subtext 0", Hex: "#6c6f85", Accent: false},
				{ID: "overlay2", Name: "Overlay 2", Hex: "#7c7f93", Accent: false},
				{ID: "overlay1", Name: "Overlay 1", Hex: "#8c8fa1", Accent: false},
				{ID: "overlay0", Name: "Overlay 0", Hex: "#9ca0b0", Accent: false},
				{ID: "surface2", Name: "Surface 2", Hex: "#acb0be", Accent: false},
				{ID: "surface1", Name: "Surface 1", Hex: "#bcc0cc", Accent: false},
				{ID: "surface0", Name: "Surface 0", Hex: "#ccd0da", Accent: false},
				{ID: "base", Name: "Base", Hex: "#eff1f5", Accent: false},
				{ID: "mantle", Name: "Mantle", Hex: "#e6e9ef", Accent: false},
				{ID: "crust", Name: "Crust", Hex: "#dce0e8", Accent: false},
			},
		},
	}
}

type ANSIMapping struct {
	NormalCode int
	BrightCode int
	Mapping    string
}

func getANSIMappings() map[string]ANSIMapping {
	return map[string]ANSIMapping{
		"black": {
			NormalCode: 0,
			BrightCode: 8,
			Mapping:    "",
		},
		"red": {
			NormalCode: 1,
			BrightCode: 9,
			Mapping:    "red",
		},
		"green": {
			NormalCode: 2,
			BrightCode: 10,
			Mapping:    "green",
		},
		"yellow": {
			NormalCode: 3,
			BrightCode: 11,
			Mapping:    "yellow",
		},
		"blue": {
			NormalCode: 4,
			BrightCode: 12,
			Mapping:    "blue",
		},
		"magenta": {
			NormalCode: 5,
			BrightCode: 13,
			Mapping:    "pink",
		},
		"cyan": {
			NormalCode: 6,
			BrightCode: 14,
			Mapping:    "teal",
		},
		"white": {
			NormalCode: 7,
			BrightCode: 15,
			Mapping:    "",
		},
	}
}

// ProcessPalette converts raw variants to final palette structure
func ProcessPalette() PaletteResult {
	rawVariants := getRawVariants()
	ansiMappings := getANSIMappings()

	result := PaletteResult{
		Version:  "1.7.1", // Match original version
		Variants: make(map[string]PaletteVariant),
	}

	// Process each variant (latte, frappe, macchiato, mocha)
	for variantIndex, rawVariant := range rawVariants {
		variant := PaletteVariant{
			Name:              rawVariant.Name,
			Emoji:             rawVariant.Emoji,
			Order:             variantIndex,
			Dark:              rawVariant.Dark,
			PaletteColors:     make(map[string]PaletteColor),
			AnsiPaletteColors: make(map[string]ANSIColor),
		}

		// Process regular colors
		for colorIndex, rawColor := range rawVariant.PaletteColors {
			variant.PaletteColors[rawColor.ID] = ProcessColor(rawColor, colorIndex)
		}

		// Process ANSI colors
		for ansiIndex, ansiName := range []string{"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white"} {
			ansiMapping := ansiMappings[ansiName]
			variant.AnsiPaletteColors[ansiName] = ProcessANSIColor(ansiName, ansiMapping, ansiIndex, variant, rawVariant.Dark)
		}

		result.Variants[rawVariant.ID] = variant
	}

	return result
}

// ProcessColor converts a raw color to final palette color
func ProcessColor(rawColor RawPaletteColor, order int) PaletteColor {
	color := NewColor(rawColor.Hex)

	// Get RGB coordinates (matches color.to("srgb").toGamut().coords.map(i => Math.round(i * 255)))
	coords := color.ToSRGBGamut()
	rgb := RGB{
		R: int(math.Round(coords[0] * 255)),
		G: int(math.Round(coords[1] * 255)),
		B: int(math.Round(coords[2] * 255)),
	}

	// Get HSL (matches tinycolor(hex).toHsl())
	hsl := TinycolorHSL(rawColor.Hex)

	return PaletteColor{
		Name:   rawColor.Name,
		Order:  order,
		Hex:    rawColor.Hex,
		RGB:    rgb,
		HSL:    hsl,
		Accent: rawColor.Accent,
	}
}

// ProcessANSIColor creates ANSI color variants (normal and bright)
func ProcessANSIColor(ansiName string, mapping ANSIMapping, order int, variant PaletteVariant, isDark bool) ANSIColor {
	var normalColor *Color
	var normalName string

	// Determine normal color based on ANSI mapping
	if ansiName == "black" {
		// Special handling for black
		if isDark {
			normalColor = NewColor(findColorHex(variant, "surface1"))
		} else {
			normalColor = NewColor(findColorHex(variant, "subtext1"))
		}
		normalName = "Black"
	} else if ansiName == "white" {
		// Special handling for white
		if isDark {
			normalColor = NewColor(findColorHex(variant, "subtext0"))
		} else {
			normalColor = NewColor(findColorHex(variant, "surface2"))
		}
		normalName = "White"
	} else {
		// Use mapped color
		normalColor = NewColor(findColorHex(variant, mapping.Mapping))
		normalName = strings.Title(ansiName) // Capitalize first letter
	}

	// Create bright color by cloning and modifying in LCH space
	brightColor := normalColor.Clone()

	// Only apply LCH modifications for non-black/white colors
	if ansiName != "black" && ansiName != "white" {
		lch := brightColor.GetLCH()

		// Apply the same transformations as gen_palette.ts
		if isDark {
			lch.SetL(lch.L() * 0.94) // Darken for dark themes
			lch.SetC(lch.C() + 8)    // Increase chroma
		} else {
			lch.SetL(lch.L() * 1.09) // Brighten for light themes
			lch.SetC(lch.C() + 0)    // No chroma change for light themes
		}
		lch.SetH(lch.H() + 2) // Always shift hue by 2 degrees
	} else {
		// For black/white, bright version uses different base colors
		if ansiName == "black" {
			if isDark {
				brightColor = NewColor(findColorHex(variant, "surface2"))
			} else {
				brightColor = NewColor(findColorHex(variant, "subtext0"))
			}
		} else { // white
			if isDark {
				brightColor = NewColor(findColorHex(variant, "subtext1"))
			} else {
				brightColor = NewColor(findColorHex(variant, "surface1"))
			}
		}
	}

	// Convert to final format
	normalVariant := colorToANSIVariant(normalColor, normalName, mapping.NormalCode)
	brightVariant := colorToANSIVariant(brightColor, "Bright "+normalName, mapping.BrightCode)

	return ANSIColor{
		Name:   normalName,
		Order:  order,
		Normal: normalVariant,
		Bright: brightVariant,
	}
}

// Helper function to find a color's hex value by ID
func findColorHex(variant PaletteVariant, colorID string) string {
	if color, exists := variant.PaletteColors[colorID]; exists {
		return color.Hex
	}
	// If not found in processed colors, search in raw colors (shouldn't happen in normal flow)
	return "#000000" // Fallback
}

// Helper function to convert Color to ANSIVariant
func colorToANSIVariant(color *Color, name string, code int) ANSIVariant {
	coords := color.ToSRGBGamut()
	rgb := RGB{
		R: int(math.Round(coords[0] * 255)),
		G: int(math.Round(coords[1] * 255)),
		B: int(math.Round(coords[2] * 255)),
	}

	hex := color.ToString()
	hsl := TinycolorHSL(hex)

	return ANSIVariant{
		Name: name,
		Hex:  hex,
		RGB:  rgb,
		HSL:  hsl,
		Code: code,
	}
}

// Custom JSON marshaling to maintain exact order like the original
func (pr PaletteResult) MarshalJSON() ([]byte, error) {
	var buf strings.Builder
	buf.WriteString("{")

	// Always start with version
	buf.WriteString(fmt.Sprintf(`"version":"%s"`, pr.Version))

	// Add variants in the correct order: latte, frappe, macchiato, mocha
	variantOrder := []string{"latte", "frappe", "macchiato", "mocha"}
	for _, variantName := range variantOrder {
		if variant, exists := pr.Variants[variantName]; exists {
			buf.WriteString(",")
			buf.WriteString(fmt.Sprintf(`"%s":`, variantName))

			// Marshal the variant with custom ordering
			variantJSON, err := variant.MarshalJSON()
			if err != nil {
				return nil, err
			}
			buf.Write(variantJSON)
		}
	}

	buf.WriteString("}")
	return []byte(buf.String()), nil
}

// Custom JSON marshaling for PaletteVariant to maintain color order
func (pv PaletteVariant) MarshalJSON() ([]byte, error) {
	var buf strings.Builder
	buf.WriteString("{")

	// Basic properties first
	buf.WriteString(fmt.Sprintf(`"name":"%s"`, pv.Name))
	buf.WriteString(fmt.Sprintf(`,"emoji":"%s"`, pv.Emoji))
	buf.WriteString(fmt.Sprintf(`,"order":%d`, pv.Order))
	buf.WriteString(fmt.Sprintf(`,"dark":%t`, pv.Dark))

	// Colors in order
	buf.WriteString(`,"colors":{`)
	colorOrder := []string{
		"rosewater", "flamingo", "pink", "mauve", "red", "maroon", "peach", "yellow",
		"green", "teal", "sky", "sapphire", "blue", "lavender", "text", "subtext1",
		"subtext0", "overlay2", "overlay1", "overlay0", "surface2", "surface1",
		"surface0", "base", "mantle", "crust",
	}

	first := true
	for _, colorName := range colorOrder {
		if color, exists := pv.PaletteColors[colorName]; exists {
			if !first {
				buf.WriteString(",")
			}
			first = false

			buf.WriteString(fmt.Sprintf(`"%s":`, colorName))
			colorJSON, err := json.Marshal(color)
			if err != nil {
				return nil, err
			}
			buf.Write(colorJSON)
		}
	}
	buf.WriteString("}")

	// ANSI colors in order
	buf.WriteString(`,"ansiColors":{`)
	ansiOrder := []string{"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white"}

	first = true
	for _, ansiName := range ansiOrder {
		if ansiColor, exists := pv.AnsiPaletteColors[ansiName]; exists {
			if !first {
				buf.WriteString(",")
			}
			first = false

			buf.WriteString(fmt.Sprintf(`"%s":`, ansiName))
			ansiJSON, err := json.Marshal(ansiColor)
			if err != nil {
				return nil, err
			}
			buf.Write(ansiJSON)
		}
	}
	buf.WriteString("}")

	buf.WriteString("}")
	return []byte(buf.String()), nil
}

// WriteJSONFile writes the palette to a JSON file
func WriteJSONFile(palette PaletteResult, filename string) error {
	// Create the directory if it doesn't exist
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal to JSON with indentation (matches the original formatting)
	jsonData, err := json.MarshalIndent(palette, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func main() {
	fmt.Println("Starting palette generation...")

	// Generate the complete palette
	palette := ProcessPalette()

	fmt.Printf("Generated palette version %s\n", palette.Version)
	fmt.Printf("Variants: %d\n", len(palette.Variants))

	// Test one variant to make sure it worked
	if latte, exists := palette.Variants["latte"]; exists {
		fmt.Printf("Latte variant: %s %s (dark: %t)\n", latte.Name, latte.Emoji, latte.Dark)
		fmt.Printf("  Colors: %d\n", len(latte.PaletteColors))
		fmt.Printf("  ANSI Colors: %d\n", len(latte.AnsiPaletteColors))

		// Show a sample color
		if red, exists := latte.PaletteColors["red"]; exists {
			fmt.Printf("  Red: %s -> RGB(%d,%d,%d) HSL(%.1f,%.3f,%.3f)\n",
				red.Hex, red.RGB.R, red.RGB.G, red.RGB.B, red.HSL.H, red.HSL.S, red.HSL.L)
		}

		// Show a sample ANSI color
		if ansiRed, exists := latte.AnsiPaletteColors["red"]; exists {
			fmt.Printf("  ANSI Red: normal=%s bright=%s\n",
				ansiRed.Normal.Hex, ansiRed.Bright.Hex)
		}
	}

	// Write the JSON file
	outputFile := "palette.json"
	fmt.Printf("\nWriting palette to %s...\n", outputFile)

	if err := WriteJSONFile(palette, outputFile); err != nil {
		fmt.Printf("Error writing JSON file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated %s!\n", outputFile)

	// Show file size
	if info, err := os.Stat(outputFile); err == nil {
		fmt.Printf("File size: %d bytes\n", info.Size())
	}
}
