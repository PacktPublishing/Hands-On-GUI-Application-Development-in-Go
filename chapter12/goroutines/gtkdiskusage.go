package main

import (
	"fmt"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"os"
	"path/filepath"

	"github.com/PacktPublishing/Hands-On-GUI-Application-Development-in-Go/chapter12/goroutines/disk"
)

func gtkReportUsage(path string, list *gtk.ListStore, totalLabel *gtk.Label) {
	f, _ := os.Open(path)
	files, _ := f.Readdir(-1)
	f.Close()

	result := make(chan disk.SizeInfo)
	for _, file := range files {
		go disk.DirSize(filepath.Join(path, file.Name()), result)
	}

	var total int64
	results := 0
	for info := range result {
		var listIter gtk.TreeIter
		total += info.Size

		gdk.ThreadsEnter()
		list.Append(&listIter)
		list.SetValue(&listIter, 0, info.Name)
		list.SetValue(&listIter, 1, disk.FormatSize(info.Size))
		gdk.ThreadsLeave()

		results++
		if results == len(files) {
			break
		}
	}

	gdk.ThreadsEnter()
	totalLabel.SetText(fmt.Sprintf("Total: %s", disk.FormatSize(total)))
	gdk.ThreadsLeave()
}

func main() {
	glib.ThreadInit(nil)
	gdk.ThreadsInit()
	gdk.ThreadsEnter()
	gtk.Init(nil)

	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("Disk Usage")
	window.Connect("destroy", func() {
		gtk.MainQuit()
	})

	path, _ := os.Getwd()

	if len(os.Args) == 2 {
		path = os.Args[1]
	}

	list := gtk.NewTreeView()
	listModel := gtk.NewListStore(gtk.TYPE_STRING, gtk.TYPE_STRING)
	list.SetModel(listModel)
	list.AppendColumn(gtk.NewTreeViewColumnWithAttributes(path, gtk.NewCellRendererText(), "text", 0))
	list.AppendColumn(gtk.NewTreeViewColumnWithAttributes("size", gtk.NewCellRendererText(), "text", 1))

	total := gtk.NewLabel("Total: calculating")
	layout := gtk.NewVBox(false, 3)
	layout.Add(list)
	layout.Add(total)

	window.Add(layout)
	window.SetBorderWidth(3)
	window.ShowAll()

	go gtkReportUsage(path, listModel, total)
	gtk.Main()
}
