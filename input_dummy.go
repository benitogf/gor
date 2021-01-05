package main

import (
	"time"
)

// DummyInput ...
type DummyInput struct {
	data chan []byte
}

// NewDummyInput ...
func NewDummyInput(options string) (di *DummyInput) {
	di = new(DummyInput)
	di.data = make(chan []byte)

	go di.emit()

	return
}

func (i *DummyInput) Read(data []byte) (int, error) {
	buf := <-i.data
	copy(data, buf)

	return len(buf), nil
}

func (i *DummyInput) emit() {
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ticker.C:
			i.data <- []byte("GET / HTTP/1.1\r\n\r\n")
		}
	}
}

func (i *DummyInput) String() string {
	return "Dummy Input"
}
