package vined

import (
	"io"
	"net"

	log "github.com/donnie4w/go-logger/logger"
)

type Client struct {
	ctx *context
}

func (c *Client) Handle(ios map[string]io.ReadWriteCloser, address string) {
	for k, _ := range ios {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			log.Error("util.NewTcp:", err)
			continue
		}

		c.ctx.vined.MsgRxChan[k] = make(chan []byte)
		c.ctx.vined.MsgTxChan[k] = make(chan []byte)

		c.ctx.vined.ClientConnChan[conn.LocalAddr().String()] = make(chan string)
		c.ctx.vined.ClientConnChan[conn.LocalAddr().String()] <- k
		log.Debug(k, conn.LocalAddr().String())
		go c.IOLoop(conn, k)
	}

}

//read&wirte with net.TCPConn
func (c *Client) IOLoop(conn net.Conn, portName string) {
	go c.Write(conn, portName)
	c.Read(conn, portName)

}

func (c *Client) Read(conn net.Conn, portName string) {
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
		msg.NodeName = c.ctx.vined.getOpts().BroadcastAddress
		msg.PortName = portName

		c.ctx.vined.MsgRxChan[portName] <- msg.Data
		c.ctx.vined.MessageTotalRxChan <- msg
		log.Debug("get data：", msg)
	}
}

func (c *Client) Write(conn net.Conn, portName string) {
	for {
		data := <-c.ctx.vined.MsgTxChan[portName]
		conn.Write(data)
		msg := Message{
			NodeName: c.ctx.vined.getOpts().BroadcastAddress,
			PortName: portName,
			Data:     data,
		}
		c.ctx.vined.MessageTotalTxChan <- msg
	}
}
