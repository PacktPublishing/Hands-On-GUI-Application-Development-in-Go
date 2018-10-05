package main

import (
	"fmt"
	"time"
)

func tick() {
	fmt.Println("Started")

	for _ = range time.NewTicker(time.Second).C {
		fmt.Print(".")
	}

	fmt.Println("Ended")
}

func main() {
	go tick()
	go tick()
	go tick()

	time.Sleep(5 * time.Second)
}
