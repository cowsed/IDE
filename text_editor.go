package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
)

var _ Widget = &TextEditor{}

type Cursor struct {
	row, col int
}
type TextEditor struct {
	image.Rectangle
	text               []string
	text_tex           *ebiten.Image
	cursor             Cursor
	scroll             float64
	last_interact_time uint64 //tick alue of the last time we interacted (used to keep cursor alive while we're editing)
	ReadOnly           bool   //can we edit this textbox
	focused            bool   //does this textbox have keyboard focus
	saved              bool   //is the file saved to disk
	uptodate           bool

	highlighter *Highlighter

	filepath string
	filename string
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
	width := font.MeasureString(CodeFontFace, te.text[te.cursor.row][:te.cursor.col]).Round()
	if ((ticks-te.last_interact_time)/40)%2 == 0 {
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
	//draw background
	te.text_tex.Fill(color.RGBA{})
	if te.highlighter != nil {
		te.DrawWithHighlighting()
	} else {
		//log.Println("uncool drawing")
		//no ability to draw with syntax highlighting
		text.Draw(te.text_tex, strings.Join(te.text, "\n"), CodeFontFace, 0, text_edit_top_padding+CodeFontPeriodFromTop, Style.FGColorMuted) //@optimmize iterate through text and draw each line individuall. saves time by skipping the join

	}
}
func (te *TextEditor) DrawWithHighlighting() {
	colors := map[string]color.Color{
		"red":         Style.RedMuted,
		"brightred":   Style.RedStrong,
		"blue":        Style.BlueMuted,
		"brightblue":  Style.BlueStrong,
		"green":       Style.GreenMuted,
		"brightgreen": Style.GreenStrong,
		"yellow":      Style.YellowMuted,
		"cyan":        Style.AquaStrong,
		"magenta":     Style.PurpleMuted,
		"brightblack": Style.Gray,
	}
	topleft := image.Pt(0, 0) //top left of the line
	for _, line := range te.text {
		lineusage := make([]bool, len(line))
		use := func(start, end int) {
			for i := max(0, start); i < min(len(lineusage), end); i++ {
				lineusage[i] = true
			}
		}
		already_used := func(start, end int) bool {
			for i := max(0, start); i < min(len(lineusage), end); i++ {
				if lineusage[i] {
					return true
				}
			}
			return false
		}
		print_usage := func() {
			s := ""
			for _, u := range lineusage {
				if u {
					s += "1"
				} else {
					s += "0"
				}
			}
			fmt.Println(s)
		}
		//draw line in fgcolor as the back, handles any parts that aren't a special color
		text.Draw(te.text_tex, line, CodeFontFace, topleft.X, text_edit_top_padding+CodeFontPeriodFromTop+topleft.Y, Style.FGColorMuted)
		//for each highlighter regex, draw all that it can
		if te.highlighter.expressions != nil && len(te.highlighter.expressions) > 0 {

			for i := len(te.highlighter.expressions) - 1; i >= 0; i-- {
				exp := te.highlighter.expressions[i]
				//for _, exp := range te.highlighter.expressions {
				fg_col, col_exists := colors[exp.fg_col]
				if !col_exists {
					fg_col = colornames.Greenyellow
				}

				indices := exp.reg.FindAllStringIndex(line, -1)
				for _, startnend := range indices {
					start := startnend[0]
					end := startnend[len(startnend)-1]
					if !already_used(start, end) {
						fmt.Println("didnt overlaped")

						advance := font.MeasureString(CodeFontFace, line[:start]).Round()
						text.Draw(te.text_tex, line[start:end], CodeFontFace, topleft.X+advance, text_edit_top_padding+CodeFontPeriodFromTop+topleft.Y, fg_col)
						use(start, end)
					}

					//before_width := font.BoundString(CodeFontFace, line[start:end]).Dx()
				}
			}
			print_usage()
			fmt.Println(line)
		}
		topleft.Y += CodeFontSize
	}
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
