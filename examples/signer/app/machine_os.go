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

func CDI() []byte {
	cdi := make([]byte, 32)
	for i := byte(0); i < 32; i++ {
		cdi[i] = i
	}
	return cdi
}
