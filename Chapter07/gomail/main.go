package main

import (
	"github.com/PacktPublishing/Hands-On-GUI-Application-Development-in-Go/client"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"os"
)

type mainUI struct {
	core.QObject
	_     func(message *client.EmailMessage) `signal:"newMail"`
	model *core.QStringListModel

	app    *widgets.QApplication
	server *client.EmailServer

	subject, to, from, date, content *widgets.QLabel
}

func (m *mainUI) buildMenu() *widgets.QMenuBar {
	menu := widgets.NewQMenuBar(nil)

	file := widgets.NewQMenu2("File", menu)
	new := file.AddAction("New")
	new.ConnectTriggered(func(_ bool) { newCompose(m.server).show() })
	file.AddSeparator()
	quit := file.AddAction("Quit")
	quit.ConnectTriggered(func(_ bool) {
		m.app.QuitDefault()
	})
	menu.AddMenu(file)

	menu.AddMenu(widgets.NewQMenu2("Edit", menu))
	menu.AddMenu(widgets.NewQMenu2("Help", menu))

	return menu
}

func (m *mainUI) buildToolbar() *widgets.QToolBar {
	toolbar := widgets.NewQToolBar("tools", nil)
	toolbar.SetToolButtonStyle(core.Qt__ToolButtonTextUnderIcon)

	docNew := toolbar.AddAction2(gui.QIcon_FromTheme2("document-new", nil), "New")
	docNew.ConnectTriggered(func(_ bool) { newCompose(m.server).show() })
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

func (m *mainUI) buildUI() *widgets.QMainWindow {
	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("GoMail")

	widget := widgets.NewQWidget(window, 0)
	widget.SetLayout(widgets.NewQVBoxLayout())
	window.SetMinimumSize2(600, 400)
	window.SetCentralWidget(widget)

	window.SetMenuBar(m.buildMenu())
	widget.Layout().AddWidget(m.buildToolbar())

	subjects := []string{}
	for _, message := range m.server.ListMessages() {
		subjects = append(subjects, message.Subject)
	}
	m.model = core.NewQStringListModel2(subjects, widget)
	m.model.SetHeaderData(0, core.Qt__Horizontal, core.NewQVariant14("Inbox"), 0)
	list := widgets.NewQTreeView(window)
	list.SetModel(m.model)
	list.ConnectSelectionChanged(func(selected *core.QItemSelection, _ *core.QItemSelection) {
		if len(selected.Indexes()) == 0 {
			return
		}

		row := selected.Indexes()[0].Row()
		m.setMessage(m.server.ListMessages()[row])
	})

	detail := widgets.NewQWidget(window, 0)
	form := widgets.NewQFormLayout(detail)
	detail.SetLayout(form)

	m.subject = widgets.NewQLabel2("subject", detail, 0)
	form.AddRow5(m.subject)
	m.from = widgets.NewQLabel2("email", detail, 0)
	form.AddRow3("From", m.from)
	m.to = widgets.NewQLabel2("email", detail, 0)
	form.AddRow3("To", m.to)
	m.date = widgets.NewQLabel2("date", detail, 0)
	form.AddRow3("Date", m.date)
	m.content = widgets.NewQLabel2("content", detail, 0)
	form.AddRow5(m.content)

	splitter := widgets.NewQSplitter(window)
	splitter.AddWidget(list)
	splitter.AddWidget(detail)
	widget.Layout().AddWidget(splitter)

	return window
}

func (m *mainUI) setMessage(message *client.EmailMessage) {
	m.subject.SetText(message.Subject)
	m.to.SetText(message.ToEmailString())
	m.from.SetText(message.FromEmailString())
	m.date.SetText(message.DateString())

	m.content.SetText(message.Content)
}

func (m *mainUI) prependEmail(message *client.EmailMessage) {
	m.model.SetStringList(append([]string{message.Subject}, m.model.StringList()...))
}

func main() {
	main := NewMainUI(nil)
	main.app = widgets.NewQApplication(len(os.Args), os.Args)
	main.server = client.NewTestServer()

	window := main.buildUI()
	main.setMessage(main.server.CurrentMessage())
	window.Show()

	main.ConnectNewMail(main.prependEmail)
	go func() {
		for email := range main.server.Incoming() {
			main.NewMail(email)
		}
	}()

	widgets.QApplication_Exec()
}
