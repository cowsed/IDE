package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/hajimehoshi/ebiten/v2"
)

type Widget interface {
	Draw(target *ebiten.Image)
	SetRect(rect image.Rectangle)
	//Keyboard events get sent to the place that last received mouse input
	TakeKeyboard(key ebiten.Key)
	//Mouse events return a pointer(interfaces are just pointers) to the widget that actually used the input
	//this is used to tell that widget when a mouse out happens
	MouseOut()
	MouseOver(x, y int) Widget
	LMouseDown(x, y int) Widget
	LMouseUp(x, y int) Widget
}

var _ Widget = &ColorRect{}
var _ Widget = &Tabs{}
var _ Widget = &HorizontalSplitter{}
var _ Widget = &MenuBar{}

type BorderShowMode int

const (
	ShowAlways BorderShowMode = iota
	ShowOnHover
)

type Tabs struct {
	image.Rectangle
	Titles          []string
	Tabs            []Widget
	TabHeaderRects  []image.Rectangle
	CurrentTab      int
	TabHeight       int
	current_hovered int
}

// TakeKeyboard implements Widget
func (t *Tabs) TakeKeyboard(key ebiten.Key) {
	log.Println("TakeKeyboard unimplemented for tabs")
}

// MouseOut implements Widget
func (t *Tabs) MouseOut() {
	t.current_hovered = -1
}

func (t *Tabs) Draw(target *ebiten.Image) {
	t.DrawTabs(target)
	//for a myriad of reasons we can't draw the current tab
	if t.CurrentTab < 0 || t.CurrentTab > len(t.Tabs) || t.Tabs[t.CurrentTab] == nil {
		return
	}
	t.Tabs[t.CurrentTab].Draw(target)
}
func (t *Tabs) DrawTabs(target *ebiten.Image) {

	for i := 0; i < len(t.Tabs); i++ {

		my_r := t.TabHeaderRects[i]
		my_c := Style.BGColorMuted
		if i == t.current_hovered {
			my_c = Style.BGColorStrong
		}

		ebitenutil.DrawRect(target, float64(my_r.Min.X), float64(my_r.Min.Y), float64(my_r.Dx()), float64(my_r.Dy()), my_c)
		text.Draw(target, t.Titles[i], MainFontFace, my_r.Min.X+tab_x_padding, my_r.Min.Y+MainFontPeriodFromTop+tab_y_padding, Style.FGColorStrong)
		if i == t.CurrentTab {
			ebitenutil.DrawLine(target, float64(my_r.Min.X), float64(my_r.Max.Y-1), float64(my_r.Max.X), float64(my_r.Max.Y-1), Style.FGColorStrong)
		}

	}
}

func (t *Tabs) LMouseDown(x int, y int) Widget {
	// over tabs
	if y < t.Rectangle.Min.Y+t.TabHeight {
		for i := 0; i < len(t.TabHeaderRects); i++ {
			if image.Pt(x, y).In(t.TabHeaderRects[i]) {
				fmt.Println("tab switched")
				t.CurrentTab = i
			}
		}
		return t
	}
	// over body
	if t.Tabs[t.CurrentTab] != nil {
		return t.Tabs[t.CurrentTab].LMouseDown(x, y)
	}
	return nil
}

func (t *Tabs) LMouseUp(x int, y int) Widget {
	// over tabs
	if x < t.Rectangle.Min.X+t.TabHeight {
		return t
	}
	//over body
	if t.Tabs[t.CurrentTab] != nil {
		return t.Tabs[t.CurrentTab].LMouseUp(x, y)
	}
	return nil
}

func (t *Tabs) MouseOver(x int, y int) Widget {
	// over tabs
	if y < t.Rectangle.Min.Y+t.TabHeight {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
		for i, r := range t.TabHeaderRects {
			if image.Pt(x, y).In(r) {
				t.current_hovered = i
			}
		}
		return t
	}
	//over body
	if t.Tabs[t.CurrentTab] != nil {
		return t.Tabs[t.CurrentTab].MouseOver(x, y)
	}
	return nil
}

// SetRect implements Widget
func (t *Tabs) SetRect(rect image.Rectangle) {
	t.Rectangle = rect
	tabs_rect := rect
	tabs_rect.Max.Y = rect.Min.Y + t.TabHeight

	if len(t.TabHeaderRects) != len(t.Tabs) {
		t.TabHeaderRects = make([]image.Rectangle, len(t.Tabs))
	}
	//Generate rects describing tabs
	start := tabs_rect.Min
	for i := 0; i < len(t.Tabs); i++ {
		//fmt.Println("----")
		//fmt.Println(t.Titles[i])
		text_width := text.BoundString(MainFontFace, t.Titles[i]).Dx()
		//fmt.Println(text_width)

		tab_rect := image.Rect(0, 0, text_width+2*tab_x_padding, MainFontSize+2*tab_y_padding)

		t.TabHeaderRects[i] = tab_rect.Add(start)

		start = start.Add(image.Pt(tab_rect.Dx(), 0))

	}
	//rect for content
	pane_rect := rect
	pane_rect.Min.Y = rect.Min.Y + t.TabHeight
	for i := 0; i < len(t.Tabs); i++ {
		if t.Tabs[i] != nil {
			t.Tabs[i].SetRect(pane_rect)
		}
	}
}

type HorizontalSplitter struct {
	image.Rectangle
	split_x           int
	Left, Right       Widget
	dragging          bool
	border_half_width int
	border_mode       BorderShowMode
	border_hovered    bool
}

// TakeKeyboard implements Widget
func (hz *HorizontalSplitter) TakeKeyboard(key ebiten.Key) {
	log.Println("TakeKeyboard unimplemented for horizontal splitter")
	//panic("unimplemented")
}

