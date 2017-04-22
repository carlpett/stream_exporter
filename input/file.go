package input

import (
	"bufio"
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
	registerInput(inputTypeFile, newFileInput)
}

func newFileInput() StreamInput {
	if *filePath == "" {
		panic("No file path set")
	}

	return FileInput{
		path: *filePath,
	}
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
