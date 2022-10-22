package main

import (
	"image/color"
	"io"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}
func init() {
	f, err := os.Open("Fonts/CascadiaCode-2111.01/ttf/CascadiaCode.ttf")
	check(err)
	defer f.Close()
	font_bytes, err := io.ReadAll(f)
	check(err)

	MainFont, err = truetype.Parse(font_bytes)
	check(err)
	MainFontOpts := truetype.Options{
		Size:              float64(MainFontSize),
		DPI:               0,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}
	MainFontFace = truetype.NewFace(MainFont, &MainFontOpts)
	MainFontDescent = MainFontFace.Metrics().Descent.Round()

	//menu bar font
	f2, err := os.Open("Fonts/CascadiaCode-2111.01/ttf/CascadiaCode.ttf")
	check(err)
	defer f2.Close()
	font_bytes2, err := io.ReadAll(f2)
	check(err)

	MenuFont, err = truetype.Parse(font_bytes2)
	check(err)
	MenuFontOpts := truetype.Options{
		Size:              float64(MenuFontSize),
		DPI:               0,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}
	MenuFontFace = truetype.NewFace(MenuFont, &MenuFontOpts)
	MenuFontDescent = MenuFontFace.Metrics().Descent.Round()
}

var MainFontSize int = 18
var MainFont *truetype.Font
var MainFontFace font.Face
var MainFontDescent int

var MenuFontSize int = 12
var MenuFont *truetype.Font
var MenuFontFace font.Face
var MenuFontDescent int

var tab_x_padding int = 13
var tab_y_padding int = 4

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
	FGColorMuted:  ParseHexColor("#FBF1C7"),
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
}

// https://stackoverflow.com/questions/54197913/parse-hex-string-to-image-color
func ParseHexColor(s string) (c color.RGBA) {
	c.A = 0xff

	if s[0] != '#' {
		return color.RGBA{0, 0, 0, 0}
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		return 0
	}

	switch len(s) {
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
	case 4:
		c.R = hexToByte(s[1]) * 17
		c.G = hexToByte(s[2]) * 17
		c.B = hexToByte(s[3]) * 17
	}
	return
}
