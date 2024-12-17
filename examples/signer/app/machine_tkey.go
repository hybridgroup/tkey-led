//go:build tinygo && tkey

package main

import "machine"

var (
	uart = machine.Serial
)

func init() {
	uart.Configure(machine.UARTConfig{})
}

func CDI() []byte {
	return machine.CDI()
}
