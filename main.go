package main

import (
	"errors"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var ticks uint64

type Editor struct {
	screenWidth  int
	screenHeight int

	should_close bool

	MainWidget             Widget
	last_mouse_consumer    Widget //widget to send mouseout to
	last_keyboard_consumer Widget //Widget to send keyboard inputs to
}

func ToggleFullscreen() {
	ebiten.SetFullscreen(!ebiten.IsFullscreen())
}
func (g *Editor) SetShouldClose() {
	g.should_close = true
}
func (g *Editor) Update() error {
	ticks++
	if g.should_close {
		return errors.New("editor closed by user")
	}
	if !ebiten.IsFocused() {
		return nil
	}

	//mouse handling
	x, y := ebiten.CursorPosition()
	mouse_consumer := g.MainWidget.MouseOver(x, y)
	if mouse_consumer != g.last_mouse_consumer {
		if g.last_mouse_consumer != nil {
			g.last_mouse_consumer.MouseOut()
		}
	}
	g.last_mouse_consumer = mouse_consumer

	key_consumer := g.last_keyboard_consumer
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		consumer := g.MainWidget.LMouseDown(x, y)
		if consumer != nil {
			key_consumer = consumer
		}
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		consumer := g.MainWidget.LMouseUp(x, y)
		if consumer != nil {
			key_consumer = consumer
		}

	}

	global_shortcuts := map[KeyShortcut]func(){
		{
			mod_shift: false,
			mod_ctrl:  false,
			mod_alt:   false,
			mod_meta:  false,
			key:       ebiten.KeyF11,
		}: ToggleFullscreen,
		{
			mod_shift: false,
			mod_ctrl:  true,
			mod_alt:   false,
			mod_meta:  false,
			key:       ebiten.KeyQ,
		}: g.SetShouldClose,
	}
	//special function keys
	ctrl_down := ebiten.IsKeyPressed(ebiten.KeyControl)
	alt_down := ebiten.IsKeyPressed(ebiten.KeyAlt)
	meta_down := ebiten.IsKeyPressed(ebiten.KeyMeta)
	shift_down := ebiten.IsKeyPressed(ebiten.KeyShift)

	for shortcut := range global_shortcuts {
		switch shortcut.key {
		case ebiten.KeyControl, ebiten.KeyShift, ebiten.KeyAlt, ebiten.KeyMeta:
			continue
		default:
			if inpututil.IsKeyJustReleased(shortcut.key) {
				executable := global_shortcuts[KeyShortcut{
					mod_shift: shift_down,
					mod_ctrl:  ctrl_down,
					mod_alt:   alt_down,
					mod_meta:  meta_down,
					key:       shortcut.key,
				}]
				if executable != nil {
					executable()
				}
			}
		}
	}

	if key_consumer != g.last_keyboard_consumer {
		if g.last_keyboard_consumer != nil {
			g.last_keyboard_consumer.KeyboardFocusLost()
		}
	}
	g.last_keyboard_consumer = key_consumer

	//Keyboard handling
	if g.last_keyboard_consumer != nil {
		g.last_keyboard_consumer.TakeKeyboard()
	}
	return nil
}

func (g *Editor) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, 0, 0, float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy()), Style.BGColorMuted)
	g.MainWidget.Draw(screen)

}

func (g *Editor) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideHeight == g.screenWidth && outsideWidth == g.screenWidth {
		//nothing changed
		return outsideWidth, outsideHeight
	}
	r := image.Rect(0, 0, outsideWidth, outsideHeight)
	g.screenWidth = outsideWidth
	g.screenHeight = outsideHeight

	g.MainWidget.SetRect(r)
	return g.screenWidth, g.screenHeight
}

/*
[]MenuItem{
	&DummyMenuItem{txt: "File", kids: []MenuItem{&DummyMenuItem{txt: "Open"}, &DummyMenuItem{txt: "Close"}, &DummyMenuItem{txt: "Quit"}}},
	&DummyMenuItem{txt: "Edit", kids: []MenuItem{&DummyMenuItem{txt: "Copy", kids: []MenuItem{&DummyMenuItem{txt: "Cut"}, &DummyMenuItem{txt: "Pasta"}}}}},
	&DummyMenuItem{txt: "Code", kids: []MenuItem{}},
},
*/

func main() {

	menu_items := []MenuItem{
		NewMenuItem("File", []MenuItem{NewMenuItem("Save", nil), NewMenuItem("Save as", nil), NewMenuItem("Open", nil), NewMenuItem("Close", nil), NewMenuItem("Quit", nil)}),
		NewMenuItem("Edit", []MenuItem{NewMenuItem("Copy", nil), NewMenuItem("Cut", nil), NewMenuItem("Pasta", nil)}),
		NewMenuItem("Code", []MenuItem{NewMenuItem("Go To", []MenuItem{NewMenuItem("Symbol Definition", nil)})}),
	}
	te1 := &TextEditor{text: strings.Split("", "\n")}

	var data_pane *TextEditor = &TextEditor{
		ReadOnly: true,
	}
	ticker := time.NewTicker(time.Second / 60)
	go func() {
		for t := range ticker.C {
			data_pane.SetText(fmt.Sprintf("\nData:\ntime: %v\nTPS: %f\nFPS: %f\nTicks: %d\nTETick: %d", t.Format(time.Kitchen), ebiten.ActualTPS(), ebiten.ActualFPS(), ticks, te1.last_interact_time))
		}
	}()

	main_view := &HorizontalSplitter{
		split_x: 200,
		Left:    data_pane,
		Right: &Tabs{
			current_hovered: -1,
			Titles:          []string{"Text editor", "Blue", "Green", "Red"},
			Tabs: []Widget{
				te1,
				NewColorRect(Style.BlueMuted),
				NewColorRect(Style.GreenMuted),
				NewColorRect(Style.RedMuted),
			},
			CurrentTab: 0,
			TabHeight:  2*tab_y_padding + MainFontSize,
		},
		border_half_width: 2,
		border_mode:       ShowOnHover,
	}
	g := &Editor{
		MainWidget: NewMenuBar(menu_items, main_view),
	}

	//
	//ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g.MainWidget.SetRect(image.Rect(0, 0, 800, 700))
	g.Layout(800, 800)
	ebiten.SetWindowSize(800, 800)

	ebiten.SetWindowTitle("IDE")
	if err := ebiten.RunGame(g); err != nil {
		log.Println(err)
	}
	f, err := os.Create("mem.pprof")
	check(err)
	defer f.Close()
	pprof.Lookup("allocs").WriteTo(f, 0)

}
