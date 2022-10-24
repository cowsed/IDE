package main

import (
	"image"
	"image/color"
	"log"
	"strings"

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
	ReadOnly           bool
	text               []string
	text_tex           *ebiten.Image
	cursor             Cursor
	uptodate           bool
	scroll             float64
	last_interact_time uint64
	focused            bool
	saved              bool
	filepath           string
	filename           string
}

// Title implements Widget
func (te *TextEditor) Title() string {
	if te.filepath == "" {
		return "untitled*"
	} else {
		s := te.filename
		if !te.saved {
			s += "*"
		}
		return s
	}
}

func (te *TextEditor) KeyboardFocusLost() {
	te.focused = false
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
	if te.ReadOnly {
		return
	}
	te.DrawCursor(target)
	if te.scroll > 0 {
		//Draw "shadow" from the top
		y := te.Rectangle.Min.Y
		for i := 1; i < 10; i++ {
			y++
			col := color.RGBA{
				A: 70 - uint8(i*5),
			}
			ebitenutil.DrawLine(target, float64(te.Min.X), float64(y), float64(te.Max.X), float64(y), col)
		}
	}
}
func (te *TextEditor) DrawCursor(target *ebiten.Image) {
	if !te.focused {
		return
	}
	y := te.cursor.row * CodeFontSize
	start := te.Min
	width := text.BoundString(CodeFontFace, te.text[te.cursor.row][:te.cursor.col]).Dx()

	num_trailing_spaces := 0
	consumed_whole := true
	for i := len(te.text[te.cursor.row][:te.cursor.col]) - 1; i >= 0; i-- {
		if te.text[te.cursor.row][:te.cursor.col][i:i+1] == " " {
			num_trailing_spaces++
		} else {
			consumed_whole = false
			break
		}
	}

	num_leading_spaces := 0 //we only count leading spaces if we didnt count them when counting trailing spaces (in a string of all spaces, if we didnt do this we would count all the spaces twice)
	if !consumed_whole {
		for _, c := range te.text[te.cursor.row][:te.cursor.col] {
			if c == rune(" "[0]) {
				num_leading_spaces++
			} else {
				break
			}
		}
	}
	if ((ticks-te.last_interact_time)/40)%2 == 0 {
		space_width, _ := CodeFontFace.GlyphAdvance(' ')
		width += (num_trailing_spaces + num_leading_spaces) * space_width.Round()
		move_over := 1
		width += move_over
		ebitenutil.DrawLine(target, float64(start.X+width), float64(start.Y+y), float64(start.X+width), float64(start.Y+y+CodeFontSize), Style.FGColorMuted)
	}
}

func (te *TextEditor) DrawTextTexture() {

	needed_dims := text.BoundString(CodeFontFace, strings.Join(te.text, "\n"))
	needed_dims.Max.Y += CodeFontPeriodFromTop
	needed_dims.Max.X = max(needed_dims.Max.X, 1)
	if te.text_tex == nil || (te.Rectangle.Dx() != te.text_tex.Bounds().Dx() || te.Rectangle.Dy() != te.text_tex.Bounds().Dy()) {
		te.text_tex = ebiten.NewImage(te.Dx(), te.Dy())
	}
	te.text_tex.Fill(color.RGBA{})
	text.Draw(te.text_tex, strings.Join(te.text, "\n"), CodeFontFace, 0, text_edit_top_padding+CodeFontPeriodFromTop, Style.FGColorMuted) //@optimmize iterate through text and draw each line individuall. saves time by skipping the join

}
func (te *TextEditor) MarkRedraw() {
	te.uptodate = false
}
func (te *TextEditor) EnterText(s string) {
	te.Interacted()
	before := te.text[te.cursor.row][:te.cursor.col]
	after := te.text[te.cursor.row][te.cursor.col:]
	te.text[te.cursor.row] = before + s + after
	te.cursor.col += len(s)
	te.MarkRedraw()
}

