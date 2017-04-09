package input

type inputType string

const (
	inputTypeFile      = "file"
	inputTypeNamedPipe = "namedpipe"
	inputTypeSocket    = "socket"
)

type InputConfig struct {
	Type                 string
	FileInputConfig      `mapstructure:",squash"`
	NamedPipeInputConfig `mapstructure:",squash"`
	SocketInputConfig    `mapstructure:",squash"`
}
