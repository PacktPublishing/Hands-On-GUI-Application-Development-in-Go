package main

import (
	"github.com/PacktPublishing/Hands-On-GUI-Application-Development-in-Go/client"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type GoMailUIBrowse struct {
	subject, from, to, date walk.Label
	emailList               *walk.TreeView
	window                  *walk.MainWindow

	emailDetail *walk.DataBinder
}

func (g *GoMailUIBrowse) buildEmailActions() []MenuItem {
	return []MenuItem{
		Action{
			Text: "New",
			OnTriggered: func() {
				new(GoMailUICompose).Show(g.window)
			},
		},
		Action{
			Text: "Reply",
		},
		Action{
			Text: "Reply All",
		},
		Separator{},
		Action{
			Text: "Delete",
		},
	}
}

func (g *GoMailUIBrowse) buildEditActions() []MenuItem {
	return []MenuItem{
		Action{
			Text: "Cut",
		},
		Action{
			Text: "Copy",
		},
		Action{
			Text: "Paste",
		},
	}
}

func (g *GoMailUIBrowse) buildMenu() []MenuItem {
	return []MenuItem{
		Menu{
			Text: "File",
			Items: append(
				g.buildEmailActions(),
				Separator{},
				Action{
					Text: "Quit",
					OnTriggered: func() {
						walk.App().Exit(0)
					},
				},
			),
		},
		Menu{
			Text:  "Edit",
			Items: g.buildEditActions(),
		},
		Menu{
			Text: "Help",
		},
	}
}

func (g *GoMailUIBrowse) buildToolbar() ToolBar {
	items := append(
		g.buildEmailActions(),
		Separator{},
	)

	for _, item := range g.buildEditActions() {
		items = append(items, item)
	}

	return ToolBar{
		Items:       items,
		ButtonStyle: ToolBarButtonTextOnly,
	}
}

func (g *GoMailUIBrowse) buildUI(model *EmailClientModel) MainWindow {
	g.emailDetail = &walk.DataBinder{}
	return MainWindow{
		Title:     "GoMail",
		AssignTo:  &g.window,
		Layout:    HBox{},
		MinSize:   Size{600, 400},
		MenuItems: g.buildMenu(),
		ToolBar:   g.buildToolbar(),
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					TreeView{
						AssignTo: &g.emailList,
						Model:    model,
						OnCurrentItemChanged: func() {
							item := g.emailList.CurrentItem()

							if email, ok := item.(*EmailModel); ok {
								g.SetMessage(email.email)
							}
						},
					},
					Composite{
						Layout: Grid{Columns: 3},
						Border: false,
						DataBinder: DataBinder{
							AssignTo:   &g.emailDetail,
							DataSource: model.Server.CurrentMessage(),
						},
						Children: []Widget{
							Label{
								Text:       Bind("Subject"),
								Font:       Font{Bold: true},
								ColumnSpan: 3,
							},
							Label{
								Text: "From",
								Font: Font{Bold: true},
							},
							Label{
								Text:       Bind("FromEmailString"),
								ColumnSpan: 2,
							},
							Label{
								Text: "To",
								Font: Font{Bold: true},
							},
							Label{
								Text:       Bind("ToEmailString"),
								ColumnSpan: 2,
							},
							Label{
								Text: "Date",
								Font: Font{Bold: true},
							},
							Label{
								Text:       Bind("DateString"),
								ColumnSpan: 2,
							},
							TextEdit{
								Text:       Bind("Content"),
								ReadOnly:   true,
								ColumnSpan: 3,
							},
						},
					},
				},
			},
		},
	}
}

func (g *GoMailUIBrowse) SetMessage(email *client.EmailMessage) {
	g.emailDetail.SetDataSource(email)
	g.emailDetail.Reset()
}

func (g *GoMailUIBrowse) Run() {
	model := NewEmailServerModel()
	model.SetServer(client.NewTestServer())
	g.buildUI(model).Run()
}

func NewMailUIBrowse() {
	new(GoMailUIBrowse).Run()
}
