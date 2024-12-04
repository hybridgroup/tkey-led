package main

import (
	"encoding/binary"
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
			if i > int(hdr.CmdLen.Bytelen()) {
				// handle command
				handleCommand(input, i, tx)

				// reset, and wait for next command
				i = 0
				break
			}

			time.Sleep(time.Millisecond * 10)
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

		led = machine.Pin(input[2])

		proto.NewFrame(rspSetLED, 2, tx)
		proto.SetFrameData(tx, []byte{proto.StatusOK})
		uart.Write(tx[:rspSetLED.CmdLen().Bytelen()+1])

	case cmdSetTiming.Code():
		timing = int(binary.LittleEndian.Uint16(input[2:]))

		proto.NewFrame(rspSetTiming, 2, tx)
		proto.SetFrameData(tx, []byte{proto.StatusOK})
		uart.Write(tx[:rspSetTiming.CmdLen().Bytelen()+1])

	case cmdBlinking.Code():
		blinking = input[2] == 1

		proto.NewFrame(rspBlinking, 2, tx)
		proto.SetFrameData(tx, []byte{proto.StatusOK})
		uart.Write(tx[:rspBlinking.CmdLen().Bytelen()+1])
	}

}
