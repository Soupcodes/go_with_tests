package main

import (
	"bytes"
	"testing"
)

func TestGreet(t *testing.T) {
	buffer := bytes.Buffer{}
	// If Greet didn't implement Fprintf under the hood, but only implemented Printf, the information would have been written to the wrong writer interface and never be printed back in these testss
	Greet(&buffer, "Chris")

	got := buffer.String()
	want := "Hello, Chris"

	if got != want {
		t.Errorf("Got: %s, Want: %s", got, want)
	}
}
