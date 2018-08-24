package main

import (
	"time"

	"github.com/andlabs/ui"

	"github.com/PacktPublishing/Hands-On-GUI-Application-Development-in-Go/client"
)

type composeUI struct {
	subject, to, content *ui.Entry
	dialog               ui.Window

	server *client.EmailServer
}

func (c *composeUI) CreateMessage() *client.EmailMessage {
	email := &client.EmailMessage{}

	email.Subject = c.subject.Text()
	email.To = client.Email(c.to.Text())
	email.Content = c.content.Text()
	email.Date = time.Now()

	return email
}

func (c *composeUI) buildUI() *ui.Window {
	window := ui.NewWindow("New GoMail", 400, 320, false)
	window.SetMargined(true)
	window.OnClosing(func(*ui.Window) bool {
		return true
	})

	c.subject = ui.NewEntry()
	c.subject.SetText("subject")

	toBox := ui.NewHorizontalBox()
	toBox.Append(ui.NewLabel("To   "), false)
	c.to = ui.NewEntry()
	c.to.SetText("email")
	toBox.Append(c.to, true)

	c.content = ui.NewEntry()
	c.content.SetText("email content")

	buttonBox := ui.NewHorizontalBox()
	buttonBox.Append(ui.NewLabel(""), true)
	cancel := ui.NewButton("Cancel")
	cancel.OnClicked(func(*ui.Button) {
		window.Hide()
	})
	buttonBox.Append(cancel, false)
	send := ui.NewButton("Send")
	send.OnClicked(func(*ui.Button) {
		email := c.CreateMessage()
		c.server.Send(email)

		window.Hide()
	})
	buttonBox.Append(send, false)

	layout := ui.NewVerticalBox()
	layout.Append(c.subject, false)
	layout.Append(toBox, false)
	layout.Append(c.content, true)
	layout.Append(buttonBox, false)

	window.SetChild(layout)
	return window
}
