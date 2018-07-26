package main

import "fmt"

func say(words string, to chan string) {
	fmt.Println("Speaking:", words)
	to <- words
}

func talk(to chan string) {
	say("Hello", to)
	say("Everyone", to)
	say("My name is...", to)
	fmt.Println("Never mind")
}

func listen(to chan string) {
	heard := <-to
	fmt.Println("I heard:", heard)
}

func main() {
	chat := make(chan string, 3)

	go talk(chat)

	listen(chat)
	listen(chat)
	fmt.Println("Bye")
}
