//go:build tinygo && !tkey

package main

import "machine"

var (
	uart = machine.Serial
)

func init() {
	uart.Configure(machine.UARTConfig{})
}

func CDI() []byte {
	cdi := make([]byte, 32)
	for i := byte(0); i < 32; i++ {
		cdi[i] = i
	}
	return cdi
}
