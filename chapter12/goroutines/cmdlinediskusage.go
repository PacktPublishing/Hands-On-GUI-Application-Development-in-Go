package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/PacktPublishing/Hands-On-GUI-Application-Development-in-Go/chapter12/goroutines/disk"
)

func reportUsage(path string) {
	f, _ := os.Open(path)
	files, _ := f.Readdir(-1)
	f.Close()

	result := make(chan disk.SizeInfo)
	for _, file := range files {
		go disk.DirSize(filepath.Join(path, file.Name()), result)
	}

	var total int64
	results := 0
	for info := range result {
		total += info.Size
		fmt.Printf("%s:\t%s\n", info.Name, disk.FormatSize(info.Size))

		results++
		if results == len(files) {
			break
		}
	}
	fmt.Printf("\nTotal:\t\t%s\n", disk.FormatSize(total))
}

func main() {
	path, _ := os.Getwd()

	if len(os.Args) == 2 {
		path = os.Args[1]
	}

	fmt.Println("Scanning", path)
	reportUsage(path)
}
