package data

import (
	"github.com/onedarktheme/onedark/openpalette/internal/model"
)

var Definitions = map[string]model.PaletteVariant{
	"latte": {
		Name:  "Latte",
		Emoji: "🌻",
		Order: 0,
		Dark:  false,
		Colors: map[string]model.ColorDefinition{
			"rosewater": {Name: "Rosewater", Order: 0, Hex: "#dc8a78", Accent: true},
			"flamingo":  {Name: "Flamingo", Order: 1, Hex: "#dd7878", Accent: true},
			"pink":      {Name: "Pink", Order: 2, Hex: "#ea76cb", Accent: true},
			"mauve":     {Name: "Mauve", Order: 3, Hex: "#8839ef", Accent: true},
			"red":       {Name: "Red", Order: 4, Hex: "#d20f39", Accent: true},
			"maroon":    {Name: "Maroon", Order: 5, Hex: "#e64553", Accent: true},
			"peach":     {Name: "Peach", Order: 6, Hex: "#fe640b", Accent: true},
			"yellow":    {Name: "Yellow", Order: 7, Hex: "#df8e1d", Accent: true},
			"green":     {Name: "Green", Order: 8, Hex: "#40a02b", Accent: true},
			"teal":      {Name: "Teal", Order: 9, Hex: "#179299", Accent: true},
			"sky":       {Name: "Sky", Order: 10, Hex: "#04a5e5", Accent: true},
			"sapphire":  {Name: "Sapphire", Order: 11, Hex: "#209fb5", Accent: true},
			"blue":      {Name: "Blue", Order: 12, Hex: "#1e66f5", Accent: true},
			"lavender":  {Name: "Lavender", Order: 13, Hex: "#7287fd", Accent: true},
			"text":      {Name: "Text", Order: 14, Hex: "#4c4f69", Accent: false},
			"subtext1":  {Name: "Subtext 1", Order: 15, Hex: "#5c5f77", Accent: false},
			"subtext0":  {Name: "Subtext 0", Order: 16, Hex: "#6c6f85", Accent: false},
			"overlay2":  {Name: "Overlay 2", Order: 17, Hex: "#7c7f93", Accent: false},
			"overlay1":  {Name: "Overlay 1", Order: 18, Hex: "#8c8fa1", Accent: false},
			"overlay0":  {Name: "Overlay 0", Order: 19, Hex: "#9ca0b0", Accent: false},
			"surface2":  {Name: "Surface 2", Order: 20, Hex: "#acb0be", Accent: false},
			"surface1":  {Name: "Surface 1", Order: 21, Hex: "#bcc0cc", Accent: false},
			"surface0":  {Name: "Surface 0", Order: 22, Hex: "#ccd0da", Accent: false},
			"base":      {Name: "Base", Order: 23, Hex: "#eff1f5", Accent: false},
			"mantle":    {Name: "Mantle", Order: 24, Hex: "#e6e9ef", Accent: false},
			"crust":     {Name: "Crust", Order: 25, Hex: "#dce0e8", Accent: false},
		},
	},
}
