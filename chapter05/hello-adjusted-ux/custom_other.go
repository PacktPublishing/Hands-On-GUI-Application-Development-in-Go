// +build !darwin

package main

import "github.com/andlabs/ui"

func layoutQuit() ui.Control {
	button := ui.NewButton("Exit")
	button.OnClicked(func(*ui.Button) {
		ui.Quit()
	})

	return button
}
