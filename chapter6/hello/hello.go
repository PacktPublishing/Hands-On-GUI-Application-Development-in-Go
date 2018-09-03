package main

import "github.com/mattn/go-gtk/gtk"

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

	vbox := gtk.NewVBox(false, 1)
	vbox.Add(gtk.NewLabel("Hello World!"))
	vbox.Add(quit)

	window.Add(vbox)
	window.ShowAll()
	gtk.Main()
}