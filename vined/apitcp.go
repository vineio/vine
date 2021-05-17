package vined

import (
	"crypto/md5"
	"fmt"
	"io"
	"net"

	log "github.com/donnie4w/go-logger/logger"
	"github.com/vineio/vine/internal/protocol"
)

type apiTcpServer struct {
	ctx *context
}

func (t *apiTcpServer) Handle(clientConn net.Conn) {
	log.Info(fmt.Sprintf("API TCP: new client(%s)", clientConn.RemoteAddr()))

	buf := make([]byte, 16)
	_, err := io.ReadFull(clientConn, buf)

	if err != nil {
		log.Info(fmt.Sprintf("failed to read protocol serial port md5 value - %s", err))
		return
	}
	var protocolMagic [md5.Size]byte
	copy(protocolMagic[:], buf)
	log.Info(fmt.Sprintf("CLIENT(%v): desired protocol magic '%v'", clientConn.RemoteAddr(), protocolMagic))

	var prot protocol.Protocol
	ok := false
	for portName, md5_value := range t.ctx.vined.serialMd5 {
		if protocolMagic == md5_value {
			t.ctx.vined.ApiTCPConn[clientConn] = portName //set up map for client conn and serial port name
			ok = true
			break
		}
	}
	if ok {
		prot = &protocolV1{ctx: t.ctx}
	} else {
		protocol.SendFramedResponse(clientConn, frameTypeError, []byte("E_BAD_PROTOCOL"))
		clientConn.Close()
		log.Error(fmt.Sprintf("client(%s) bad protocol serial md5 '%v'", clientConn.RemoteAddr(), protocolMagic))
		return
	}

	err = prot.IOLoop(clientConn)
	if err != nil {
		log.Error(fmt.Sprintf("client(%s) - %s", clientConn.RemoteAddr(), err))
		return
	}
}
