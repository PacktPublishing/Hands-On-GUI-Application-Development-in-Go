package main

import "fmt"

/*
void loadMenu();
*/
import "C"

func loadMenu() {
	C.loadMenu()
}

//export menuNewClicked
func menuNewClicked() {
	fmt.Println("New clicked")
}
