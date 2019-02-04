package main

import "log"

import "github.com/gobuffalo/packr"

func main() {
	box := packr.NewBox("./data")
	data, err := box.Find("shiny-hall.jpg")

	log.Println("datLen", len(data), "err", err)
}

