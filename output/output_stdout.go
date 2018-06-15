package output

import "fmt"

// StdOutput used for debugging, prints all incoming requests
type StdOutput struct {
}

// NewStdOutput constructor for StdOutput
func NewStdOutput() (i *StdOutput) {
	i = new(StdOutput)
	return
}

func (i *StdOutput) Write(data []byte) (int, error) {
	fmt.Println(string(data))
	return len(data), nil
}

func (i *StdOutput) String() string {
	return "Stdout Output"
}
