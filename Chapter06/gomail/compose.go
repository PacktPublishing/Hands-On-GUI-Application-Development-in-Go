package main

import (
	"github.com/PacktPublishing/Hands-On-GUI-Application-Development-in-Go/client"
	"github.com/mattn/go-gtk/gtk"
	"time"
)

type composeUI struct {
	server client.EmailServer

	subject, to *gtk.Entry
	content *gtk.TextView
	window *gtk.Window
}

func (c *composeUI) createEmail() *client.EmailMessage {
	email := &client.EmailMessage{}

	email.Subject = c.subject.GetText()
	email.To = client.Email(c.to.GetText())
	email.Date = time.Now()

	var start, end gtk.TextIter
	c.content.GetBuffer().GetBounds(&start, &end)
	email.Content = c.content.GetBuffer().GetText(&start, &end, false)

	return email
}

func (c *composeUI) cancelClicked() {
	c.window.Destroy()
}

func (c *composeUI) sendClicked() {
	c.server.Send(c.createEmail())
	c.window.Destroy()
}

func showCompose(server client.EmailServer) {
	compose := new(composeUI)
	compose.server = server

	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("New GoMail")
	window.Connect("destroy", func() {
		window.Destroy()
	})
	compose.window = window

	vbox := gtk.NewVBox(false, padding)
	compose.subject = gtk.NewEntry()
	compose.subject.SetText("subject")
	vbox.PackStart(compose.subject, false, true, 0)
	toBox := gtk.NewHBox(false, padding)
	toBox.PackStart(gtk.NewLabel("To"), false, true, 0)
	compose.to = gtk.NewEntry()
	compose.to.SetText("email")
	toBox.Add(compose.to)
	vbox.PackStart(toBox, false, true, 0)

	compose.content = gtk.NewTextView()
	compose.content.GetBuffer().SetText("email content")
	compose.content.SetEditable(true)
	vbox.Add(compose.content)

	buttonBox := gtk.NewHBox(false, padding)
	cancel := gtk.NewButtonWithLabel("Cancel")
	buttonBox.PackEnd(cancel, false, true, 0)
	cancel.Clicked(compose.cancelClicked)
	send := gtk.NewButtonWithLabel("Send")
	buttonBox.PackEnd(send, false, true, 0)
	send.Clicked(compose.sendClicked)
	vbox.PackEnd(buttonBox, false, true, 0)

	window.Add(vbox)
	window.SetBorderWidth(padding)
	window.Resize(400, 320)
	window.ShowAll()
}
