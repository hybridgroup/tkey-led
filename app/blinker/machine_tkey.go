//go:build tinygo

package main

import "machine"

var (
	uart = machine.Serial

	led = machine.LED_GREEN
)

func init() {
	uart.Configure(machine.UARTConfig{})
}

func allLEDOff() {
	machine.LED_RED.Low()
	machine.LED_GREEN.Low()
	machine.LED_BLUE.Low()
}

func changeLED(p uint8) {
	allLEDOff()

	led = machine.Pin(p)
}

func ledSet(on bool) {
	led.Set(on)
}
