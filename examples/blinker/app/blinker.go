package main

import (
	"encoding/binary"
	"time"

	"github.com/hybridgroup/tinygo-tkey/pkg/proto"
)

var (
	blinking = true
	timing   = 500

	cmdSetLED    = proto.NewAppCmd(0x01, "cmdSetLED", proto.CmdLen32)
	rspSetLED    = proto.NewAppCmd(0x02, "rspSetLED", proto.CmdLen4)
	cmdSetTiming = proto.NewAppCmd(0x03, "cmdSetTiming", proto.CmdLen32)
	rspSetTiming = proto.NewAppCmd(0x04, "rspSetTiming", proto.CmdLen4)
	cmdBlinking  = proto.NewAppCmd(0x05, "cmdBlinking", proto.CmdLen4)
	rspBlinking  = proto.NewAppCmd(0x06, "rspBlinking", proto.CmdLen4)
)

func main() {
	rx := make([]byte, 256)
	tx := make([]byte, 256)
	i := 0
	on := true

	for {
		for uart.Buffered() > 0 {
			data, _ := uart.ReadByte()
			rx[i] = data
			i++

			hdr, err := proto.ParseFramingHdr(rx[0])
			if err != nil {
				// reset, and wait for next command
				i = 0
				break
			}

			// did we receive full command?
			if i > int(hdr.Len()) {
				handleCommand(rx, tx)

				// reset, and wait for next command
				i = 0
				break
			}

			// wait for more data
			time.Sleep(time.Millisecond * 10)
		}

		if blinking {
			ledSet(on)
			on = !on
		}

		time.Sleep(time.Millisecond * time.Duration(timing))
	}
}

func handleCommand(rx []byte, tx []byte) (err error) {
	var response proto.Frame

	switch rx[1] {
	case cmdSetLED.Code():
		changeLED(rx[2])

		response, err = proto.NewFrame(rspSetLED, 2, []byte{proto.StatusOK})

	case cmdSetTiming.Code():
		timing = int(binary.LittleEndian.Uint16(rx[2:]))

		response, err = proto.NewFrame(rspSetTiming, 2, []byte{proto.StatusOK})

	case cmdBlinking.Code():
		blinking = rx[2] == 1

		response, err = proto.NewFrame(rspBlinking, 2, []byte{proto.StatusOK})

	default:
		response, err = proto.NewFrame(proto.NewAppCmd(0x00, "cmdUnknown", proto.CmdLen1), 2, []byte{proto.StatusBad})

	}

	if err != nil {
		return err
	}

	// read response into tx buffer
	response.Read(tx)

	// write tx buffer with response
	uart.Write(tx[:response.Len()+1])

	return nil
}
