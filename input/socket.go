package input

import (
	"bufio"
	"flag"
	"net"
)

type SocketInput struct {
	family     string
	listenAddr string
}

var (
	socketFamily     = flag.String("input.socket.family", "tcp", "Socket family (tcp/udp/etc)")
	socketListenAddr = flag.String("input.socket.listenaddr", "", "Listening address of socket")
)

func init() {
	registerInput(inputTypeSocket, newSocketInput)
}

func newSocketInput() StreamInput {
	if *socketFamily == "" {
		panic("Socket family not set")
	} else if *socketFamily != "tcp" && *socketFamily != "udp" && *socketFamily != "domain" {
		panic("Invalid socket family")
	}

	if *socketListenAddr == "" {
		panic("Socket listening address not set")
	}

	return SocketInput{
		family:     *socketFamily,
		listenAddr: *socketListenAddr,
	}
}

func (socket SocketInput) StartStream(ch chan<- string) {
	l, err := net.Listen(socket.family, socket.listenAddr)
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
