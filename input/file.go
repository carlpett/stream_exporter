package input

import (
	"bufio"
	"io"
	"os"
)

type FileInput struct {
	scanner *bufio.Scanner
}

func newFileInput(config InputConfig) FileInput {
	file, err := os.Open("/tmp/myfile")
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	return FileInput{scanner: scanner}
}

func (file FileInput) ReadLine() (string, error) {
	if file.scanner.Scan() {
		return file.scanner.Text(), nil
	} else {
		return "", io.EOF
	}
}
