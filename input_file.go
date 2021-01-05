package main

import (
	"encoding/gob"
	"log"
	"os"
	"time"
)

// FileInput ...
type FileInput struct {
	data    chan []byte
	path    string
	decoder *gob.Decoder
}

// NewFileInput ...
func NewFileInput(path string) (i *FileInput) {
	i = new(FileInput)
	i.data = make(chan []byte)
	i.path = path
	i.Init(path)

	go i.emit()

	return
}

// Init ...
func (i *FileInput) Init(path string) {
	file, err := os.Open(path)

	if err != nil {
		log.Fatal(i, "Cannot open file %q. Error: %s", path, err)
	}

	i.decoder = gob.NewDecoder(file)
}

func (i *FileInput) Read(data []byte) (int, error) {
	buf := <-i.data
	copy(data, buf)

	return len(buf), nil
}

func (i *FileInput) String() string {
	return "File input: " + i.path
}

func (i *FileInput) emit() {
	var lastTime int64

	for {
		raw := new(RawRequest)
		err := i.decoder.Decode(raw)

		if err != nil {
			return
		}

		if lastTime != 0 {
			time.Sleep(time.Duration(raw.Timestamp - lastTime))
		}
		lastTime = raw.Timestamp

		i.data <- raw.Request
	}
}
