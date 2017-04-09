package input

import (
	"bufio"
	"io"
	"os"
)

type FileInputConfig struct {
	FilePath string `mapstructure:"file_path"`
}

type FileInput struct {
	scanner *bufio.Scanner
}

func init() {
	registerInput(inputTypeFile, newFileInput)
}

func newFileInput(config InputConfig) StreamInput {
	file, err := os.Open(config.FilePath)
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
