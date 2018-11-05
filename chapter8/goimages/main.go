package main

import (
	"flag"
	"fmt"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/iconvg"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/exp/shiny/unit"
	"golang.org/x/exp/shiny/widget"
	"golang.org/x/exp/shiny/widget/node"
	"golang.org/x/exp/shiny/widget/theme"
	"golang.org/x/image/draw"
	"image"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

const space = 10

var padSize = unit.DIPs(space * 2)
var spaceSize = unit.DIPs(space)

var name *widget.Label
var view *scaledImage
var images []image.Image
var names []string
var index = 0

func changeImage(offset int) {
	newidx := index + offset
	if newidx < 0 || newidx >= len(images) {
		return
	}

	chooseImage(newidx, images[newidx])
}

func previousImage() {
	changeImage(-1)
}

func nextImage() {
	changeImage(1)
}

func chooseImage(idx int, img image.Image) {
	index = idx
	view.SetImage(img)
	name.Text = names[idx]
	name.Mark(node.MarkNeedsMeasureLayout)
	name.Mark(node.MarkNeedsPaintBase)

	//	win.Send(paint.Event{})
}

func expandSpace() node.Node {
	return widget.WithLayoutData(widget.NewSpace(),
		widget.FlowLayoutData{ExpandAlong: true, ExpandAcross: true, AlongWeight: 1})
}

func makeBar() node.Node {
	prev := newButton("Previous", previousImage)
	next := newButton("Next", nextImage)
	name = widget.NewLabel("Filename")

	flow := widget.NewFlow(widget.AxisHorizontal, prev, expandSpace(),
		widget.NewPadder(widget.AxisBoth, padSize, name), expandSpace(), next)

	bar := widget.NewUniform(theme.Neutral, flow)

	return widget.WithLayoutData(bar,
		widget.FlowLayoutData{ExpandAlong: true, ExpandAcross: true})
}

func loadDirIcon() image.Image {
	var raster iconvg.Rasterizer
	bounds := image.Rect(0, 0, iconSize, iconSize)
	icon := image.NewRGBA(bounds)
	raster.SetDstImage(icon, bounds, draw.Over)

	iconvg.Decode(&raster, icons.FileFolder, nil)
	return icon
}

func makeCell(idx int, name string, icon image.Image) node.Node {
	var onClick func()
	if icon == nil {
		icon = loadDirIcon()
	} else {
		onClick = func() { chooseImage(idx, icon) }
	}

	return newCell(icon, name, onClick)
}

func makeList(dir string, files []string) node.Node {
	parent := makeCell(-1, filepath.Base(dir), nil)
	children := []node.Node{parent}

	for idx, name := range files {
		img := loadImage(path.Join(dir, name))
		cell := makeCell(idx, name, img)
		children = append(children, cell)
		images = append(images, img)
	}

	return widget.NewFlow(widget.AxisVertical, children...)
}

func loadImage(name string) image.Image {
	reader, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	return img
}

func scaleImage(src image.Image, width, height int) image.Image {
	ret := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.ApproxBiLinear.Scale(ret, ret.Bounds(), src, src.Bounds(), draw.Src, nil)

	return ret
}

func getImageList(dir string) []string {
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

func loadUI(dir string) {
	driver.Main(func(s screen.Screen) {
		var img image.Image
		names := getImageList(dir)

		view = newScaledImage(img)
		scaledImage := widget.WithLayoutData(view,
			widget.FlowLayoutData{ExpandAlong: true, ExpandAcross: true, AlongWeight: 4})

		body := widget.NewFlow(widget.AxisHorizontal, makeList(dir, names),
			widget.NewPadder(widget.AxisHorizontal, spaceSize, nil), scaledImage)
		expanding := widget.WithLayoutData(widget.NewPadder(widget.AxisBoth, padSize, body),
			widget.FlowLayoutData{ExpandAlong: true, ExpandAcross: true, AlongWeight: 4})
		container := widget.NewFlow(widget.AxisVertical, makeBar(), expanding)
		sheet := widget.NewSheet(widget.NewUniform(theme.Background, container))

		if len(images) > 0 {
			chooseImage(0, images[0])
		}

		container.Measure(theme.Default, 0, 0)
		if err := widget.RunWindow(s, sheet, &widget.RunWindowOptions{
			NewWindowOptions: screen.NewWindowOptions{
				Title:  "GoImages",
				Width:  container.MeasuredSize.X,
				Height: container.MeasuredSize.Y,
			},
		}); err != nil {
			log.Fatal(err)
		}
	})
}

func main() {
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
	loadUI(dir)
}
