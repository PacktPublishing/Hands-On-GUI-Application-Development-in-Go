package main

import (
	"fmt"
	"time"
)

func tick(message string) {
	for i := 0; i < 5; i++ {
		time.Sleep(10 * time.Millisecond)
		fmt.Println(message)
	}
}

func main() {
	go tick("goroutine")
	tick("function")
}