func (te *TextEditor) Backspace() {
	te.Interacted()

	if te.cursor.col == 0 && te.cursor.row == 0 {
		return
	}
	if te.cursor.col == 0 {
		var prev_line_len = 0
		if te.cursor.row >= 1 {
			prev_line_len = len(te.text[te.cursor.row-1])
		}
		//combine this line with previous
		this_line := te.text[te.cursor.row]
		up_to := te.text[0:te.cursor.row]
		after := []string{}
		if te.cursor.row+1 < len(te.text) {
			after = te.text[te.cursor.row+1:]
		}
		te.text = append(up_to, after...)
		te.text[te.cursor.row-1] += this_line
		te.cursor.row--
		te.cursor.col = prev_line_len
		te.MarkRedraw()
		return
	}
	//just change this line
	//make sure we're not out in left field
	te.cursor.col = min(len(te.text[te.cursor.row]), te.cursor.col)
	line := te.text[te.cursor.row]
	before := line[:te.cursor.col-1]
	after := ""
	if te.cursor.col < len(line) {
		after = line[te.cursor.col:]
	}
	line = before + after
	te.text[te.cursor.row] = line
	te.cursor.col--
	te.MarkRedraw()
	te.Interacted()
}
func (te *TextEditor) CursorLeft() {
	te.Interacted()
	if te.cursor.col == 0 && te.cursor.row == 0 {
		return
	}
	if te.cursor.col == 0 {
		prev_line_end := len(te.text[te.cursor.row-1])
		te.cursor.row--
		te.cursor.col = prev_line_end
		return
	}
	te.cursor.col--
}
func (te *TextEditor) CursorRight() {
	te.Interacted()
	//if at the end of a line
	if te.cursor.col > len(te.text[te.cursor.row])-1 {
		//if at the end of the file, cant go to the next line
		if te.cursor.row >= len(te.text)-1 {
			return
		}
		te.cursor.row++
		te.cursor.col = 0
		return
	}
	//just go right
	if te.cursor.col < len(te.text[te.cursor.row]) {
		te.cursor.col++
	}

}
func (te *TextEditor) CursorDown() {
	te.Interacted()

	//already at the bottom of the file
	if te.cursor.row >= len(te.text)-1 {
		te.cursor.col = len(te.text[te.cursor.row])
		return
	}
	te.cursor.row++
	te.cursor.col = min(te.cursor.col, len(te.text[te.cursor.row]))

}
func (te *TextEditor) CursorUp() {
	te.Interacted()
	//already at the top of the file
	if te.cursor.row <= 0 {
		te.cursor.col = 0
		return
	}
	te.cursor.row--
	te.cursor.col = min(te.cursor.col, len(te.text[te.cursor.row]))
}
func (te *TextEditor) Newline() {
	te.Interacted()

	line := te.text[te.cursor.row]
	line_before := line[:te.cursor.col]
	line_after := line[te.cursor.col:]

	lines_before := te.text[:te.cursor.row] //up to line with cursor
	lines_after := []string{}
	if te.cursor.row+1 < len(te.text) {
		lines_after = te.text[te.cursor.row+1:] //the rest starting at the line after where the cursor is
	}
	newtext := make([]string, len(lines_before), len(te.text)+1)
	copy(newtext, lines_before)
	newtext = append(newtext, line_before, line_after)
	newtext = append(newtext, lines_after...)
	te.text = newtext
	te.cursor.row++
	te.cursor.col = 0
}
func (te *TextEditor) SetText(s string) {
	te.text = strings.Split(s, "\n")
}

func (te *TextEditor) handle_shortcuts() {
	local_shortcuts := map[KeyShortcut]func(){
		{key: ebiten.KeyEnd}:               te.EndLine,
		{key: ebiten.KeyHome}:              te.StartLine,
		{key: ebiten.KeyBackspace}:         te.Backspace,
		{key: ebiten.KeyTab}:               te.Tab,
		{key: ebiten.KeyEnter}:             te.Newline,
		{key: ebiten.KeyLeft}:              te.CursorLeft,
		{key: ebiten.KeyRight}:             te.CursorRight,
		{key: ebiten.KeyUp}:                te.CursorUp,
		{key: ebiten.KeyDown}:              te.CursorDown,
		{mod_ctrl: true, key: ebiten.KeyA}: te.SelectAll,
	}
	ctrl_state := ebiten.IsKeyPressed(ebiten.KeyControl)
	shift_state := ebiten.IsKeyPressed(ebiten.KeyShift)
	alt_state := ebiten.IsKeyPressed(ebiten.KeyAlt)
	meta_state := ebiten.IsKeyPressed(ebiten.KeyMeta)

	for shortcut := range local_shortcuts {
		if KeyJustPressedOrKeyRepeated(shortcut.key) {
			executable := local_shortcuts[KeyShortcut{
				mod_shift: shift_state,
				mod_ctrl:  ctrl_state,
				mod_alt:   alt_state,
				mod_meta:  meta_state,
				key:       shortcut.key,
			}]
			if executable != nil {
				executable()
			}
		}
	}
}
func (te *TextEditor) TakeKeyboard() {
	te.handle_shortcuts()

	if te.ReadOnly {
		return
	}

	if ebiten.IsKeyPressed(ebiten.KeyControl) || ebiten.IsKeyPressed(ebiten.KeyAlt) || ebiten.IsKeyPressed(ebiten.KeyMeta) {
		return
	}
	var b []rune
	b = ebiten.AppendInputChars(b[:0])
	if len(b) > 0 {
		te.EnterText(string(b))
	}

}
func (te *TextEditor) Interacted() {
	te.last_interact_time = ticks
}

func (te *TextEditor) LMouseDown(x int, y int) Widget {
	te.focused = true
	//local_y := y - text_edit_top_padding - te.Min.Y
	return te
}

// LMouseUp implements Widget
func (te *TextEditor) LMouseUp(x int, y int) Widget {
	te.focused = true
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

/*
Shortcut functions
Press the key combo and do these common actions
*/
func (te *TextEditor) Tab() {
	te.EnterText("    ")
}
func (te *TextEditor) SelectAll() {
	log.Println("Selectall unimplemented")
	te.Interacted()
}
func (te *TextEditor) EndLine() {
	te.cursor.col = len(te.text[te.cursor.row])
	te.Interacted()
}
func (te *TextEditor) StartLine() {
	te.cursor.col = 0
	te.Interacted()
}
