package vined

import (
	"crypto/md5"
	"hash/crc32"
	"io"
	"log"
	"os"

	"github.com/vineio/vine/serial"
)

type Options struct {
	ID int64 `flag:"node-id" cfg:"id"`

	TCPAddress     string `flag:"tcp-address"`
	ApiTCPAddress  string `flag:"api-tcp-address"`
	ApiHTTPAddress string `flag:"api-http-address"`

	Optserials []serial.OptionSerial `flag:"option-serial"`

	DataPath         string `flag:"data-path"`
	BroadcastAddress string `flag:"broadcast-address"`
}

func NewOptions() *Options {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	h := md5.New()
	io.WriteString(h, hostname)
	defaultID := int64(crc32.ChecksumIEEE(h.Sum(nil)) % 1024)
	return &Options{
		ID:         defaultID,
		TCPAddress: "0.0.0.0:4201",

		BroadcastAddress: hostname,
	}
}
