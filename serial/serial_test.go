package serial

import (
	"fmt"
	"testing"
)

func Test_New(t *testing.T) {
	opts := make([]OptionSerial, 0)
	opts = append(opts, OptionSerial{"com1", 9600, 4}, OptionSerial{"com3", 9600, 4})
	fmt.Println(opts)
	fmt.Println(New(opts))

}
