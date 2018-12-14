package main

import (
	"github.com/PacktPublishing/Hands-On-GUI-Application-Development-in-Go/client"
	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/layout"
	"github.com/fyne-io/fyne/widget"
	"time"
)

type composeUI struct {
	app    fyne.App
	server *client.EmailServer

	list *widget.Group

	content, subject, to, from, date *widget.Label
}

func (c *composeUI) loadUI() fyne.Window {
	compose := c.app.NewWindow("GoMail Compose")

	subject := widget.NewEntry()
	subject.SetText("subject")
	toLabel := widget.NewLabel("To")
	to := widget.NewEntry()
	to.SetText("email")

	message := widget.NewEntry()
	message.SetText("content")

	send := widget.NewButton("Send", func() {
		email := client.NewMessage(c.subject.Text, c.content.Text,
			client.Email(c.to.Text), "", time.Now())
		c.server.Send(email)
		compose.Close()
	})
	send.Style = widget.PrimaryButton
	buttons := widget.NewHBox(
		layout.NewSpacer(),
		widget.NewButton("Cancel", func() {
			compose.Close()
		}),
		send)

	top := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(subject, nil, toLabel, nil),
		subject, toLabel, to)
	content := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(top, buttons, nil, nil),
		top, message, buttons)
	compose.SetContent(content)

	compose.Resize(fyne.NewSize(400, 320))
	return compose
}

func newCompose(mailApp fyne.App, server *client.EmailServer) *composeUI {
	ui := &composeUI{
		app:    mailApp,
		server: server,
	}

	return ui
}
