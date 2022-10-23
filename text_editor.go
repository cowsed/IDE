package main

import (
	"image"
	"image/color"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"
)

var _ Widget = &TextEditor{}

type Cursor struct {
	row, col int
}

type TextEditor struct {
	image.Rectangle
	text     []string
	text_tex *ebiten.Image
	cursor   Cursor
	uptodate bool
}

// Draw implements Widget
func (te *TextEditor) Draw(target *ebiten.Image) {
	ebitenutil.DrawRect(target, float64(te.Min.X), float64(te.Min.Y), float64(te.Dx()), float64(te.Dy()), Style.BGColorMuted)
	te.DrawTextTexture()
	//}
	geo := ebiten.GeoM{}
	geo.Translate(float64(te.Min.X), float64(te.Min.Y))
	target.DrawImage(te.text_tex, &ebiten.DrawImageOptions{
		GeoM:          geo,
		ColorM:        ebiten.ColorM{},
		CompositeMode: 0,
		Filter:        0,
	})

	//Draw "shadow" from the top
	for i := 1; i < 10; i++ {
		bounds := te.Rectangle.Inset(i)
		col := color.RGBA{
			A: 70 - uint8(i*5),
		}
		tl, tr := TopLeft(bounds), TopRight(bounds)
		ebitenutil.DrawLine(target, float64(tl.X), float64(tl.Y), float64(tr.X), float64(tr.Y), col)
	}

}

func (te *TextEditor) DrawTextTexture() {

	needed_dims := text.BoundString(CodeFontFace, strings.Join(te.text, "\n"))
	needed_dims.Max.Y += 11
	if te.text_tex == nil || needed_dims != te.Bounds() {
		te.text_tex = ebiten.NewImage(needed_dims.Dx(), needed_dims.Dy())

	}

	text.Draw(te.text_tex, strings.Join(te.text, "\n"), CodeFontFace, 0, 11, colornames.Navajowhite)

}
func (te *TextEditor) MarkRedraw() {
	te.uptodate = false
}
func (te *TextEditor) EnterText(s string) {
	te.text[te.cursor.row] += s
	te.MarkRedraw()
}

func (te *TextEditor) Backspace() {
	if te.cursor.col == 0 && te.cursor.row == 0 {
		return
	}

	te.MarkRedraw()
}
func (te *TextEditor) Newline() {
	//te.text[te.cursor.row:] = append("", te.text[te.cursor.row:])
	//te.text = append(append(te.text[0:te.cursor.row], ""), te.text[te.cursor.row:]...)
	newtext := make([]string, len(te.text)+1)
	old_index := 0
	new_index := 0
	for old_index < len(te.text) {

		newtext[new_index] = te.text[old_index]

		old_index++
		new_index++
	}
	te.text = newtext
	te.cursor.row++
}
func (te *TextEditor) TakeKeyboard(key ebiten.Key) {

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		te.Backspace()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		te.Newline()
	}

	var b []rune
	b = ebiten.AppendInputChars(b[:0])
	if len(b) > 0 {
		te.EnterText(string(b))
	}

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
