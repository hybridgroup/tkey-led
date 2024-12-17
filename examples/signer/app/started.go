package main

import (
	"encoding/binary"

	"github.com/hybridgroup/tinygo-tkey/pkg/proto"
)

func handleStartedCommand(rx []byte, tx []byte) (err error) {
	var response proto.Frame

	switch rx[1] {
	case cmdGetPublicKey.Code():
		response, err = proto.NewFrame(rspGetPublicKey, 0, publicKey)

	case cmdSetSize.Code():
		response, err = proto.NewFrame(rspSetSize, 0, []byte{proto.StatusOK})

	case cmdGetNameVersion.Code():
		result := make([]byte, 32)
		copy(result[0:], []byte(app_name0))
		copy(result[4:], []byte(app_name1))
		binary.LittleEndian.PutUint32(result[8:], app_version)

		response, err = proto.NewFrame(rspGetNameVersion, 0, result)

	case cmdGetFirmwareHash.Code():
		response, err = proto.NewFrame(rspGetFirmwareHash, 0, []byte{proto.StatusOK})

	default:
		response, err = proto.NewFrame(proto.NewAppCmd(cmdNone, "cmdUnknown", proto.CmdLen1), 0, []byte{proto.StatusBad})

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
