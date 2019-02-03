package input

import (
	"bufio"
	"os"
)

type StdinInput struct {
}

func init() {
	registerInput("stdin", newStdinInput)
}

func newStdinInput() (StreamInput, error) {
	return StdinInput{}, nil
}

func (input StdinInput) StartStream(ch chan<- string) {
	reader := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		ch <- scanner.Text()
	}
}
