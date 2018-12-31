package main

import "log"

type Storage interface {
	Read(name string) string
	Write(name, content string)
}

type diskStore struct {
}

func (*diskStore) Read(name string) string {
	log.Println("TODO Read an actual file")
	return "TODO"
}

func (*diskStore) Write(name, content string) {
	log.Println("TODO write to an actual file")
}

func NewPermanentStorage() Storage {
	return &diskStore{}
}
