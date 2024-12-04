package main

import (
	"machine"
	"time"

	"github.com/hybridgroup/tkey-led/pkg/proto"
)

var (
	uart = machine.Serial

	blinking = true
	led      = machine.LED_GREEN
	timing   = 500

	cmdSetLED    = proto.NewAppCmd(0x01, "cmdSetLED", proto.CmdLen32)
	rspSetLED    = proto.NewAppCmd(0x02, "rspSetLED", proto.CmdLen4)
	cmdSetTiming = proto.NewAppCmd(0x03, "cmdSetTiming", proto.CmdLen32)
	rspSetTiming = proto.NewAppCmd(0x04, "rspSetTiming", proto.CmdLen4)
	cmdBlinking  = proto.NewAppCmd(0x05, "cmdBlinking", proto.CmdLen4)
	rspBlinking  = proto.NewAppCmd(0x06, "rspBlinking", proto.CmdLen4)
)

func main() {
	// use default settings for UART
	uart.Configure(machine.UARTConfig{})
	input := make([]byte, 256)
	tx := make([]byte, 256)
	i := 0
	on := true

	for {
		for uart.Buffered() > 0 {
			data, _ := uart.ReadByte()
			input[i] = data
			i++

			hdr, err := proto.Parse(input[0])
			if err != nil {
				// TODO: handle error
			}

			// wait for full command
			if i < int(hdr.CmdLen.Bytelen()) {
				time.Sleep(time.Millisecond * 10)
				continue
			}

			// handle command
			handleCommand(input, i, tx)

			// reset, and wait for next command
			i = 0
			break
		}

		if blinking {
			led.Set(on)
			on = !on
		}

		time.Sleep(time.Millisecond * time.Duration(timing))
	}
}

func handleCommand(input []byte, i int, tx []byte) {
	switch input[1] {
	case cmdSetLED.Code():
		machine.LED_RED.Low()
		machine.LED_GREEN.Low()
		machine.LED_BLUE.Low()

		switch input[2] {
		case 0:
			led = machine.LED_RED
		case 1:
			led = machine.LED_GREEN
		case 2:
			led = machine.LED_BLUE
		}
		proto.NewFrame(rspSetLED, 2, tx)

	case cmdSetTiming.Code():
		timing = int(input[2])
		proto.NewFrame(rspSetTiming, 2, tx)

	case cmdBlinking.Code():
		blinking = input[2] == 1
		proto.NewFrame(rspBlinking, 2, tx)
	}

	proto.SetFrameData(tx, []byte{proto.StatusOK})

	// send response
	uart.Write(tx[:rspSetLED.CmdLen().Bytelen()])
}
