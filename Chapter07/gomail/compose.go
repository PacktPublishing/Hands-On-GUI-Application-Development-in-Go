package main

import (
	"github.com/PacktPublishing/Hands-On-GUI-Application-Development-in-Go/client"
	"github.com/therecipe/qt/widgets"
	"time"
)

type composeUI struct {
	server *client.EmailServer

	dialog      *widgets.QDialog
	to, subject *widgets.QLineEdit
	content     *widgets.QTextEdit
}

func (c *composeUI) createEmail() *client.EmailMessage {
	email := &client.EmailMessage{}

	email.Subject = c.subject.Text()
	email.To = client.Email(c.to.Text())
	email.Date = time.Now()
	email.Content = c.content.ToPlainText()

	return email
}

func (c *composeUI) show() {
	dialog := widgets.NewQDialog(nil, 0)
	dialog.SetModal(false)
	dialog.SetWindowTitle("New GoMail")
	c.dialog = dialog

	form := widgets.NewQFormLayout(dialog)
	dialog.SetLayout(form)
	dialog.SetMinimumSize2(400, 320)

	form.AddRow5(widgets.NewQLineEdit2("subject", dialog))
	form.AddRow3("To", widgets.NewQLineEdit2("email", dialog))
	form.AddRow5(widgets.NewQTextEdit2("content", dialog))

	buttons := widgets.NewQWidget(dialog, 0)
	buttons.SetLayout(widgets.NewQHBoxLayout())
	buttons.Layout().AddItem(widgets.NewQSpacerItem(0, 0, widgets.QSizePolicy__Expanding, 0))
	cancel := widgets.NewQPushButton2("Cancel", buttons)
	cancel.ConnectClicked(func(_ bool) {
		c.dialog.Close()
	})
	buttons.Layout().AddWidget(cancel)
	send := widgets.NewQPushButton2("Send", buttons)
	send.SetDefault(true)
	send.ConnectClicked(func(_ bool) {
		email := c.createEmail()
		c.server.Send(email)
		c.dialog.Close()
	})
	buttons.Layout().AddWidget(send)
	form.AddRow5(buttons)

	dialog.Show()
}

func newCompose(server *client.EmailServer) *composeUI {
	return &composeUI{server: server}
}
