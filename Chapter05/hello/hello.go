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

		button := ui.NewButton("Quit")
		button.OnClicked(func(*ui.Button) {
			ui.Quit()
		})
		box := ui.NewVerticalBox()
		box.Append(ui.NewLabel("Hello World!"), false)
		box.Append(button, false)

		window.SetChild(box)
		window.Show()
	})
	if err != nil {
		panic(err)
	}
}
