package main

import (
	"log"
	"os"
)

type app struct {
}

func (a *app) openFile(file string) {
	running := getFirstInstance(a)
	if running != nil {
		log.Println("Found running app, opening document", file)
		running.openFile(file)
		os.Exit(0)
	} else {
		newWindow(openDocument(file)).Show()
	}
}

func (a *app) openBlank() {
	running := getFirstInstance(a)
	if running != nil {
		log.Println("Found running app, opening blank document")
		running.openBlank()
		os.Exit(0)
	} else {
		newWindow(newDocument()).Show()
	}
}
