package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type SingleApp struct {
	app App
}

func (a *SingleApp) OpenFile(file string, _ *Doc) error {
	newWindow(openDocument(file)).Show()
	return nil
}

func (a *SingleApp) OpenBlank(string, *Doc) error {
	newWindow(newDocument()).Show()
	return nil
}

func newServer(app App) {
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}

	server := &SingleApp{app}
	rpc.Register(server)
	rpc.HandleHTTP()
	go http.Serve(l, nil)
}

func getServer() *rpc.Client {
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Println("dialing:", err)
		return nil
	}

	return client
}

type wrapper struct {
	client *rpc.Client
}

func (w *wrapper) openFile(file string) {
	var reply Doc
	err := w.client.Call("SingleApp.OpenFile", file, &reply)

	if err != nil {
		log.Fatal("Open error", err)
	}
}

func (w *wrapper) openBlank() {
	var reply Doc
	err := w.client.Call("SingleApp.OpenBlank", "", &reply)

	if err != nil {
		log.Fatal("Blank error", err)
	}
}

func newWrapper(c *rpc.Client) *wrapper {
	return &wrapper{c}
}

func getFirstInstance(a App) App {
	s := getServer()
	if s != nil {
		return newWrapper(s)
	}

	newServer(a)
	return nil
}
