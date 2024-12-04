package proto

import (
	"errors"
)

type Endpoint byte

const (
	// destAFPGA endpoint = 1
	DestFW  Endpoint = 2
	DestApp Endpoint = 3
)

const (
	StatusOK byte = iota
	StatusBad
)

type FramingHdr struct {
	ID            byte
	Endpoint      Endpoint
	CmdLen        CmdLen
	ResponseNotOK bool
}

func Parse(b byte) (FramingHdr, error) {
	var f FramingHdr

	if (b & 0b1000_0000) != 0 {
		return f, errors.New("reserved bit #7 is not zero")
	}

	// If bit #2 is set
	if (b & 0b0000_0100) != 0 {
		f.ResponseNotOK = true
	}

	f.ID = byte((b & 0b0110_0000) >> 5)
	f.Endpoint = Endpoint((b & 0b0001_1000) >> 3)
	f.CmdLen = CmdLen(b & 0b0000_0011)

	return f, nil
}

// NewFrame writes the framing protocol header byte into a slice of bytes.
// The cmd parameter is used to get the endpoint and command length, which
// together with id parameter are encoded as the header byte. The
// header byte is placed in the first byte in the returned buffer. The
// command code from cmd is placed in the buffer's second byte.
//
// Header byte (used for both command and response frame):
//
// Bit [7] (1 bit). Reserved - possible protocol version.
//
// Bits [6..5] (2 bits). Frame ID tag.
//
// Bits [4..3] (2 bits). Endpoint number:
//
//	00 == reserved
//	01 == HW in application_fpga
//	10 == FW in application_fpga
//	11 == SW (application) in application_fpga
//
// Bit [2] (1 bit). Usage:
//
//	Command: Unused. MUST be zero.
//	Response: 0 == OK, 1 == Not OK (NOK)
//
// Bits [1..0] (2 bits). Command/Response data length:
//
//	00 == 1 byte
//	01 == 4 bytes
//	10 == 32 bytes
//	11 == 128 bytes
//
// Note that the number of bytes indicated by the command data length
// field does **not** include the header byte. This means that a
// complete command frame, with a header indicating a command length
// of 128 bytes, is 128+1 bytes in length.
func NewFrame(cmd Cmd, id int, buf []byte) error {
	if id > 3 {
		return errors.New("frame ID must be 0..3")
	}
	if cmd.Endpoint() > 3 {
		return errors.New("endpoint must be 0..3")
	}
	if cmd.CmdLen() > 3 {
		return errors.New("cmdlen must be 0..3")
	}
	if len(buf) < 1+cmd.CmdLen().Bytelen() {
		return errors.New("buffer too small for frame")
	}

	buf[0] = (byte(id) << 5) | (byte(cmd.Endpoint()) << 3) | byte(cmd.CmdLen())

	// Set command code
	buf[1] = cmd.Code()

	return nil
}

// SetFrameData sets the data in a frame buffer. The data is copied
// into the buffer starting at the third byte.
func SetFrameData(buf []byte, data []byte) error {
	copy(buf[2:], data)

	return nil
}
