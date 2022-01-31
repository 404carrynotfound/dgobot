package main

import "testing"

func TestValidURL(t *testing.T) {
	if !validURL("https://www.youtube.com/watch?v=dQw4w9WgXcQ") {
		t.Error("String is valid.")
	}

	if validURL("foo") {
		t.Error("String is invalid")
	}
}
