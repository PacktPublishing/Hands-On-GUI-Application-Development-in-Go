package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type GoMailUICompose struct {
	cancel, send *walk.PushButton
	dialog       *walk.Dialog
}

func (g *GoMailUICompose) buildUI() Dialog {
	return Dialog{
		Title:    "New GoMail",
		AssignTo: &g.dialog,
		Layout:   HBox{},
		MinSize:  Size{400, 320},
		Children: []Widget{
			GroupBox{
				Layout: Grid{Columns: 3},
				Children: []Widget{
					LineEdit{
						Text:       "subject",
						Font:       Font{Bold: true},
						ColumnSpan: 3,
					},
					Label{
						Text: "To",
						Font: Font{Bold: true},
					},
					LineEdit{
						Text:       "email",
						ColumnSpan: 2,
					},
					TextEdit{
						Text:       "email content",
						ColumnSpan: 3,
					},
					GroupBox{
						Layout:     HBox{},
						ColumnSpan: 3,
						Children: []Widget{
							HSpacer{},
							PushButton{
								Text:     "Cancel",
								AssignTo: &g.cancel,
								OnClicked: func() {
									g.dialog.Cancel()
								},
							},
							PushButton{
								Text:     "Send",
								AssignTo: &g.send,
							},
						},
					},
				},
			},
		},
		DefaultButton: &g.send,
		CancelButton:  &g.cancel,
	}
}

func (g *GoMailUICompose) Show(owner walk.Form) {
	g.buildUI().Run(owner)
}
