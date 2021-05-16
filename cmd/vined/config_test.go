package main

import (
	"testing"
)

type Serial struct {
	PortName        string
	BaudRate        int
	MinimumReadSize int
}

func Test_init(t *testing.T) {
	t.Log("test config")

	ss := make([]Serial, 0)
	Unmarshal("serial", &ss)

	t.Log(ss)

}
