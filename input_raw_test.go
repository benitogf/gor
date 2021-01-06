package main

import (
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestRAWInput(t *testing.T) {

	wg := new(sync.WaitGroup)
	quit := make(chan int)

	listener := startHTTP(func(req *http.Request) {})

	input := NewRAWInput(listener.Addr().String())
	output := NewTestOutput(func(data []byte) {
		wg.Done()
	})

	Plugins.Inputs = []io.Reader{input}
	Plugins.Outputs = []io.Writer{output}

	address := strings.Replace(listener.Addr().String(), "[::]", "127.0.0.1", -1)

	go StartEmmiter(quit)
	// wait for the emmiter to start TODO: make the state available, maybe make the function block until start?
	// without the sleep the test is flaky
	time.Sleep(500 * time.Millisecond)

	for i := 0; i < 100; i++ {
		res, err := http.Get("http://" + address)
		if err != nil {
			log.Fatal("request failed", err)
		}
		wg.Add(1)
		res.Body.Close()
	}

	wg.Wait()

	close(quit)
}
