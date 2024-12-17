package main

import (
	"github.com/hybridgroup/tinygo-tkey/pkg/proto"
)

const (
	cmdNone = iota
	cmdIdGetPublicKey
	rspIdGetPublicKey
	cmdIDSetSize
	rspIDSetSize
	cmdIDLoadData
	rspIDLoadData
	cmdIDGetSig
	rspIDGetSig
	cmdIDGetNameVersion
	rspIDGetNameVersion
	cmdIDGetFirmwareHash
	rspIDGetFirmwareHash
)

var (
	cmdGetPublicKey    = proto.NewAppCmd(cmdIdGetPublicKey, "cmdGetPublicKey", proto.CmdLen1)
	rspGetPublicKey    = proto.NewAppCmd(rspIdGetPublicKey, "rspGetPublicKey", proto.CmdLen128)
	cmdSetSize         = proto.NewAppCmd(cmdIDSetSize, "cmdSetSize", proto.CmdLen32)
	rspSetSize         = proto.NewAppCmd(rspIDSetSize, "rspSetSize", proto.CmdLen4)
	cmdLoadData        = proto.NewAppCmd(cmdIDLoadData, "cmdLoadData", proto.CmdLen128)
	rspLoadData        = proto.NewAppCmd(rspIDLoadData, "rspLoadData", proto.CmdLen4)
	cmdGetSig          = proto.NewAppCmd(cmdIDGetSig, "cmdGetSig", proto.CmdLen1)
	rspGetSig          = proto.NewAppCmd(rspIDGetSig, "rspGetSig", proto.CmdLen128)
	cmdGetNameVersion  = proto.NewAppCmd(cmdIDGetNameVersion, "cmdGetNameVersion", proto.CmdLen1)
	rspGetNameVersion  = proto.NewAppCmd(rspIDGetNameVersion, "rspGetNameVersion", proto.CmdLen32)
	cmdGetFirmwareHash = proto.NewAppCmd(cmdIDGetFirmwareHash, "cmdGetFirmwareHash", proto.CmdLen32)
	rspGetFirmwareHash = proto.NewAppCmd(rspIDGetFirmwareHash, "rspGetFirmwareHash", proto.CmdLen128)
)
