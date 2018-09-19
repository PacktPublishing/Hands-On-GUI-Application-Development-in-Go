package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"os"
)

/*
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
*/
func buildMenu(app *widgets.QApplication) *widgets.QMenuBar {
	menu := widgets.NewQMenuBar(nil)

	file := widgets.NewQMenu2("File", menu)
	new := file.AddAction("New")
	new.ConnectTriggered(func(_ bool){showCompose()})
	file.AddSeparator()
	quit := file.AddAction("Quit")
	quit.ConnectTriggered(func (_ bool) {
		app.QuitDefault()
	})
	menu.AddMenu(file)

	menu.AddMenu(widgets.NewQMenu2("Edit", menu))
	menu.AddMenu(widgets.NewQMenu2("Help", menu))

	return menu
}

func buildToolbar() *widgets.QToolBar {
	toolbar := widgets.NewQToolBar("tools", nil)
	toolbar.SetToolButtonStyle(core.Qt__ToolButtonTextUnderIcon)

	new := toolbar.AddAction2(gui.QIcon_FromTheme2("document-new", nil), "New")
	new.ConnectTriggered(func(_ bool){showCompose()})
	toolbar.AddAction("Reply")
	toolbar.AddAction("Reply All")
	toolbar.AddSeparator()
	toolbar.AddAction2(gui.QIcon_FromTheme2("edit-delete", nil), "Delete")

	toolbar.AddSeparator()
	toolbar.AddAction2(gui.QIcon_FromTheme2("edit-cut", nil), "Cut")
	toolbar.AddAction2(gui.QIcon_FromTheme2("edit-copy", nil), "Copy")
	toolbar.AddAction2(gui.QIcon_FromTheme2("edit-paste", nil), "Paste")

	return toolbar
}


func main() {
	app := widgets.NewQApplication(len(os.Args), os.Args)

	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("GoMail")

	widget := widgets.NewQWidget(window, 0)
	widget.SetLayout(widgets.NewQVBoxLayout())
	window.SetMinimumSize2(600, 400)
	window.SetCentralWidget(widget)

	window.SetMenuBar(buildMenu(app))
	widget.Layout().AddWidget(buildToolbar())

	model := core.NewQStringListModel2([]string{"email1", "email2"}, widget)
	model.SetHeaderData(0, core.Qt__Horizontal, core.NewQVariant14("Inbox"), 0)
	list := widgets.NewQTreeView(window)
	list.SetModel(model)

	detail := widgets.NewQWidget(window, 0)
	form := widgets.NewQFormLayout(detail)
	detail.SetLayout(form)
	form.AddRow5(widgets.NewQLabel2("subject", detail, 0))
	form.AddRow3("From", widgets.NewQLabel2("email", detail, 0))
	form.AddRow3("To", widgets.NewQLabel2("email", detail, 0))
	form.AddRow3("Date", widgets.NewQLabel2("date", detail, 0))
	form.AddRow5(widgets.NewQLabel2("content", detail, 0))

	splitter := widgets.NewQSplitter(window)
	splitter.AddWidget(list)
	splitter.AddWidget(detail)
	widget.Layout().AddWidget(splitter)

	window.Show()
	widgets.QApplication_Exec()
}
