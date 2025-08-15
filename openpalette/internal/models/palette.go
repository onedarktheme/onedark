package models

type RGB struct {
	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
}

type HSL struct {
	H float64 `json:"h"`
	S float64 `json:"s"`
	L float64 `json:"l"`
}

type ColorDefinition struct {
	Name   string `json:"name"`
	Hex    string `json:"hex"`
	RGB    RGB    `json:"rgb"`
	HSL    HSL    `json:"hsl"`
	Accent bool   `json:"accent"`
	Order  int    `json:"order"`
}

type ColorDefinitionWithCode struct {
	ColorDefinition
	Code int `json:"code"`
}

type ANSI struct {
	Name   string                  `json:"name"`
	Order  int                     `json:"order"`
	Normal ColorDefinitionWithCode `json:"normal"`
	Bright ColorDefinitionWithCode `json:"bright"`
}

type PaletteVariant struct {
	Name       string                     `json:"name"`
	Emoji      string                     `json:"emoji"`
	Order      int                        `json:"order"`
	Dark       bool                       `json:"dark"`
	Colors     map[string]ColorDefinition `json:"colors"`
	AnsiColors map[string]ANSI            `json:"ansiColors"`
}

type Palette struct {
	Version  string                    `json:"version"`
	Variants map[string]PaletteVariant `json:"variants"`
}

type PaletteBuilder struct {
	palette Palette
}
