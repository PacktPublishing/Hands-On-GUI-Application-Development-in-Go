package main

import (
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"log"
)

func clicked(ctx *glib.CallbackContext) {
	button := ctx.Data().(*gtk.Button)
	log.Println("Button clicked was:", button.GetLabel())
}

func addButton(label string) *gtk.Button {
	button := gtk.NewButton()
	button.SetLabel(label)
	button.Clicked(clicked, button)

	return button
}

func main() {
	gtk.Init(nil)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("Hello")
	window.Connect("destroy", func() {
		gtk.MainQuit()
	})

	quit := gtk.NewButton()
	quit.SetLabel("Quit")
	quit.Clicked(func() {
		gtk.MainQuit()
	})

	vbox := gtk.NewVBox(false, 3)
	vbox.Add(addButton("One"))
	vbox.Add(addButton("Two"))
	vbox.Add(addButton("Three"))
	vbox.Add(gtk.NewHSeparator())
	vbox.Add(quit)

	window.Add(vbox)
	window.SetBorderWidth(3)
	window.ShowAll()
	gtk.Main()
}
