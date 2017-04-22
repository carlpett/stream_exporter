package input

import (
	"bufio"
	"errors"
	"flag"
	"os"
)

type FileInput struct {
	path string
}

var (
	filePath = flag.String("input.file.path", "", "Path to file to read")
)

func init() {
	registerInput("file", newFileInput)
}

func newFileInput() (StreamInput, error) {
	if *filePath == "" {
		return nil, errors.New("-input.file.path not set")
	}

	return FileInput{
		path: *filePath,
	}, nil
}

func (input FileInput) StartStream(ch chan<- string) {
	defer close(ch)

	file, err := os.Open(input.path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		ch <- scanner.Text()
	}
}
