package main

import (
	"log"
	"os/exec"
)

func notify(title, content string) {
	cmd := exec.Command("osascript", "-e", "display notification \""+content+
		"\" with title \""+title+"\"")
	err := cmd.Run()

	if err != nil {
		log.Printf("Error showing notification: %v", err)
	}
}
