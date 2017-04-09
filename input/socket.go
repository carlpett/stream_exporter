package input

import (
	"bufio"
	"net"
)

type SocketInputConfig struct {
	SocketType string `mapstructure:"socket_type"`
	ListenAddr string `mapstructure:"socket_listenaddr"`
}

type SocketInput struct {
	config SocketInputConfig
}

func init() {
	registerInput(inputTypeSocket, newSocketInput)
}

func newSocketInput(config InputConfig) StreamInput {
	// TODO: Validate config
	return SocketInput{
		config: config.SocketInputConfig,
	}
}

func (socket SocketInput) StartStream(ch chan<- string) {
	l, err := net.Listen(socket.config.SocketType, socket.config.ListenAddr)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go func(c net.Conn) {
			// TODO: Timeout + metrics for failed reads
			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				ch <- scanner.Text()
			}
			c.Close()
		}(conn)
	}
}
