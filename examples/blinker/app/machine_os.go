//go:build !tinygo

package main

var (
	uart = &UART{}
)

type UART struct {
}

func (u *UART) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (u *UART) ReadByte() (byte, error) {
	return 0, nil
}

func (u *UART) Buffered() int {
	return 0
}

func ledSet(on bool) {
}

func changeLED(p uint8) {
}
