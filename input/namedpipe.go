package input

import (
	"bufio"
)

type NamedPipeInputConfig struct {
	PipePath string `mapstructure:"namedpipe_path"`
}

type NamedPipeInput struct {
	config  NamedPipeInputConfig
	scanner *bufio.Scanner
}
