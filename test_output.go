package main

type writeCallback func(data []byte)

// TestOutput ...
type TestOutput struct {
	cb writeCallback
}

// NewTestOutput ...
func NewTestOutput(cb writeCallback) (i *TestOutput) {
	i = new(TestOutput)
	i.cb = cb

	return
}

func (i *TestOutput) Write(data []byte) (int, error) {
	i.cb(data)

	return len(data), nil
}

func (i *TestOutput) String() string {
	return "Test Input"
}
