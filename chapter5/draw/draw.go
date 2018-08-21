package main

import (
	"github.com/andlabs/ui"
	"log"
)

type areaHandler struct {
}

func (areaHandler) Draw(a *ui.Area, dp *ui.AreaDrawParams) {
 	p := ui.NewPath(ui.Winding)
 	p.NewFigure(10, 10)
	p.LineTo(dp.ClipWidth - 10, 10)
 	p.LineTo(dp.ClipWidth - 10, dp.ClipHeight - 10)
 	p.LineTo(10, dp.ClipHeight - 10)
	p.CloseFigure()
	p.End()

	dp.Context.Fill(p, &ui.Brush{Type:ui.Solid, R:.75, G:.25, B:0, A:1})
	dp.Context.Stroke(p, &ui.Brush{Type:ui.Solid, R:.25, G:.25, B:.75, A:.5},
		&ui.StrokeParams{Thickness: 4, Dashes: []float64{10, 6}, Cap:ui.RoundCap})
 	p.Free()
}

func (areaHandler) MouseEvent(a *ui.Area, me *ui.AreaMouseEvent) {
	log.Println("Mouse at", me.X, me.Y)
}

func (areaHandler) MouseCrossed(a *ui.Area, left bool) {
	// ignore
}

func (areaHandler) DragBroken(a *ui.Area) {
	// ignore
}

func (areaHandler) KeyEvent(a *ui.Area, ke *ui.AreaKeyEvent) (handled bool) {
	log.Println("Key code", ke.Key)

	return false
}

func main() {
	err := ui.Main(func() {
		window := ui.NewWindow("Draw", 200, 150, false)
		window.SetMargined(false)
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})

		handler := new(areaHandler)
		box := ui.NewVerticalBox()
		box.Append(ui.NewArea(handler), true)

		window.SetChild(box)
		window.Show()
	})
	if err != nil {
		panic(err)
	}
}
