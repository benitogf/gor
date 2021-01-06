package main

import (
	"errors"
	"io"
)

// StartEmmiter ...
func StartEmmiter(stop chan int) {
	for _, in := range Plugins.Inputs {
		go Broadcast(in, Plugins.Outputs...)
	}

	select {
	case <-stop:
		return
	}
}

// Broadcast from 1 reader to multiple writers
func Broadcast(src io.Reader, writers ...io.Writer) (err error) {
	buf := make([]byte, 32*1024)
	wIndex := 0

	for {
		nr, er := src.Read(buf)
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
		if nr == 0 {
			continue
		}
		if len(buf) <= nr {
			err = errors.New("tcp buffer exceeded")
			break
		}

		Debug("Sending", src, ": ", string(buf[0:nr]))
		if Settings.splitOutput {
			// Simple round robin
			writers[wIndex].Write(buf[0:nr])

			wIndex++

			if wIndex >= len(writers) {
				wIndex = 0
			}
			continue
		}

		for _, dst := range writers {
			dst.Write(buf[0:nr])
		}
	}
	return err
}
