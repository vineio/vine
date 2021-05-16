package main

import (
	"github.com/vineio/vine/serial"
	"github.com/vineio/vine/vined"

	log "github.com/donnie4w/go-logger/logger"
)

func main() {
	optserials := make([]serial.OptionSerial, 0)

	Unmarshal("serial", &optserials)
	if len(optserials) <= 0 {
		log.Error("no serial class")
		return
	}

	opts := vined.NewOptions()
	opts.Optserials = optserials

	v := vined.New(opts)

	v.Main()
}
