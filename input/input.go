package input

type StreamInput interface {
	StartStream(ch chan<- string)
}

var inputTypes = make(map[string]func() StreamInput)

func registerInput(inputType string, factory func() StreamInput) {
	inputTypes[inputType] = factory
}

func GetAvailableInputs() []string {
	inputs := make([]string, 0, len(inputTypes))
	for key := range inputTypes {
		inputs = append(inputs, key)
	}
	return inputs
}

func NewInput(inputType string) StreamInput {
	factory, registered := inputTypes[inputType]
	if !registered {
		panic("Unknown input type")
	}
	return factory()
}
