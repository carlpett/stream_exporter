package input

type StreamInput interface {
	ReadLine() (string, error)
}

func NewInput(config InputConfig) StreamInput {
	return newFileInput(config)
}
