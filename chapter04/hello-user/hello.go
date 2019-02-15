package main

import (
	"fmt"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	var message *walk.Label
	var userName *walk.TextEdit

	MainWindow{
		Title:  "Hello",
		Layout: VBox{},
		Children: []Widget{
			Label{
				AssignTo: &message,
				Text:     "Hello World!",
			},
			TextEdit{
				AssignTo: &userName,
				OnTextChanged: func() {
					welcome := fmt.Sprintf("Hello %s!", userName.Text())
					message.SetText(welcome)
				},
			},
			PushButton{
				Text: "Quit",
				OnClicked: func() {
					walk.App().Exit(0)
				},
			},
		},
	}.Run()
}
