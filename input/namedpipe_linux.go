package input

import (
	"bufio"
	"io"
	"os"
	"syscall"
)

func init() {
	registerInput(inputTypeNamedPipe, newNamedPipeInput)
}

func rmfifo(path string) {
	os.Remove(path)
}

func newNamedPipeInput(config InputConfig) StreamInput {
	err := syscall.Mkfifo(config.PipePath, 0666)
	if err != nil {
		panic(err)
	}

	pipe, err := os.OpenFile(config.PipePath, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(pipe)
	scanner := bufio.NewScanner(reader)

	return NamedPipeInput{
		scanner: scanner,
		config:  config,
	}
}

func (pipe NamedPipeInput) ReadLine() (string, error) {
	if pipe.scanner.Scan() {
		return pipe.scanner.Text(), nil
	} else {
		return "", io.EOF
	}
}

// TODO
func (pipe NamedPipeInput) Close() {
	rmfifo(pipe.config.Path)
}
