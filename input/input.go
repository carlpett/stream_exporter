package input

type StreamInput interface {
	StartStream(ch chan<- string)
}

var inputTypes = make(map[string]func(InputConfig) StreamInput)

func registerInput(inputType string, factory func(InputConfig) StreamInput) {
	inputTypes[inputType] = factory
}

func NewInput(config InputConfig) StreamInput {
	inputType := config.Type
	factory, registered := inputTypes[inputType]
	if !registered {
		panic("Unknown input type")
	}
	return factory(config)
}
