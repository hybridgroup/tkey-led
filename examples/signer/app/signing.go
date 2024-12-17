package main

import "github.com/hybridgroup/tinygo-tkey/pkg/proto"

func handleSigningCommand(rx []byte, tx []byte) (err error) {
	var response proto.Frame

	switch rx[1] {
	case cmdGetSig.Code():
		response, err = proto.NewFrame(rspGetSig, 2, []byte{proto.StatusOK})

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