// MouseOut implements Widget
func (*HorizontalSplitter) MouseOut() {
	fmt.Println("lost mouse")
}

func (hz *HorizontalSplitter) LMouseUp(x, y int) Widget {
	if hz.dragging {
		hz.dragging = false //cant be dragging if we let go
		fmt.Print("stoped draggin\n")
	} else if x < hz.split_x-hz.border_half_width {
		if hz.Left != nil {
			return hz.Left.LMouseUp(x, y)
		}
	} else if x > hz.split_x+hz.border_half_width {
		if hz.Right != nil {
			return hz.Right.LMouseUp(x-hz.split_x, y)
		}
	}
	return hz
}

func (hz *HorizontalSplitter) LMouseDown(x, y int) Widget {
	if x > hz.split_x-hz.border_half_width && x < hz.split_x+hz.border_half_width {
		hz.dragging = true

	}
	if x < hz.split_x-hz.border_half_width {
		//left
		if hz.Left != nil {
			return hz.Left.LMouseDown(x, y)
		}
	}
	if x > hz.split_x+hz.border_half_width {
		//left
		if hz.Right != nil {
			return hz.Right.LMouseDown(x, y)
		}
	}
	return hz
}

func (hz *HorizontalSplitter) MouseOver(x, y int) Widget {
	screenspaceDividerX := hz.split_x + hz.Rectangle.Min.X
	var consumer Widget = nil
	//pass mouse down to the left
	if x < screenspaceDividerX-hz.border_half_width {
		hz.border_hovered = false
		if hz.Left != nil {
			consumer = hz.Left.MouseOver(x, y)
		} else {
			ebiten.SetCursorShape(ebiten.CursorShapeDefault)
		}
	} else if x > screenspaceDividerX+hz.border_half_width {
		//pass mouse down to the right
		hz.border_hovered = false
		if hz.Right != nil {
			consumer = hz.Right.MouseOver(x, y)
		} else {
			ebiten.SetCursorShape(ebiten.CursorShapeDefault)
		}
	} else {
		ebiten.SetCursorShape(ebiten.CursorShapeEWResize)
		hz.border_hovered = true
		consumer = hz

	}
	if hz.dragging {
		fmt.Println("Draggin")
		ebiten.SetCursorShape(ebiten.CursorShapeEWResize)
		hz.split_x = x
		//stop from going too far that you can't reach the handle
		if hz.split_x < hz.Min.X+5 {
			hz.split_x = hz.Min.X + 5
		} else if hz.split_x > hz.Max.X-5 {
			hz.split_x = hz.Max.X - 5
		}
		hz.SetRect(hz.Rectangle)
		consumer = hz

	}
	return consumer
}

func (hz *HorizontalSplitter) Draw(target *ebiten.Image) {
	if hz.Left != nil {
		hz.Left.Draw(target)
	}
	if hz.Right != nil {
		hz.Right.Draw(target)
	} else {
		println("no right")
	}
	//draw divider

	if hz.border_mode == ShowAlways || hz.border_hovered || hz.dragging {
		border_min_x := hz.Rectangle.Min.X + hz.split_x - hz.border_half_width
		ebitenutil.DrawRect(target, float64(border_min_x), float64(hz.Rectangle.Min.Y), float64(hz.border_half_width)*2, float64(hz.Rectangle.Dy()), Style.FGColorMuted)
	}

}
func (hz *HorizontalSplitter) SetRect(r image.Rectangle) {
	fmt.Println("recalc hz")

	old_width := hz.Rectangle.Dx()
	hz.Rectangle = r
	if old_width == 0 {
		fmt.Println("skupping recalc")
		return
	}
	x_percent := float64(hz.split_x) / float64(old_width)

	hz.split_x = int(x_percent * float64(r.Dx()))
	fmt.Println(hz.split_x, x_percent)

	left_rect_min_x := hz.Rectangle.Min.X
	left_rect_max_x := hz.Rectangle.Min.X + hz.split_x  // - hz.border_half_width
	right_rect_min_x := hz.Rectangle.Min.X + hz.split_x // + hz.border_half_width
	right_rect_max_x := hz.Rectangle.Max.X

	leftRect := image.Rect(left_rect_min_x, hz.Min.Y, left_rect_max_x, hz.Max.Y)
	rightRect := image.Rect(right_rect_min_x, hz.Min.Y, right_rect_max_x, hz.Max.Y)

	if hz.Left != nil {
		hz.Left.SetRect(leftRect)
	}
	if hz.Right != nil {
		hz.Right.SetRect(rightRect)
	}
}

func NewColorRect(col color.Color) *ColorRect {
	return &ColorRect{color: col}
}

type ColorRect struct {
	image.Rectangle
	color color.Color
}

// TakeKeyboard implements Widget
func (cr *ColorRect) TakeKeyboard(key ebiten.Key) {
	log.Println("TakeKeyboard unimplemented for color rect")
	//panic("unimplemented")
}

// MouseOut implements Widget
func (*ColorRect) MouseOut() {
}

// Draw implements Widget
func (cr *ColorRect) Draw(target *ebiten.Image) {
	ebitenutil.DrawRect(target, float64(cr.Min.X), float64(cr.Min.Y), float64(cr.Dx()), float64(cr.Dy()), cr.color)
}

// LMouseDown implements Widget
func (cr *ColorRect) LMouseDown(x int, y int) Widget {
	return cr
}

// LMouseUp implements Widget
func (cr *ColorRect) LMouseUp(x int, y int) Widget {
	return cr
}

// MouseOver implements Widget
func (cr *ColorRect) MouseOver(x int, y int) Widget {
	ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	return cr
}

func (cr *ColorRect) SetRect(rect image.Rectangle) {
	cr.Rectangle = rect

}
