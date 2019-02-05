package main

import "github.com/andlabs/ui"

func layoutQuit() ui.Control {
	button := ui.NewButton("Quit")
	button.OnClicked(func(*ui.Button) {
		ui.Quit()
	})

	box := ui.NewHorizontalBox()
	box.Append(ui.NewLabel(""), true)
	box.Append(button, false)

	return box
}
