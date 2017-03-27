package input

import (
	"bufio"
	"io"
	"os"
	"syscall"
)

func rmfifo() {
	os.Remove("/tmp/myfifo")
}

func newNamedPipeInput(config InputConfig) NamedPipeInput {
	err := syscall.Mkfifo("/tmp/myfifo", 0666)
	if err != nil {
		panic(err)
	}

	pipe, err := os.OpenFile("/tmp/myfifo", os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(pipe)
	scanner := bufio.NewScanner(reader)

	return NamedPipeInput{scanner: scanner}
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
	rmfifo()
}
