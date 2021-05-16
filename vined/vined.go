package vined

import (
	"io"
	"net"
	"os"
	"strings"
	"sync/atomic"

	"github.com/vineio/vine/internal/protocol"
	"github.com/vineio/vine/serial"
	"github.com/vineio/vine/util"

	log "github.com/donnie4w/go-logger/logger"
)

type VINED struct {
	clientIDSequence int64
	tcpListener      net.Listener

	opts      atomic.Value
	waitGroup util.WaitGroupWrapper

	serialIO       map[string]io.ReadWriteCloser
	ClientConnChan map[string]chan string
}

func New(opts *Options) *VINED {

	v := &VINED{
		serialIO:       make(map[string]io.ReadWriteCloser),
		ClientConnChan: make(map[string]chan string),
	}
	v.swapOpts(opts)
	return v
}

func (v *VINED) Main() {

	var err error
	v.tcpListener, err = net.Listen("tcp", v.getOpts().TCPAddress)
	if err != nil {
		log.Error("listen (%s) failed - %s", v.getOpts().TCPAddress, err)
		os.Exit(1)
	}

	var ctx = &context{v}

	log.Debug("tcp listen on:", v.getOpts().TCPAddress)
	tcpServer := &tcpServer{ctx: ctx}
	v.waitGroup.Wrap(func() {
		protocol.TCPServer(v.tcpListener, tcpServer)
	})

	ios, err := serial.New(v.getOpts().Optserials)
	if err != nil {
		log.Error("serial.New():", err)
		os.Exit(1)
	}
	log.Debug("serial opeded on:", v.getOpts().Optserials)

	v.serialIO = ios
	client := &Client{ctx, make(chan Message, 1024)}
	v.waitGroup.Wrap(func() {

		n := strings.Index(v.getOpts().TCPAddress, ":")
		tcpServerPort := v.getOpts().TCPAddress[n:]
		client.Handle(ios, tcpServerPort)
	})

	v.waitGroup.Wait()
}

func (v *VINED) getOpts() *Options {
	return v.opts.Load().(*Options)
}

func (v *VINED) swapOpts(opts *Options) {
	v.opts.Store(opts)
}
