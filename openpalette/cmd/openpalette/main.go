package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/onedarktheme/onedark/openpalette/internal/color"
	"github.com/onedarktheme/onedark/openpalette/internal/data"
	"github.com/onedarktheme/onedark/openpalette/internal/model"
)

func marshalWithOrder(data map[string]interface{}) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("{\n")

	buf.WriteString(`  "version": "` + data["version"].(string) + `"`)

	for variantKey, variantData := range data {
		if variantKey == "version" {
			continue
		}

		variant := variantData.(model.PaletteVariant)
		buf.WriteString(",\n")
		buf.WriteString(`  "` + variantKey + `": {` + "\n")
		buf.WriteString(`    "name": "` + variant.Name + `",` + "\n")
		buf.WriteString(`    "emoji": "` + variant.Emoji + `",` + "\n")
		buf.WriteString(fmt.Sprintf(`    "order": %d,`, variant.Order) + "\n")
		buf.WriteString(fmt.Sprintf(`    "dark": %t,`, variant.Dark) + "\n")
		buf.WriteString(`    "colors": {` + "\n")

		type colorEntry struct {
			key   string
			color model.ColorDefinition
		}
		var colors []colorEntry
		for k, c := range variant.Colors {
			colors = append(colors, colorEntry{k, c})
		}
		sort.Slice(colors, func(i, j int) bool {
			return colors[i].color.Order < colors[j].color.Order
		})

		for i, entry := range colors {
			if i > 0 {
				buf.WriteString(",\n")
			}
			colorJSON, _ := json.MarshalIndent(entry.color, "      ", "  ")
			buf.WriteString(`      "` + entry.key + `": ` + string(colorJSON))
		}

		buf.WriteString("\n    },\n")
		buf.WriteString(`    "ansiColors": {` + "\n")

		type ansiEntry struct {
			key  string
			ansi model.ANSI
		}
		var ansiColors []ansiEntry
		for k, a := range variant.AnsiColors {
			ansiColors = append(ansiColors, ansiEntry{k, a})
		}
		sort.Slice(ansiColors, func(i, j int) bool {
			return ansiColors[i].ansi.Order < ansiColors[j].ansi.Order
		})

		for i, entry := range ansiColors {
			if i > 0 {
				buf.WriteString(",\n")
			}
			ansiJSON, _ := json.MarshalIndent(entry.ansi, "      ", "  ")
			buf.WriteString(`      "` + entry.key + `": ` + string(ansiJSON))
		}

		buf.WriteString("\n    }\n")
		buf.WriteString("  }")
	}

	buf.WriteString("\n}")
	return buf.Bytes(), nil
}

func main() {
	output := make(map[string]interface{})
	output["version"] = "1.7.1"

	for variantKey, variant := range data.Definitions {
		processedVariant := model.PaletteVariant{
			Name:       variant.Name,
			Emoji:      variant.Emoji,
			Order:      variant.Order,
			Dark:       variant.Dark,
			Colors:     make(map[string]model.ColorDefinition),
			AnsiColors: generateAnsiColors(variant),
		}

		for colorKey, colorDef := range variant.Colors {
			r, g, b, err := color.HexToRGB(colorDef.Hex)
			if err != nil {
				fmt.Printf("Error converting hex %s to RGB: %v\n", colorDef.Hex, err)
				continue
			}

			h, s, l := color.RGBToHSL(r, g, b)

			processedColorDef := model.ColorDefinition{
				Name:   colorDef.Name,
				Order:  colorDef.Order,
				Hex:    colorDef.Hex,
				RGB:    model.RGB{R: r, G: g, B: b},
				HSL:    model.HSL{H: h, S: s, L: l},
				Accent: colorDef.Accent,
			}

			processedVariant.Colors[colorKey] = processedColorDef
		}

		output[strings.ToLower(variantKey)] = processedVariant
	}

	dataJSON, err := marshalWithOrder(output)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	if err := os.WriteFile("palette.json", dataJSON, 0644); err != nil {
		fmt.Println("Error writing palette.json:", err)
		return
	}

	fmt.Println("palette.json successfully created!")
}

func generateAnsiColors(variant model.PaletteVariant) map[string]model.ANSI {
	ansiColors := make(map[string]model.ANSI)

	ansiMappings := []struct {
		name      string
		order     int
		normalKey string
		brightHex string
		code      int
	}{
		{"black", 0, "subtext1", "", 0},
		{"red", 1, "red", "", 1},
		{"green", 2, "green", "", 2},
		{"yellow", 3, "yellow", "", 3},
		{"blue", 4, "blue", "", 4},
		{"magenta", 5, "pink", "", 5},
		{"cyan", 6, "teal", "", 6},
		{"white", 7, "surface2", "", 7},
	}

	for _, mapping := range ansiMappings {
		normalColor, exists := variant.Colors[mapping.normalKey]
		if !exists {
			continue
		}

		r, g, b, err := color.HexToRGB(normalColor.Hex)
		if err != nil {
			continue
		}
		h, s, l := color.RGBToHSL(r, g, b)

		normalColorDef := model.AnsiColorDefinition{
			Name: mapping.name,
			Hex:  normalColor.Hex,
			RGB:  model.RGB{R: r, G: g, B: b},
			HSL:  model.HSL{H: h, S: s, L: l},
			Code: mapping.code,
		}

		brightHex := mapping.brightHex
		if brightHex == "" {
			adjustedHex, err := color.AdjustBrightness(normalColor.Hex, variant.Dark)
			if err == nil {
				brightHex = adjustedHex
			} else {
				brightHex = normalColor.Hex
			}
		}

		br, bg, bb, err := color.HexToRGB(brightHex)
		if err != nil {
			continue
		}
		bh, bs, bl := color.RGBToHSL(br, bg, bb)

		brightColorDef := model.AnsiColorDefinition{
			Name: "Bright " + strings.Title(mapping.name),
			Hex:  brightHex,
			RGB:  model.RGB{R: br, G: bg, B: bb},
			HSL:  model.HSL{H: bh, S: bs, L: bl},
			Code: mapping.code + 8,
		}

		if mapping.name == "black" {
			normalColorDef.Name = "Black"
			if subtext0, exists := variant.Colors["subtext0"]; exists {
				br, bg, bb, err := color.HexToRGB(subtext0.Hex)
				if err == nil {
					bh, bs, bl := color.RGBToHSL(br, bg, bb)
					brightColorDef = model.AnsiColorDefinition{
						Name: "Bright Black",
						Hex:  subtext0.Hex,
						RGB:  model.RGB{R: br, G: bg, B: bb},
						HSL:  model.HSL{H: bh, S: bs, L: bl},
						Code: 8,
					}
				}
			}
		} else if mapping.name == "white" {
			normalColorDef.Name = "White"
			if surface1, exists := variant.Colors["surface1"]; exists {
				br, bg, bb, err := color.HexToRGB(surface1.Hex)
				if err == nil {
					bh, bs, bl := color.RGBToHSL(br, bg, bb)
					brightColorDef = model.AnsiColorDefinition{
						Name: "Bright White",
						Hex:  surface1.Hex,
						RGB:  model.RGB{R: br, G: bg, B: bb},
						HSL:  model.HSL{H: bh, S: bs, L: bl},
						Code: 15,
					}
				}
			}
		} else {
			normalColorDef.Name = strings.Title(mapping.name)
		}

		ansiColors[mapping.name] = model.ANSI{
			Name:   strings.Title(mapping.name),
			Order:  mapping.order,
			Normal: normalColorDef,
			Bright: brightColorDef,
		}
	}

	return ansiColors
}
