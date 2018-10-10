package main

import (
	"encoding/json"
	"fmt"

	"github.com/PacktPublishing/Hands-On-GUI-Application-Development-in-Go/chapter12/network/remote"
)

type Person struct {
	Title     string `json:"title,omitempty"`
	Firstname string `json:"firstname"`
	Surname   string `json:"surname"`

	Username string `json:"username"`
	Password string `json:"-"`
}

func main() {
	fmt.Println("Downloading...")
	stream := remote.ReadStream("http://echo.jsontest.com/title/Sir/" +
		"firstname/Tom/surname/Jones/username/singer1/")
	defer stream.Close()

	person := &Person{}
	json.NewDecoder(stream).Decode(person)
	fmt.Println("Decoded:", person)
}
