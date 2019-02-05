package main

import "testing"

func TestLogic(t *testing.T) {
	if true == false {
		t.Error("It's illogical")
	}
}
