package main

import (
	"errors"
	"fmt"
	"image"
	_ "image/png"
	"log"

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

	return nil
}

func (g *Editor) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, 0, 0, float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy()), Style.BGColorMuted)
	g.MainWidget.Draw(screen)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %f\nFPS: %f", ebiten.ActualTPS(), ebiten.ActualFPS()), 12, 222)

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

func main() {
	g := &Editor{
		MainWidget: &MenuBar{TopLevelItems: []MenuItem{
			&DummyMenuItem{txt: "File", kids: []MenuItem{&DummyMenuItem{txt: "Open"}, &DummyMenuItem{txt: "Close"}, &DummyMenuItem{txt: "Quit"}}},
			&DummyMenuItem{txt: "Edit", kids: []MenuItem{&DummyMenuItem{txt: "Copy", kids: []MenuItem{&DummyMenuItem{txt: "Cut"}, &DummyMenuItem{txt: "Pasta"}}}}},
			&DummyMenuItem{txt: "Code", kids: []MenuItem{}}},
			WidgetIApplyTo: &HorizontalSplitter{
				split_x: 200,
				Left:    &ColorRect{color: Style.RedMuted},
				Right: &Tabs{
					current_hovered: -1,
					Titles:          []string{"Text editor", "Blue", "Green", "Red"},
					Tabs: []Widget{
						&TextEditor{},
						&ColorRect{color: Style.BGColorMuted},
						&ColorRect{color: Style.GreenStrong},
						&ColorRect{color: Style.RedStrong},
					},
					CurrentTab: 1,
					TabHeight:  2*tab_y_padding + MainFontSize,
				},
				border_half_width: 2,
				border_mode:       ShowOnHover},
		},
	}

	//
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g.MainWidget.SetRect(image.Rect(0, 0, 800, 700))
	g.Layout(800, 800)
	ebiten.SetWindowSize(800, 800)

	ebiten.SetWindowTitle("IDE")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
