package main

import (
	"github.com/PacktPublishing/Hands-On-GUI-Application-Development-in-Go/client"
	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/app"
	"github.com/fyne-io/fyne/layout"
	"github.com/fyne-io/fyne/theme"
	"github.com/fyne-io/fyne/widget"
)

type mainUI struct {
	app    fyne.App
	server *client.EmailServer

	list *widget.Group

	content, subject, to, from, date *widget.Label
}

func (m *mainUI) buildToolbar() *widget.Toolbar {
	return widget.NewToolbar(
		widget.NewToolbarAction(theme.MailComposeIcon(), func() {
			newCompose(m.app, m.server).loadUI().Show()
		}),
		widget.NewToolbarAction(theme.MailReplyIcon(), func() {
		}),
		widget.NewToolbarAction(theme.MailReplyAllIcon(), func() {
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.DeleteIcon(), func() {
		}),
		widget.NewToolbarAction(theme.CutIcon(), func() {
		}),
		widget.NewToolbarAction(theme.CopyIcon(), func() {
		}),
		widget.NewToolbarAction(theme.PasteIcon(), func() {
		}),
	)
}

func (m *mainUI) loadUI() fyne.Window {
	browse := m.app.NewWindow("GoMail")
	toolbar := m.buildToolbar()
	m.list = widget.NewGroup("Inbox")
	for _, email := range m.server.ListMessages() {
		m.list.Append(m.addEmail(email))
	}
	m.content = widget.NewLabel("content")
	m.subject = widget.NewLabel("subject")
	m.subject.TextStyle = fyne.TextStyle{Bold: true}
	meta := widget.NewForm()
	m.to = widget.NewLabel("email")
	meta.Append("To", m.to)
	m.from = widget.NewLabel("email")
	meta.Append("From", m.from)
	m.date = widget.NewLabel("date")
	meta.Append("Date", m.date)
	box := widget.NewVBox(meta, m.content)
	detail := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(m.subject, nil, box, nil),
		m.subject, box)
	container := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(toolbar, nil, m.list, nil),
		toolbar, m.list, detail,
	)

	browse.SetContent(container)
	browse.Resize(fyne.NewSize(600, 400))
	return browse
}

func (m *mainUI) setMessage(email *client.EmailMessage) {
	m.subject.SetText(email.Subject)

	m.to.SetText(email.ToEmailString())
	m.from.SetText(email.FromEmailString())
	m.date.SetText(email.DateString())

	m.content.SetText(email.Content)
}

func (m *mainUI) addEmail(email *client.EmailMessage) fyne.CanvasObject {
	return widget.NewButton(email.Subject, func() {
		m.setMessage(email)
	})
}

func newMainUI(mailApp fyne.App, server *client.EmailServer) *mainUI {
	ui := &mainUI{
		app:    mailApp,
		server: server,
	}

	return ui
}

func main() {
	mailApp := app.New()
	server := client.NewTestServer()

	ui := newMainUI(mailApp, server)
	win := ui.loadUI()

	go func() {
		for email := range server.Incoming() {
			ui.list.Prepend(ui.addEmail(email))
		}
	}()

	ui.setMessage(server.CurrentMessage())
	win.ShowAndRun()
}
