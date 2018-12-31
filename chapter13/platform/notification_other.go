// +build !darwin

package main

import "log"

func notify(title, content string) {
	log.Println("Notifications not supported")
}
