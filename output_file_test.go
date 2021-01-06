package main

import (
	"io"
	"sync"
	"testing"
)

func TestFileOutput(t *testing.T) {
	wg := new(sync.WaitGroup)
	quit := make(chan int)

	input := NewTestInput()
	output := NewFileOutput("/tmp/test_requests.gor")

	Plugins.Inputs = []io.Reader{input}
	Plugins.Outputs = []io.Writer{output}

	go StartEmmiter(quit)

	for i := 0; i < 100; i++ {
		wg.Add(2)
		input.EmitGET()
		input.EmitPOST()
	}
	close(quit)

	quit = make(chan int)

	input2 := NewFileInput("/tmp/test_requests.gor")
	output2 := NewTestOutput(func(data []byte) {
		wg.Done()
	})

	Plugins.Inputs = []io.Reader{input2}
	Plugins.Outputs = []io.Writer{output2}

	go StartEmmiter(quit)

	wg.Wait()
	close(quit)
}
