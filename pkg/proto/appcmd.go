package proto

type AppCmd struct {
	code   byte
	name   string
	cmdLen CmdLen
}

func NewAppCmd(code byte, name string, cmdLen CmdLen) AppCmd {
	return AppCmd{code, name, cmdLen}
}

func (c AppCmd) Code() byte {
	return c.code
}

func (c AppCmd) CmdLen() CmdLen {
	return c.cmdLen
}

func (c AppCmd) Endpoint() Endpoint {
	return DestApp
}

func (c AppCmd) String() string {
	return c.name
}
