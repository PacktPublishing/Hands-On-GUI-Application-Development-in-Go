package main

import "github.com/andlabs/ui"

func main() {
	err := ui.Main(func() {
		window := ui.NewWindow("Hello", 100, 50, false)
		window.SetMargined(true)
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})

		button := layoutQuit()
		box := ui.NewVerticalBox()
		box.SetPadded(true)
		box.Append(ui.NewLabel("Hello World!"), false)
		box.Append(button, false)

		window.SetChild(box)
		window.Show()
	})
	if err != nil {
		panic(err)
	}
}
