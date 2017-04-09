package input

type NamedPipeInputConfig struct {
	PipePath string `mapstructure:"namedpipe_path"`
}

type NamedPipeInput struct {
	config NamedPipeInputConfig
}
