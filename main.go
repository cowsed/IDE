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

type Editor struct {
	screenWidth  int
	screenHeight int

	MainWidget          Widget
	last_mouse_consumer Widget
}

func (g *Editor) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return errors.New("game ended by player")
	}
	x, y := ebiten.CursorPosition()
	consumer := g.MainWidget.MouseOver(x, y)
	if consumer != g.last_mouse_consumer {
		if g.last_mouse_consumer != nil {
			g.last_mouse_consumer.MouseOut()
		}
	}
	g.last_mouse_consumer = consumer

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.MainWidget.LMouseDown(x, y)
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.MainWidget.LMouseUp(x, y)
	}

	//Keyboard
	consumer.TakeKeyboard(ebiten.Key0)
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
	data_pane := &TextEditor{
		ReadOnly: true,
	}
	ticker := time.NewTicker(time.Second / 60)
	go func() {
		for t := range ticker.C {
			data_pane.SetText(fmt.Sprintf("time: %v\nTPS: %f\nFPS: %f", t.Format(time.Kitchen), ebiten.ActualTPS(), ebiten.ActualFPS()))
		}
	}()
	main_view := &HorizontalSplitter{
		split_x: 200,
		Left:    data_pane,
		Right: &Tabs{
			current_hovered: -1,
			Titles:          []string{"Text editor", "Blue", "Green", "Red"},
			Tabs: []Widget{
				&TextEditor{text: strings.Split("", "\n")},
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
