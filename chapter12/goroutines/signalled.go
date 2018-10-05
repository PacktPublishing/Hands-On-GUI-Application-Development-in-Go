package main

import (
	"fmt"
	"time"
)

func tickUntil(stop chan (struct{})) {
	fmt.Println("Started")

	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case <-ticker:
			fmt.Print(".")
		case <-stop:
			fmt.Println("Ended")
			return
		}
	}
}

func main() {
	stop := make(chan (struct{}))

	go tickUntil(stop)
	go tickUntil(stop)
	go tickUntil(stop)

	time.Sleep(5 * time.Second)
	close(stop)

	time.Sleep(10 * time.Millisecond)
}
