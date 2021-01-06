package main

import (
	"bufio"
	"io"
	"log"
	"net"
)

// TCPInput ...
// Can be tested using nc tool:
//    echo "asdad" | nc 127.0.0.1 27017
//
type TCPInput struct {
	data     chan []byte
	address  string
	listener net.Listener
}

// NewTCPInput ...
func NewTCPInput(address string) (i *TCPInput) {
	i = new(TCPInput)
	i.data = make(chan []byte)
	i.address = address

	i.listen(address)

	return
}

func (i *TCPInput) Read(data []byte) (int, error) {
	buf := <-i.data
	copy(data, buf)

	return len(buf), nil
}

func (i *TCPInput) listen(address string) {
	listener, err := net.Listen("tcp", address)
	i.listener = listener

	if err != nil {
		log.Fatal("Can't start:", err)
	}

	go func() {
		for {
			conn, err := listener.Accept()

			if err != nil {
				log.Println("Error while Accept()", err)
				continue
			}

			go i.handleConnection(conn)
		}
	}()
}

func (i *TCPInput) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		buf, err := reader.ReadBytes('¶')
		bufLen := len(buf)
		if bufLen > 0 {
			newBufLen := len(buf) - 2
			if newBufLen > 0 {
				newBuf := make([]byte, newBufLen)
				copy(newBuf, buf[:newBufLen])
				i.data <- newBuf
				if err != nil {
					log.Println("tcp connection error: ", err)
					if err != io.EOF {
						break
					}
				}
			}
		}
	}
}

func (i *TCPInput) String() string {
	return "TCP input: " + i.address
}
