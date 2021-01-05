package main

import (
	"fmt"
)

// DummyOutput ...
type DummyOutput struct {
}

// NewDummyOutput ...
func NewDummyOutput(options string) (di *DummyOutput) {
	di = new(DummyOutput)

	return
}

func (i *DummyOutput) Write(data []byte) (int, error) {
	fmt.Println("Writing message: ", data)

	return len(data), nil
}

func (i *DummyOutput) String() string {
	return "Dummy Output"
}
