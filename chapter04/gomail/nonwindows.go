// +build !windows

package main

import "log"

func NewMailUIBrowse() {
	log.Fatalln("GoMail with Walk only works on windows")
}
