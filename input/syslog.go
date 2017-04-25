package input

import (
	"errors"
	"flag"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mcuadros/go-syslog.v2"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

type SyslogInput struct {
	listenAddr   string
	listenFamily string
	format       format.Format
}

var (
	syslogListenFamily = flag.String("input.syslog.listenfamily", "", "Listening protocol family (tcp/udp/unix)")
	syslogListenAddr   = flag.String("input.syslog.listenaddr", "", "Listening address of syslog server")
	syslogFormatFlag   = flag.String("input.syslog.format", "autodetect", "Format of incoming syslog data (rfc3164/rfc5424/rfc6587/autodetect)")
)

func init() {
	registerInput("syslog", newSyslogInput)
}

func newSyslogInput() (StreamInput, error) {
	if *syslogListenFamily == "" {
		return nil, errors.New("-input.syslog.listenfamily not set")
	} else if *syslogListenFamily != "tcp" && *syslogListenFamily != "udp" && *syslogListenFamily != "unix" {
		return nil, errors.New(fmt.Sprintf("%q is not a valid value for -input.syslog.listenfamily", *syslogListenFamily))
	}

	if *syslogListenAddr == "" {
		return nil, errors.New("-input.syslog.listenaddr not set")
	}

	var syslogFormat format.Format
	switch *syslogFormatFlag {
	case "":
		syslogFormat = syslog.Automatic
	case "autodetect":
		syslogFormat = syslog.Automatic
	case "rfc3164":
		syslogFormat = syslog.RFC3164
	case "rfc5424":
		syslogFormat = syslog.RFC5424
	case "rfc6587":
		syslogFormat = syslog.RFC6587
	default:
		return nil, errors.New(fmt.Sprintf("%q is not a valid value for -input.syslog.format", *syslogFormatFlag))
	}

	return SyslogInput{
		listenAddr:   *syslogListenAddr,
		listenFamily: *syslogListenFamily,
		format:       syslogFormat,
	}, nil
}

func (input SyslogInput) StartStream(ch chan<- string) {
	syslogChannel := make(syslog.LogPartsChannel)
	logHandler := syslog.NewChannelHandler(syslogChannel)

	server := syslog.NewServer()

	var err error
	switch input.listenFamily {
	case "tcp":
		err = server.ListenTCP(input.listenAddr)
	case "udp":
		err = server.ListenUDP(input.listenAddr)
	case "unix":
		err = server.ListenUnixgram(input.listenAddr)
	default:
		log.Fatalf("Unknown listen family %q", input.listenFamily)
	}
	if err != nil {
		log.Fatal(err)
	}

	server.SetHandler(logHandler)
	server.SetFormat(input.format)

	err = server.Boot()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Syslog server started listening at %s", input.listenAddr)

	go func(lineIn syslog.LogPartsChannel) {
		for parts := range lineIn {
			ch <- parts["content"]
			/* parts is a map with the following keys:
			-hostname
			-tag
			-content (Message part)
			-facility
			-tls_peer
			-timestamp
			-severity
			-client
			-priority
			Should these be joined up? Should it be possible to ask for a particular field? User-specified Sprintf-format?
			For now, just send the message itself.
			*/
		}
	}(syslogChannel)

	server.Wait()
	log.Info("Syslog server shutting down")
}
