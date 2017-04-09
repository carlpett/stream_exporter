package input

type inputType string

const (
	inputTypeFile      = "file"
	inputTypeNamedPipe = "namedpipe"
)

type InputConfig struct {
	Type                 string
	FileInputConfig      `mapstructure:",squash"`
	NamedPipeInputConfig `mapstructure:",squash"`
}
