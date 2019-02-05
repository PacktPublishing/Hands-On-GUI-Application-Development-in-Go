package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type toolbarLabel struct {
}

func (t *toolbarLabel) ToolbarObject() fyne.CanvasObject {
	label = widget.NewLabel("filename")
	return label
}

type asyncImage struct {
	path   string
	image  *canvas.Image
	pixels image.Image
}

func (a *asyncImage) ColorModel() color.Model {
	if a.pixels == nil {
		return color.RGBAModel
	}

	return a.pixels.ColorModel()
}

func (a *asyncImage) Bounds() image.Rectangle {
	if a.pixels == nil {
		return image.ZR
	}

	return a.pixels.Bounds()
}

func (a *asyncImage) At(x, y int) color.Color {
	if a.pixels == nil {
		return color.Transparent
	}

	return a.pixels.At(x, y)
}

func (a *asyncImage) load() {
	if a.path == "" {
		return
	}
	reader, err := os.Open(a.path)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	a.pixels, _, err = image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	canvas.Refresh(a.image)
}

func (a *asyncImage) loadPath(path string) {
	a.path = path
	go a.load()
}

func newAsyncImage(path string) *asyncImage {
	async := &asyncImage{}
	async.image = canvas.NewImageFromImage(async)
	async.loadPath(path)

	return async
}

var images []string
var index int

var async *asyncImage
var label *widget.Label

func previousImage() {
	if index == 0 {
		return
	}

	chooseImage(index - 1)
}

func nextImage() {
	if index == len(images)-1 {
		return
	}

	chooseImage(index + 1)
}

func chooseImage(id int) {
	path := images[id]
	label.SetText(filepath.Base(path))
	async.loadPath(path)
	index = id
}

func makeRow(id int, path string) fyne.CanvasObject {
	filename := filepath.Base(path)
	button := widget.NewButton(filename, func() {
		chooseImage(id)
	})

	preview := newAsyncImage(path).image
	iconHeight := button.MinSize().Height
	preview.SetMinSize(fyne.NewSize(int(float32(iconHeight)*1.5), iconHeight))

	return fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, preview, nil),
		preview, button)
}

func checkerColor(x, y, _, _ int) color.Color {
	xr := x / 10
	yr := y / 10

	if xr%2 == yr%2 {
		return color.RGBA{0xc0, 0xc0, 0xc0, 0xff}
	} else {
		return color.RGBA{0x99, 0x99, 0x99, 0xff}
	}
}

func getImageList(dir string) []string {
	var names []string
	files, _ := ioutil.ReadDir(dir)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(file.Name()))
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
			names = append(names, file.Name())
		}
	}

	return names
}

func makeList(dir string) *widget.Group {
	files := getImageList(dir)
	group := widget.NewGroup(filepath.Base(dir))

	for idx, name := range files {
		path := filepath.Join(dir, name)
		images = append(images, path)

		group.Append(makeRow(idx, path))
	}

	return group
}

func parseArgs() string {
	dir, _ := os.Getwd()

	flag.Usage = func() {
		fmt.Println("goimages takes a single, optional, directory parameter")
	}
	flag.Parse()

	if len(flag.Args()) > 1 {
		flag.Usage()
		os.Exit(2)
	} else if len(flag.Args()) == 1 {
		dir = flag.Args()[0]

		if _, err := ioutil.ReadDir(dir); os.IsNotExist(err) {
			fmt.Println("Directory", dir, "does not exist or could not be read")
			os.Exit(1)
		}
	}

	return dir
}

func main() {
	imageApp := app.New()
	win := imageApp.NewWindow("GoImages")

	navBar := widget.NewToolbar(
		widget.NewToolbarAction(theme.NavigateBackIcon(), previousImage),
		widget.NewToolbarSpacer(),
		&toolbarLabel{},
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.NavigateNextIcon(), nextImage))
	fileList := makeList(parseArgs())

	checkers := canvas.NewRaster(checkerColor)
	async = newAsyncImage("")
	async.image.FillMode = canvas.ImageFillContain
	chooseImage(0)

	container := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(navBar, nil, fileList, nil),
		navBar, fileList, checkers, async.image,
	)

	win.SetContent(container)
	win.Resize(fyne.NewSize(640, 480))

	win.ShowAndRun()
}
