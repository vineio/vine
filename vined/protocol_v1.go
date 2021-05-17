package vined

import (
	"net"

	log "github.com/donnie4w/go-logger/logger"
)

const (
	frameTypeResponse int32 = 0
	frameTypeError    int32 = 1
	frameTypeMessage  int32 = 2
)

type protocolV1 struct {
	ctx *context
}

func (p *protocolV1) IOLoop(conn net.Conn) error {
	p.Read(conn)
	return nil
}

func (p *protocolV1) Read(conn net.Conn) {
	for {
		bs := make([]byte, 1024)
		n, err := conn.Read(bs)
		if err != nil {
			log.Error("(c *Client) Read:", err)
			return
		}

		if n <= 0 {
			continue
		}
		var msg Message
		msg.Data = bs[:n]
		msg.NodeName = p.ctx.vined.getOpts().BroadcastAddress
		msg.PortName = p.ctx.vined.ApiTCPConn[conn]

		p.ctx.vined.MsgTxChan[msg.PortName] <- msg.Data //put the user's data to tx buffer
		log.Debug("user send data：", msg)
	}
}

func (p *protocolV1) Write(conn net.Conn) {
	portName := p.ctx.vined.ApiTCPConn[conn]
	for {
		data := <-p.ctx.vined.MsgRxChan[portName] //send the serial data to user
		conn.Write(data)
	}
}
