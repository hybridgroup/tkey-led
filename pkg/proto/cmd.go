package proto

// Length of command data that follows the first 1 byte frame header
type CmdLen byte

const (
	CmdLen1   CmdLen = 0
	CmdLen4   CmdLen = 1
	CmdLen32  CmdLen = 2
	CmdLen128 CmdLen = 3
)

// Bytelen returns the number of bytes corresponding to the specific
// CmdLen value.
func (l CmdLen) Bytelen() int {
	switch l {
	case CmdLen1:
		return 1
	case CmdLen4:
		return 4
	case CmdLen32:
		return 32
	case CmdLen128:
		return 128
	}
	return 0
}

type Cmd interface {
	Code() byte
	String() string

	CmdLen() CmdLen
	Endpoint() Endpoint
}
