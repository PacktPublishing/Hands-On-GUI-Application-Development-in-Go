package main

import "log"

func runApp(storage Storage) {
	log.Println("Writing README.txt")
	storage.Write("README.txt", "overwrite")

	log.Println("Reading README.txt")
	log.Println(storage.Read("README.txt"))
}

func main() {
	runApp(NewPermanentStorage())
}
