package input

import (
	"bufio"
	"io"
	"os"
	"syscall"
)

func init() {
	registerInput("namedpipe", newNamedPipeInput)
}

var (
	pipePath = flag.String("input.namedpipe.path", "", "Path where pipe should be created")
)

func newNamedPipeInput() StreamInput {
	return NamedPipeInput{
		path: *pipePath,
	}
}

func (input NamedPipeInput) StartStream(ch chan<- string) {
	err := syscall.Mkfifo(input.path, 0666)
	if err != nil {
		panic(err)
	}

	pipe, err := os.OpenFile(input.path, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		panic(err)
	}
	defer os.Remove(input.path)

	reader := bufio.NewReader(pipe)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		ch <- scanner.Text()
	}
}
