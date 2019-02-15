package main

import (
	"fmt"
	"github.com/PacktPublishing/Hands-On-GUI-Application-Development-in-Go/client"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

const padding = 3

type mainUI struct {
	content *gtk.TextView
	subject, to, from, date *gtk.Label

	server    client.EmailServer
	listModel *gtk.ListStore
	listIter  gtk.TreeIter
}

func (m *mainUI) buildMenu() *gtk.MenuBar {
	menubar := gtk.NewMenuBar()

	fileMenu := gtk.NewMenuItemWithLabel("File")
	menu := gtk.NewMenu()
	item := gtk.NewMenuItemWithLabel("New")
	item.Connect("activate", m.showCompose)
	menu.Append(item)
	menu.Append(gtk.NewSeparatorMenuItem())
	item = gtk.NewMenuItemWithLabel("Quit")
	item.Connect("activate", func() {
		gtk.MainQuit()
	})
	menu.Append(item)
	fileMenu.SetSubmenu(menu)

	menubar.Append(fileMenu)
	menubar.Append(gtk.NewMenuItemWithLabel("Edit"))
	menubar.Append(gtk.NewMenuItemWithLabel("Help"))

	return menubar
}

func (m *mainUI) buildToolbar() *gtk.Toolbar {
	toolbar := gtk.NewToolbar()
	toolbar.SetStyle(gtk.TOOLBAR_BOTH)

	item := gtk.NewToolButtonFromStock(gtk.STOCK_NEW)
	item.OnClicked(m.showCompose)
	toolbar.Add(item)
	toolbar.Add(gtk.NewToolButton(nil, "Reply"))
	toolbar.Add(gtk.NewToolButton(nil, "Reply All"))
	toolbar.Add(gtk.NewSeparatorToolItem())
	toolbar.Add(gtk.NewToolButtonFromStock(gtk.STOCK_DELETE))
	toolbar.Add(gtk.NewSeparatorToolItem())
	toolbar.Add(gtk.NewToolButtonFromStock(gtk.STOCK_CUT))
	toolbar.Add(gtk.NewToolButtonFromStock(gtk.STOCK_COPY))
	toolbar.Add(gtk.NewToolButtonFromStock(gtk.STOCK_PASTE))

	return toolbar
}

func newLeftLabel(text string) *gtk.Label {
	label := gtk.NewLabel(text)
	label.SetAlignment(0, 0)

	return label
}

func newBoldLeftLabel(text string) *gtk.Label {
	label := newLeftLabel(text)
	label.SetMarkup(fmt.Sprintf("<b>%s</b>", text))

	return label
}

func (m *mainUI) prependEmail(message *client.EmailMessage) {
	m.listModel.Prepend(&m.listIter)
	m.listModel.SetValue(&m.listIter, 0, message.Subject)
}

func (m *mainUI) setEmail(message *client.EmailMessage) {
	m.subject.SetText(message.Subject)
	m.to.SetText(message.ToEmailString())
	m.from.SetText(message.FromEmailString())
	m.date.SetText(message.DateString())

	m.content.GetBuffer().SetText(message.Content)
}

func (m *mainUI) Select(selection *gtk.TreeSelection, model *gtk.TreeModel, path *gtk.TreePath, selected bool) bool {
	if selected { // already selected, just return
		return true
	}

	row := path.GetIndices()[0]
	email := m.server.ListMessages()[row]

	m.setEmail(email)
	return true
}

func (m *mainUI) showMain(server client.EmailServer) {
	gtk.Init(nil)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("GoMail")
	window.Connect("destroy", func() {
		gtk.MainQuit()
	})

	toolbar := gtk.NewToolbar()
	toolbar.SetStyle(gtk.TOOLBAR_BOTH)

	m.listModel = gtk.NewListStore(gtk.TYPE_STRING)
	list := gtk.NewTreeView()
	list.SetModel(m.listModel)
	list.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Inbox", gtk.NewCellRendererText(), "text", 0))

	var workaround gtk.GtkTreeSelecter
	workaround = m
	list.GetSelection().SetSelectFunction(&workaround)

	messages := server.ListMessages()
	for i := len(messages)-1; i >= 0; i-- {
		m.prependEmail(messages[i])
	}

	m.subject = newBoldLeftLabel("subject")
	m.subject.SetAlignment(0, 0)

	meta := gtk.NewHBox(false, padding)
	labels := gtk.NewVBox(true, padding)
	labels.Add(newBoldLeftLabel("To"))
	labels.Add(newBoldLeftLabel("From"))
	labels.Add(newBoldLeftLabel("Date"))
	values := gtk.NewVBox(true, padding)
	m.to = newLeftLabel("email")
	values.Add(m.to)
	m.from = newLeftLabel("email")
	values.Add(m.from)
	m.date = newLeftLabel("date")
	values.Add(m.date)
	meta.PackStart(labels, false, true, 0)
	meta.Add(values)

	m.content = gtk.NewTextView()
	m.content.SetEditable(false)

	detail := gtk.NewVBox(false, padding)
	detail.PackStart(m.subject, false, true, 0)
	detail.PackStart(meta, false, true, 0)
	detail.Add(m.content)

	split := gtk.NewHPaned()
	split.Add1(list)
	split.Add2(detail)

	vbox := gtk.NewVBox(false, padding)
	vbox.PackStart(m.buildMenu(), false, true, 0)
	vbox.PackStart(m.buildToolbar(), false, true, 0)
	vbox.Add(split)

	window.Add(vbox)
	window.SetBorderWidth(padding)
	window.Resize(600, 400)
	window.ShowAll()
}

func (m *mainUI) showCompose() {
	showCompose(m.server)
}

func main() {
	glib.ThreadInit(nil)
	gdk.ThreadsInit()
	gdk.ThreadsEnter()
	gtk.Init(nil)

	server := client.NewTestServer()
	main := new(mainUI)
	main.server = server
	main.showMain(server)

	main.setEmail(server.CurrentMessage())
	go func() {
		for email := range server.Incoming() {
			gdk.ThreadsEnter()
			main.prependEmail(email)
			gdk.ThreadsLeave()
		}
	}()

	gtk.Main()
}