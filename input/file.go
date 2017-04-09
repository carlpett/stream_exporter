package input

import (
	"bufio"
	"os"
)

type FileInputConfig struct {
	FilePath string `mapstructure:"file_path"`
}

type FileInput struct {
	config FileInputConfig
}

func init() {
	registerInput(inputTypeFile, newFileInput)
}

func newFileInput(config InputConfig) StreamInput {
	return FileInput{
		config: config.FileInputConfig,
	}
}

func (input FileInput) StartStream(ch chan<- string) {
	defer close(ch)

	file, err := os.Open(input.config.FilePath)
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
