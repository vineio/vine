package serial

import (
	"crypto/md5"
	"io"

	log "github.com/donnie4w/go-logger/logger"
	"github.com/jacobsa/go-serial/serial"
)

type OptionSerial struct {
	PortName        string
	BaudRate        uint
	MinimumReadSize uint
}

func NewSerial(portnum string, baudrate uint) (io.ReadWriteCloser, error) {
	opt := serial.OpenOptions{
		PortName:        portnum,
		BaudRate:        baudrate,
		DataBits:        8,
		StopBits:        1,
		ParityMode:      serial.PARITY_NONE,
		MinimumReadSize: 8,
	}

	var err error
	var tempIO io.ReadWriteCloser

	tempIO, err = serial.Open(opt)

	return tempIO, err
}

func New(opts []OptionSerial) (map[string]io.ReadWriteCloser, map[string][md5.Size]byte, error) {
	// opts := make([]OptionSerial, 0)

	// config.Unmarshal("serial", &opts)
	// if len(opts) <= 0 {
	// return nil, errors.New("no serial class")
	// }

	tempios := make(map[string]io.ReadWriteCloser)
	tempmd5 := make(map[string][md5.Size]byte)

	for _, v := range opts {
		var io io.ReadWriteCloser
		var err error

		if io, err = NewSerial(v.PortName, v.BaudRate); err != nil {
			log.Error("NewSerial:", err)
			continue
		}
		tempios[v.PortName] = io
		tempmd5[v.PortName] = md5.Sum([]byte(v.PortName))
	}

	return tempios, tempmd5, nil
}
