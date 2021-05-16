package vined

import (
	"net"

	"github.com/vineio/vine/util"

	log "github.com/donnie4w/go-logger/logger"
)

type tcpServer struct {
	ctx *context
}

func (t *tcpServer) Handle(clientConn net.Conn) {
	portName := <-t.ctx.vined.ClientConnChan[clientConn.RemoteAddr().String()] //wait for client to connect

	log.Debug("tunneled the serial port to net,serial port name,net ip :", portName, clientConn.RemoteAddr().String())
	util.Join(clientConn, t.ctx.vined.serialIO[portName]) //tunnel the serial to net

}
