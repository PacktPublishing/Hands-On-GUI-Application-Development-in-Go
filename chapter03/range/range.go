package main

import "fmt"

func download(file string, c chan int) {
	fmt.Println("Downloading", file)

	c <- 10
	c <- 40
	c <- 65
	c <- 100

	close(c)
}

func main() {
	c := make(chan int)
	go download("myfile.jpg", c)

	for i := range c {
		fmt.Printf("Progress %d%%...\n", i)
	}
	fmt.Println("Download complete")
}
