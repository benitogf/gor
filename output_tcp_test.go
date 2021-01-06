package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"sync"
	"testing"
)

func TestTCPOutput(t *testing.T) {
	wg := new(sync.WaitGroup)
	quit := make(chan int)

	listener := startTCP(func(data []byte) {
		wg.Done()
	})
	input := NewTestInput()
	output := NewTCPOutput(listener.Addr().String())

	Plugins.Inputs = []io.Reader{input}
	Plugins.Outputs = []io.Writer{output}

	go StartEmmiter(quit)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		input.EmitGET()
	}

	wg.Wait()

	close(quit)
}

func startTCP(cb func([]byte)) net.Listener {
	listener, err := net.Listen("tcp", ":0")

	if err != nil {
		log.Fatal("Can't start:", err)
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("failed TCP connection", err)
			}

			go func() {
				reader := bufio.NewReader(conn)
				for {
					buf, err := reader.ReadBytes('¶')
					newBufLen := len(buf) - 2
					newBuf := make([]byte, newBufLen)
					copy(newBuf, buf[:newBufLen])
					if err != nil {
						log.Println("error reading bytes from TCP", err)
						if err != io.EOF {
							conn.Close()
							break
						}
					}
					cb(newBuf)
				}
			}()
		}
	}()

	return listener
}

func BenchmarkTCPOutput(b *testing.B) {
	wg := new(sync.WaitGroup)
	quit := make(chan int)

	listener := startTCP(func(data []byte) {
		wg.Done()
	})
	input := NewTestInput()
	output := NewTCPOutput(listener.Addr().String())

	Plugins.Inputs = []io.Reader{input}
	Plugins.Outputs = []io.Writer{output}

	go StartEmmiter(quit)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		input.EmitGET()
	}

	wg.Wait()

	close(quit)
}
