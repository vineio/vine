package vined

import (
	"crypto/md5"
	"fmt"
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

const (
	MAX_TOTAL_CHAN_SIZE  = 1024
	MAX_SINGLE_CHAN_SIZE = 1024
)

type VINED struct {
	clientIDSequence int64
	tcpListener      net.Listener
	apiTcpListener   net.Listener
	opts             atomic.Value
	waitGroup        util.WaitGroupWrapper

	serialIO           map[string]io.ReadWriteCloser
	serialMd5          map[string][md5.Size]byte //key string is serial port name ,value is the name's md5 value
	ClientConnChan     map[string]chan string    //key string is:client ip address,value string is serial portname
	MsgTxChan          map[string]chan []byte    //key string is:serial port name,value is rx or tx serial data
	MsgRxChan          map[string]chan []byte
	MessageTotalRxChan chan Message
	MessageTotalTxChan chan Message
	ApiTCPConn         map[net.Conn]string //key string si :api tcp client conn,value is serial port name
}

func New(opts *Options) *VINED {

	v := &VINED{
		serialIO:           make(map[string]io.ReadWriteCloser),
		serialMd5:          make(map[string][md5.Size]byte),
		ClientConnChan:     make(map[string]chan string),
		ApiTCPConn:         make(map[net.Conn]string),
		MsgTxChan:          make(map[string]chan []byte),
		MsgRxChan:          make(map[string]chan []byte),
		MessageTotalRxChan: make(chan Message, MAX_TOTAL_CHAN_SIZE),
		MessageTotalTxChan: make(chan Message, MAX_TOTAL_CHAN_SIZE),
	}

	v.swapOpts(opts)
	return v
}

func (v *VINED) Main() {

	var err error
	var ctx = &context{v}

	v.tcpListener, err = net.Listen("tcp", v.getOpts().TCPAddress)
	if err != nil {
		log.Error(fmt.Sprintf("listen (%s) failed - %s", v.getOpts().TCPAddress, err))
		os.Exit(1)
	}
	v.apiTcpListener, err = net.Listen("tcp", v.getOpts().ApiTCPAddress)
	if err != nil {
		log.Error(fmt.Sprintf("listen (%s) failed - %s", v.getOpts().ApiTCPAddress, err))
		os.Exit(1)
	}

	log.Debug("tcp listen on:", v.getOpts().TCPAddress)
	tcpServer := &tcpServer{ctx: ctx}
	v.waitGroup.Wrap(func() {
		protocol.TCPServer(v.tcpListener, tcpServer)
	})

	log.Debug("api tcp listen on:", v.getOpts().ApiTCPAddress)
	apiTcpServer := &apiTcpServer{ctx: ctx}
	v.waitGroup.Wrap(func() {
		protocol.TCPServer(v.apiTcpListener, apiTcpServer)
	})

	ios, md5s, err := serial.New(v.getOpts().Optserials)
	if err != nil {
		log.Error("serial.New():", err)
		os.Exit(1)
	}
	log.Debug(md5s)
	log.Debug("serial opeded on:", v.getOpts().Optserials)

	v.serialMd5 = md5s
	v.serialIO = ios
	client := &Client{ctx}
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
