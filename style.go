package main

import (
	"image/color"
	"io"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

func init() {
	MainFont = LoadFont("Fonts/Source_Code_Pro/SourceCodePro-Regular.ttf")
	MenuFont = LoadFont("Fonts/Source_Code_Pro/SourceCodePro-Regular.ttf")
	CodeFont = LoadFont("Fonts/Source_Code_Pro/SourceCodePro-Regular.ttf")

	MainFontFace, MainFontPeriodFromTop = MakeFace(MainFont, MainFontSize)
	MenuFontFace, MenuFontPeriodFromTop = MakeFace(MenuFont, MenuFontSize)
	CodeFontFace, CodeFontPeriodFromTop = MakeFace(CodeFont, CodeFontSize)
}
func LoadFont(path string) (font *truetype.Font) {
	f, err := os.Open(path)
	check(err)
	defer f.Close()
	font_bytes, err := io.ReadAll(f)
	check(err)
	font, err = truetype.Parse(font_bytes)
	check(err)

	return font
}
func MakeFace(font *truetype.Font, size int) (font.Face, int) {
	FontOpts := truetype.Options{
		Size:              float64(size),
		DPI:               0,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}
	face := truetype.NewFace(font, &FontOpts)
	PeriodFromTop := face.Metrics().Height.Round() - face.Metrics().Descent.Round()
	return face, PeriodFromTop
}

const MinFontSize = 6

var MainFontSize int = 16
var MainFont *truetype.Font
var MainFontFace font.Face
var MainFontPeriodFromTop int

var CodeFontSize int = 16
var CodeFont *truetype.Font
var CodeFontFace font.Face
var CodeFontPeriodFromTop int

var MenuFontSize int = 14
var MenuFont *truetype.Font
var MenuFontFace font.Face
var MenuFontPeriodFromTop int

var text_edit_top_padding = 3

var tab_x_padding int = 13
var tab_y_padding int = 8

var menu_bar_x_padding int = 8
var menu_bar_y_padding int = 4

var menu_x_padding int = 10
var menu_y_padding int = 10

var BGColor0Hard = ParseHexColor("#1D2021")
var FGColor0Hard = ParseHexColor("#FBF1C7")

var BGColor0Soft = ParseHexColor("#32302F")
var FGColor0Soft = ParseHexColor("#FBF1C7")

// https://www.reddit.com/r/gruvbox/comments/np5ylp/official_resources/
var Style = StyleColors{
	BGColorStrong: ParseHexColor("#1D2019"),
	FGColorStrong: ParseHexColor("#FBF1C7"),
	BGColorMuted:  ParseHexColor("#32302F"),
	FGColorMuted:  ParseHexColor("#BDAE93"),
	RedStrong:     ParseHexColor("#FB4934"),
	RedMuted:      ParseHexColor("#CC241D"),
	GreenStrong:   ParseHexColor("#B8BB26"),
	GreenMuted:    ParseHexColor("#98971A"),
	YellowStrong:  ParseHexColor("#FABD2F"),
	YellowMuted:   ParseHexColor("#D79921"),
	BlueStrong:    ParseHexColor("#83A598"),
	BlueMuted:     ParseHexColor("#458588"),
	PurpleStrong:  ParseHexColor("#D3869B"),
	PurpleMuted:   ParseHexColor("#B16286"),
	AquaStrong:    ParseHexColor("#8EC07C"),
	AquaMuted:     ParseHexColor("#689D6A"),
	OrangeStrong:  ParseHexColor("#FE8019"),
	OrangeMuted:   ParseHexColor("#D65D0E"),
	Gray:          ParseHexColor("#a89984"),
	White:         ParseHexColor("ebdbb2"),
}

type StyleColors struct {
	BGColorStrong color.Color
	FGColorStrong color.Color

	BGColorMuted color.Color
	FGColorMuted color.Color

	RedStrong color.Color
	RedMuted  color.Color

	GreenStrong color.Color
	GreenMuted  color.Color

	YellowStrong color.Color
	YellowMuted  color.Color

	BlueStrong color.Color
	BlueMuted  color.Color

	PurpleStrong color.Color
	PurpleMuted  color.Color

	AquaStrong color.Color
	AquaMuted  color.Color

	OrangeStrong color.Color
	OrangeMuted  color.Color

	White color.Color
	Gray  color.Color
}
