package main

import (
	_ "image/png"

	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/canvas"
	"github.com/fyne-io/fyne/desktop"

	"github.com/PacktPublishing/Hands-On-GUI-Application-Development-in-Go/chapter12/network/remote"
)

func main() {
	app := desktop.NewApp()
	w := app.NewWindow("Remote Image")

	stream := remote.CacheStream("https://golang.org/doc/gopher/frontpage.png")
	img := canvas.NewImageFromImage(remote.RemoteImage(stream))
	img.SetMinSize(fyne.Size{180, 250})
	w.SetContent(img)
	w.Show()
}
