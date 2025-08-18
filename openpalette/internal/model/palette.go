package model

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

type ColorDefinition struct {
	Name   string `json:"name"`
	Order  int    `json:"order"`
	Hex    string `json:"hex"`
	RGB    RGB    `json:"rgb"`
	HSL    HSL    `json:"hsl"`
	Accent bool   `json:"accent"`
}

type AnsiColorDefinition struct {
	Name string `json:"name"`
	Hex  string `json:"hex"`
	RGB  RGB    `json:"rgb"`
	HSL  HSL    `json:"hsl"`
	Code int    `json:"code"`
}

type ANSI struct {
	Name   string              `json:"name"`
	Order  int                 `json:"order"`
	Normal AnsiColorDefinition `json:"normal"`
	Bright AnsiColorDefinition `json:"bright"`
}

type PaletteVariant struct {
	Name       string                     `json:"name"`
	Emoji      string                     `json:"emoji"`
	Order      int                        `json:"order"`
	Dark       bool                       `json:"dark"`
	Colors     map[string]ColorDefinition `json:"colors"`
	AnsiColors map[string]ANSI            `json:"ansiColors"`
}

type Palette map[string]PaletteVariant

type Root struct {
	Version string `json:"version"`
	Palette
}
