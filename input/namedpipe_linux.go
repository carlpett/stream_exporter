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

func newNamedPipeInput(config InputConfig) StreamInput {
	return NamedPipeInput{
		config: config,
	}
}

func (socket SocketInput) StartStream(ch chan<- string) {
	err := syscall.Mkfifo(socket.config.PipePath, 0666)
	if err != nil {
		panic(err)
	}

	pipe, err := os.OpenFile(socket.config.PipePath, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		panic(err)
	}
	defer os.Remove(path)

	reader := bufio.NewReader(pipe)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		ch <- scanner.Text()
	}
}
