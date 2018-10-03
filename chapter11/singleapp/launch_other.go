// +build !darwin

package main

type app struct {
}

func (a *app) openFile(file string) {
	newWindow(openDocument(file)).Show()
}

func (a *app) openBlank() {
	newWindow(newDocument()).Show()
}
