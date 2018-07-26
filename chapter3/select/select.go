package main

import "fmt"

func square(c, quit chan int) {
	sq := 2
	for {
		select {
		case c <- sq:
			sq *= sq
		case <-quit:
			fmt.Println("quitting")
			return
		}
	}
}

func main() {
	c := make(chan int)
	quit := make(chan int)
	go square(c, quit)

	func() {
		for i := 0; i < 5; i++ {
			fmt.Println("Square", <-c)
		}
		quit <- 1
	}()
}
