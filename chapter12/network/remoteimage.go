package main

import (
	_ "image/png"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"

	"github.com/PacktPublishing/Hands-On-GUI-Application-Development-in-Go/chapter12/network/remote"
)

func main() {
	app := app.New()
	w := app.NewWindow("Remote Image")

	stream := remote.ReadStream("https://golang.org/doc/gopher/frontpage.png")
	img := canvas.NewImageFromImage(remote.RemoteImage(stream))
	img.SetMinSize(fyne.Size{180, 250})
	w.SetContent(img)
	w.ShowAndRun()
}
