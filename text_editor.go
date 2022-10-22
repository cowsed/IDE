package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

var _ Widget = &TextEditor{}

type Cursor struct {
	row, col int
}

type TextEditor struct {
	image.Rectangle
	show_line_numbers bool
	cursors           []Cursor
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

// Draw implements Widget
func (te *TextEditor) Draw(target *ebiten.Image) {
	ebitenutil.DrawRect(target, float64(te.Min.X), float64(te.Min.Y), float64(te.Dx()), float64(te.Dy()), Style.BGColorMuted)
	text.Draw(target, "TEXT HERE", MainFontFace, te.Min.X+10, te.Min.Y+20, Style.FGColorStrong)

	for i := 1; i < 10; i++ {
		bounds := te.Rectangle.Inset(i)
		col := color.RGBA{
			A: 50 - uint8(i*5),
		}
		tl, tr := TopLeft(bounds), TopRight(bounds)
		ebitenutil.DrawLine(target, float64(tl.X), float64(tl.Y), float64(tr.X), float64(tr.Y), col)
		//Draw(target, bounds, col)
	}
	DrawBorders(target, te.Rectangle, Style.BGColorStrong)
}

// LMouseDown implements Widget
func (te *TextEditor) LMouseDown(x int, y int) Widget {
	log.Println("LMOUSEDOWN ON TEXT EDITOR UNIMPLEMENTED")
	//panic("unimplemented")
	return te
}

// LMouseUp implements Widget
func (te *TextEditor) LMouseUp(x int, y int) Widget {
	log.Println("LMOUSEUP ON TEXT EDITOR UNIMPLEMENTED")
	return te
}

// MouseOut implements Widget
func (*TextEditor) MouseOut() {
}

// MouseOver implements Widget
func (te *TextEditor) MouseOver(x int, y int) Widget {
	ebiten.SetCursorShape(ebiten.CursorShapeText)
	return te
}

// SetRect implements Widget
func (te *TextEditor) SetRect(rect image.Rectangle) {
	te.Rectangle = rect
}
