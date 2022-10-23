package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func text_stats(lines []string) (int, int) {
	//returns largest line (x) and num lines (y)
	max_width := 0

	for _, line := range lines {
		max_width = max(max_width, len(line))
	}
	return len(lines), max_width
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func DrawRect(target *ebiten.Image, r image.Rectangle, color color.Color) {
	ebitenutil.DrawRect(target, float64(r.Min.X), float64(r.Min.Y), float64(r.Dx()), float64(r.Dy()), color)
}
func TopLeft(r image.Rectangle) image.Point {
	return image.Pt(r.Min.X, r.Min.Y)
}
func TopRight(r image.Rectangle) image.Point {
	return image.Pt(r.Max.X, r.Min.Y)
}
func BottomLeft(r image.Rectangle) image.Point {
	return image.Pt(r.Min.X, r.Max.Y)
}
func BottomRight(r image.Rectangle) image.Point {
	return image.Pt(r.Max.X, r.Max.Y)
}
func check(err error) {
	if err != nil {
		panic(err)
	}
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

func DrawBorders(target *ebiten.Image, r image.Rectangle, color color.Color) {
	tl, tr, bl, br := TopLeft(r), TopRight(r), BottomLeft(r), BottomRight(r)
	//top
	ebitenutil.DrawLine(target, float64(tl.X), float64(tl.Y), float64(tr.X), float64(tr.Y), color)
	//bottom
	ebitenutil.DrawLine(target, float64(bl.X), float64(bl.Y), float64(br.X), float64(br.Y), color)
	//left
	ebitenutil.DrawLine(target, float64(tl.X+1), float64(tl.Y), float64(bl.X+1), float64(bl.Y), color)
	//right
	ebitenutil.DrawLine(target, float64(tr.X), float64(tr.Y), float64(br.X), float64(br.Y), color)
}
