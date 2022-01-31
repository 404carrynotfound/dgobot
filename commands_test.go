package main

import "testing"

// Checks if for every command there's a function to handle that
func TestCommands(t *testing.T) {
	for _, command := range commands {
		if commandHandlers[command.Name] == nil {
			t.Errorf("Declared command %s in application command slice, but there's no handler.", command.Name)
		}
	}
}
