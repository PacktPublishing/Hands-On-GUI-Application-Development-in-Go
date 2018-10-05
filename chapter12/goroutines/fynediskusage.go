package main

import (
	"fmt"
	"github.com/fyne-io/fyne/desktop"
	"github.com/fyne-io/fyne/widget"
	"os"
	"path/filepath"

	"github.com/PacktPublishing/Hands-On-GUI-Application-Development-in-Go/chapter12/goroutines/disk"
)

func fyneReportUsage(path string, names, sizes *widget.Box, totalLabel *widget.Label) {
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
		total += info.Size

		names.Append(widget.NewLabel(info.Name))
		sizes.Append(widget.NewLabel(fmt.Sprintf("%s   ", disk.FormatSize(info.Size))))

		results++
		if results == len(files) {
			break
		}
	}

	totalLabel.SetText(fmt.Sprintf("Total: %s", disk.FormatSize(total)))
}

func main() {
	path, _ := os.Getwd()

	if len(os.Args) == 2 {
		path = os.Args[1]
	}

	app := desktop.NewApp()
	win := app.NewWindow("Disk Usage")

	nameList := widget.NewVBox()
	sizeList := widget.NewVBox()
	total := widget.NewLabel("Total: calculating")
	win.SetContent(widget.NewVBox(
		widget.NewLabel(path),
		widget.NewHBox(nameList, sizeList),
		total,
	))

	go fyneReportUsage(path, nameList, sizeList, total)
	win.Show()
}
