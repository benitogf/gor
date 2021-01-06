package main

import (
	"io"
	"sync"
	"testing"
)

func TestLimiter(t *testing.T) {
	wg := new(sync.WaitGroup)
	quit := make(chan int)

	input := NewTestInput()
	output := NewLimiter(NewTestOutput(func(data []byte) {
		wg.Done()
	}), 10)
	wg.Add(10)

	Plugins.Inputs = []io.Reader{input}
	Plugins.Outputs = []io.Writer{output}

	go StartEmmiter(quit)

	for i := 0; i < 100; i++ {
		input.EmitGET()
	}

	wg.Wait()

	close(quit)
}
