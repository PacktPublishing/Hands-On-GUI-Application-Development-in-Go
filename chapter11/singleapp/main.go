package main

import (
	"log"
	"os"
	"time"
)

// App is the main interface for our application functions that will be
// made available over RPC so that other instances can call them.
type App interface {
	openFile(string)
	openBlank()
}

type Doc struct {
	File string
}

func openDocument(file string) *Doc {
	return &Doc{file}
}

func newDocument() *Doc {
	return &Doc{}
}

type win struct {
	doc *Doc
}

func (w *win) Show() {
	if w.doc.File == "" {
		log.Println("Showing window for blank document")
	} else {
		log.Println("Showing window for document", w.doc.File)
	}

}

func newWindow(doc *Doc) *win {
	return &win{doc}
}

func main() {
	app := &app{}

	if len(os.Args) <= 1 {
		app.openBlank()
	} else {
		app.openFile(os.Args[1])
	}

	time.Sleep(10 * time.Second)
}
