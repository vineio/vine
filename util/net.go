package util

import (
	"io"
	"sync"

	log "github.com/donnie4w/go-logger/logger"
)

// will block until connection close
func Join(c1 io.ReadWriteCloser, c2 io.ReadWriteCloser) {
	var wait sync.WaitGroup
	pipe := func(to io.ReadWriteCloser, from io.ReadWriteCloser) {
		defer to.Close()
		defer from.Close()
		defer wait.Done()

		var err error
		_, err = io.Copy(to, from)
		if err != nil {
			log.Warn("join conns error, %v", err)
		}
	}

	wait.Add(2)
	go pipe(c1, c2)
	go pipe(c2, c1)
	wait.Wait()
	return
}
