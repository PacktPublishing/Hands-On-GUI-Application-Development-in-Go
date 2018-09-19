package main

import (
	"github.com/therecipe/qt/widgets"
)

func showCompose() {
	dialog := widgets.NewQDialog(nil, 0)
	dialog.SetModal(false)
	dialog.SetWindowTitle("New GoMail")

	form := widgets.NewQFormLayout(dialog)
	dialog.SetLayout(form)
	dialog.SetMinimumSize2(400, 320)

	form.AddRow5(widgets.NewQLineEdit2("subject", dialog))
	form.AddRow3("To", widgets.NewQLineEdit2("email", dialog))
	form.AddRow5(widgets.NewQTextEdit2("content", dialog))

	buttons := widgets.NewQWidget(dialog, 0)
	buttons.SetLayout(widgets.NewQHBoxLayout())
	buttons.Layout().AddItem(widgets.NewQSpacerItem(0, 0, widgets.QSizePolicy__Expanding, 0))
	buttons.Layout().AddWidget(widgets.NewQPushButton2("Cancel", buttons))
	send := widgets.NewQPushButton2("Send", buttons)
	send.SetDefault(true)
	buttons.Layout().AddWidget(send)
	form.AddRow5(buttons)

	dialog.Show()
}
