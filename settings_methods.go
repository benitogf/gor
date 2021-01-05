package main

import (
	"fmt"
	"strings"
)

// HTTPMethods ...
type HTTPMethods []string

func (h *HTTPMethods) String() string {
	return fmt.Sprint(*h)
}

// Set ...
func (h *HTTPMethods) Set(value string) error {
	*h = append(*h, strings.ToUpper(value))
	return nil
}

// Contains ...
func (h *HTTPMethods) Contains(value string) bool {
	for _, method := range *h {
		if value == method {
			return true
		}
	}
	return false
}
