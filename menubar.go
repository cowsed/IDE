package main

import (
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type KeyShortcut struct {
	mods []ebiten.Key
	key  ebiten.Key
}

func (ks *KeyShortcut) String() string {
	s := ""
	for i := range ks.mods {
		s += ks.mods[i].String()
		s += " + "
	}
	s += ks.key.String()
	return s
}

type MenuItem interface {
	Text() string
	Children() []MenuItem
	Execute()
	Shortcut() KeyShortcut
	DrawOpen(target *ebiten.Image, topleft image.Point)
	SpaceUsed(topleft image.Point) []image.Rectangle
	MouseOver(x, y int)
}

var _ MenuItem = &DummyMenuItem{txt: "File"}

func NewMenuItem(name string, children []MenuItem) *DummyMenuItem {
	return &DummyMenuItem{
		txt:               name,
		currently_hovered: -1,
		width:             0,
		kids:              children,
		ks:                KeyShortcut{},
	}
}

type DummyMenuItem struct {
	txt               string
	currently_hovered int
	width             int
	kids              []MenuItem
	itemrects         []image.Rectangle
	ks                KeyShortcut
}

func (dmi *DummyMenuItem) MouseOver(x, y int) {
	for i, r := range dmi.itemrects {
		if image.Pt(x, y).In(r) {
			dmi.currently_hovered = i
		}
	}
	if dmi.currently_hovered != -1 {
		dmi.kids[dmi.currently_hovered].MouseOver(x, y)
	}
}
func (dmi *DummyMenuItem) Shortcut() KeyShortcut {
	return dmi.ks
}
func (dmi *DummyMenuItem) Children() []MenuItem {
	return []MenuItem{}
}

// Calculates the size of this menu if it were drawn
func (dmi *DummyMenuItem) SpaceUsed(topleft image.Point) []image.Rectangle {
	if len(dmi.kids) == 0 {
		return []image.Rectangle{}
	}
	biggest_width := 0
	height_needed := menu_y_padding
	text_rect_tops := make([]int, len(dmi.kids))
	dmi.itemrects = make([]image.Rectangle, len(dmi.kids))
	for i := range dmi.kids {
		text_rect_tops[i] = height_needed
		text_r := text.BoundString(MenuFontFace, dmi.kids[i].Text())
		biggest_width = max(biggest_width, text_r.Dx())
		height_needed += text_r.Dy() + menu_y_padding
	}
	biggest_width += menu_bar_x_padding * 2
	dmi.width = biggest_width
	box_h := MenuFontSize + menu_y_padding
	ascent := MenuFontFace.Metrics().Ascent.Round()

	for i, y := range text_rect_tops {
		dmi.itemrects[i] = image.Rect(topleft.X, y+menu_y_padding+ascent, topleft.X+biggest_width, y+box_h+menu_y_padding+ascent)
	}

	var child_rects = []image.Rectangle{}
	if dmi.currently_hovered >= 0 && dmi.currently_hovered < len(dmi.kids) && dmi.kids[dmi.currently_hovered] != nil {
		//collect space used by open children
		child_rects = dmi.kids[dmi.currently_hovered].SpaceUsed(image.Pt(biggest_width-1, dmi.currently_hovered*(MenuFontSize+menu_y_padding)).Add(topleft))
	}
	my_space := append(child_rects, image.Rect(topleft.X, topleft.Y, topleft.X+biggest_width, topleft.Y+height_needed))
	return my_space
}

// just draws the text/content of the menu, drawing the background box is taken care of in menubar.Draw
func (dmi *DummyMenuItem) DrawOpen(target *ebiten.Image, topleft image.Point) {
	if len(dmi.kids) == 0 {
		return
	}

	start := topleft
	start.X += menu_x_padding
	start.Y += MenuFontFace.Metrics().Ascent.Round() + menu_y_padding
	for i, mi := range dmi.kids {

		if dmi.currently_hovered == i {
			//draw this one brighter
			DrawRect(target, dmi.itemrects[i], Style.RedMuted)
			dmi.kids[dmi.currently_hovered].DrawOpen(target, image.Pt(start.X+dmi.width-menu_x_padding, start.Y-MenuFontFace.Metrics().Ascent.Round()-menu_y_padding))
		}
		text.Draw(target, mi.Text(), MenuFontFace, start.X, start.Y, Style.FGColorStrong)

		start.Y += MenuFontSize + menu_y_padding

	}

}

// Execute implements MenuItem
func (dmi *DummyMenuItem) Execute() {
	panic("unimplemented")
}

// Text implements MenuItem
func (dmi *DummyMenuItem) Text() string {
	return dmi.txt
}
func NewMenuBar(Items []MenuItem, SubWidget Widget) *MenuBar {
	return &MenuBar{
		currently_hovered: -1,
		currently_open:    -1,
		WidgetIApplyTo:    SubWidget,
		TopLevelItems:     Items,
	}
}

type MenuBar struct {
	image.Rectangle
	currently_hovered int
	currently_open    int
	TopLevelRects     []image.Rectangle
	TopLevelItems     []MenuItem

	WidgetIApplyTo Widget
}

// TakeKeyboard implements Widget
func (mb *MenuBar) TakeKeyboard(key ebiten.Key) {
	log.Println("TakeKeyboard unimplemented for menubar")
	//panic("unimplemented")
}

// MouseOut implements Widget
func (mb *MenuBar) MouseOut() {
	mb.currently_hovered = -1
	mb.currently_open = -1
}

// Draw implements Widget
func (mb *MenuBar) Draw(target *ebiten.Image) {
	//draw child widget
	if mb.WidgetIApplyTo != nil {
		mb.WidgetIApplyTo.Draw(target)
	}

	//Draw backgrounds of open menus, if they exist
	if mb.currently_open >= 0 {
		topleft_of_menu := BottomLeft(mb.TopLevelRects[mb.currently_open])
		bg_rects := mb.TopLevelItems[mb.currently_open].SpaceUsed(topleft_of_menu)
		for _, r := range bg_rects {
			ebitenutil.DrawRect(target, float64(r.Min.X), float64(r.Min.Y), float64(r.Dx()), float64(r.Dy()), Style.BGColorMuted)
			DrawBorders(target, r, Style.FGColorMuted)

		}

	}

	for i := 0; i < len(mb.TopLevelItems); i++ {

		my_r := mb.TopLevelRects[i]
		my_c := Style.BGColorMuted
		if i == mb.currently_hovered {
			my_c = Style.BGColorStrong
		}
		ebitenutil.DrawRect(target, float64(my_r.Min.X), float64(my_r.Min.Y), float64(my_r.Dx()), float64(my_r.Dy()), my_c)

		text.Draw(target, mb.TopLevelItems[i].Text(), MenuFontFace, my_r.Min.X+menu_bar_x_padding, my_r.Max.Y-menu_bar_y_padding-MenuFontDescent/2, Style.FGColorStrong)

		if i == mb.currently_open {
			mb.TopLevelItems[i].DrawOpen(target, BottomLeft(my_r))
		}

	}

}

// LMouseDown implements Widget
func (mb *MenuBar) LMouseDown(x int, y int) Widget {
	split_y_ss := mb.Rectangle.Min.Y + MenuFontSize + 2*menu_bar_y_padding
	if y < split_y_ss { // captured by the menu
		for i, r := range mb.TopLevelRects {
			if image.Pt(x, y).In(r) { //clicked onto menu item, open it
				mb.currently_open = i
			}
		}
		return mb
	}
	if mb.WidgetIApplyTo != nil {
		return mb.WidgetIApplyTo.LMouseDown(x, y)
	}
	return nil
}

// LMouseUp implements Widget
func (mb *MenuBar) LMouseUp(x int, y int) Widget {
	split_y_ss := mb.Rectangle.Min.Y + MenuFontSize + 2*menu_bar_y_padding
	if y < split_y_ss {
		return mb
	}
	if mb.WidgetIApplyTo != nil {
		return mb.WidgetIApplyTo.LMouseUp(x, y)
	}
	return nil
}

// MouseOver implements Widget
func (mb *MenuBar) MouseOver(x int, y int) Widget {
	ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	split_y_ss := mb.Rectangle.Min.Y + MenuFontSize + 2*menu_bar_y_padding
	if y < split_y_ss {
		for i, r := range mb.TopLevelRects {
			if image.Pt(x, y).In(r) {
				mb.currently_hovered = i
				//If we have a menu item open and move our mouse to the next menu item, show that one
				if mb.currently_open != -1 && mb.currently_hovered != mb.currently_open {
					mb.currently_open = mb.currently_hovered
				}
			}
		}
		return mb
	}
	//mouse in an open menu
	sub_menu_space := []image.Rectangle{}
	if mb.currently_open >= 0 {
		top_level_rect := mb.TopLevelRects[mb.currently_open]
		top_left_of_open_menu := image.Pt(top_level_rect.Min.X, top_level_rect.Max.Y)
		sub_menu_space = mb.TopLevelItems[mb.currently_open].SpaceUsed(top_left_of_open_menu)
	}
	for i := range sub_menu_space {
		if image.Pt(x, y).In(sub_menu_space[i]) {
			mb.TopLevelItems[mb.currently_open].MouseOver(x, y)

			return mb
		}
		//if we got here, mouse is out

	}
	if mb.WidgetIApplyTo != nil {
		return mb.WidgetIApplyTo.MouseOver(x, y)
	}
	return nil
}

// SetRect implements Widget
func (mb *MenuBar) SetRect(rect image.Rectangle) {
	//consume the top bit for me
	split_y := rect.Min.Y + MenuFontSize + 2*menu_bar_y_padding
	mb.Rectangle = image.Rectangle{
		Min: rect.Min,
		Max: image.Point{rect.Max.X, split_y},
	}
	//figure out item rect sizes
	if len(mb.TopLevelRects) != len(mb.TopLevelItems) {
		mb.TopLevelRects = make([]image.Rectangle, len(mb.TopLevelItems))
	}
	start_x := mb.Rectangle.Min.X
	y_start := mb.Min.Y
	y_end := mb.Max.Y
	for i, item := range mb.TopLevelItems {
		width_needed := text.BoundString(MenuFontFace, item.Text()).Dx() + 2*menu_bar_x_padding

		mb.TopLevelRects[i] = image.Rect(start_x, y_start, start_x+width_needed, y_end)
		start_x += width_needed
	}
	//give the rest to child
	if mb.WidgetIApplyTo != nil {
		child_rect := image.Rectangle{
			Min: image.Point{rect.Min.X, split_y},
			Max: rect.Max,
		}
		fmt.Println("Set child to", child_rect)
		mb.WidgetIApplyTo.SetRect(child_rect)
	}
}
