package main

import (
	"bytes"
	"os"
	"testing"
)

func Test_process(t *testing.T) {
	input, err := os.Open("../../enwiktionary-2000.xml")
	if err != nil {
		t.Fatal(err)
	}

	writer := &bytes.Buffer{}
	process(input, writer)
}
