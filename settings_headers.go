package main

import (
	"errors"
	"fmt"
	"strings"
)

// HTTPHeaders ...
type HTTPHeaders []HTTPHeader

// HTTPHeader ...
type HTTPHeader struct {
	Name  string
	Value string
}

func (h *HTTPHeaders) String() string {
	return fmt.Sprint(*h)
}

// Set ...
func (h *HTTPHeaders) Set(value string) error {
	v := strings.SplitN(value, ":", 2)
	if len(v) != 2 {
		return errors.New("Expected `Key: Value`")
	}

	header := HTTPHeader{
		strings.TrimSpace(v[0]),
		strings.TrimSpace(v[1]),
	}

	*h = append(*h, header)
	return nil
}
