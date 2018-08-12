package main

import (
	"log"

	"github.com/lxn/walk"
)

var marginSize = 9

func buildWindow() (*walk.MainWindow, error) {
	win, err := walk.NewMainWindowWithName("Hello")
	if err != nil {
		return nil, err
	}
	layout := walk.NewVBoxLayout()
	layout.SetMargins(walk.Margins{marginSize, marginSize, marginSize, marginSize})
	layout.SetSpacing(marginSize)
	win.SetLayout(layout)

	label, err := walk.NewLabel(win)
	if err != nil {
		return win, err
	}
	label.SetText("Hello World!")

	button, err := walk.NewPushButton(win)
	if err != nil {
		return win, err
	}
	button.SetText("Quit")
	button.Clicked().Attach(func() {
		walk.App().Exit(0)
	})

	return win, nil
}

func main() {
	win, err := buildWindow()
	if err != nil {
		log.Fatalln(err)
	}

	win.SetVisible(true)
	win.Run()
}
