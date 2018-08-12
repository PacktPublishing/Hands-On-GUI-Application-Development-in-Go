package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	MainWindow{
		Title:  "Hello",
		Layout: VBox{},
		Children: []Widget{
			Label{Text: "Hello World!"},
			PushButton{
				Text: "Quit",
				OnClicked: func() {
					walk.App().Exit(0)
				},
			},
		},
	}.Run()
}
