package protocol

import (
	"net"
	"runtime"
	"strings"

	log "github.com/donnie4w/go-logger/logger"
)

type TCPHandler interface {
	Handle(net.Conn)
}

func TCPServer(listener net.Listener, handler TCPHandler) {
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				log.Warn("temporary Accept() failure - %s", err)
				runtime.Gosched()
				continue
			}
			// theres no direct way to detect this error because it is not exposed
			if !strings.Contains(err.Error(), "use of closed network connection") {
				log.Error("listener.Accept() - %s", err)
			}
			break
		}
		go handler.Handle(clientConn)
	}

	log.Info("TCP: closing %s", listener.Addr())
}
