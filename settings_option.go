package main

import (
	"fmt"
)

// MultiOption ...
type MultiOption []string

// String ...
func (h *MultiOption) String() string {
	return fmt.Sprint(*h)
}

// Set ...
func (h *MultiOption) Set(value string) error {
	*h = append(*h, value)
	return nil
}
