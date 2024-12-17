package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"time"

	"github.com/hybridgroup/tinygo-tkey/pkg/proto"
)

type state int

const (
	stateStarted state = iota
	stateLoading
	stateSigning
	stateFailed
)

var (
	currentState state = stateStarted
	publicKey    ed25519.PublicKey
	privateKey   ed25519.PrivateKey
)

func main() {
	time.Sleep(3 * time.Second)
	println("going to gen public key")

	privateKey = ed25519.NewKeyFromSeed(CDI())
	publicKey = privateKey.Public().(ed25519.PublicKey)
	println("publicKey:", hex.EncodeToString(publicKey))

	rx := make([]byte, 256)
	tx := make([]byte, 256)
	i := 0

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

			// did we receive a full command?
			if i > int(hdr.Len()) {
				handleCommand(rx, tx)

				// reset, and wait for next command
				i = 0
				break
			}

			// wait for more data
			time.Sleep(10 * time.Millisecond)
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func handleCommand(rx []byte, tx []byte) (err error) {
	// clear tx buffer
	for i := 0; i < len(tx); i++ {
		tx[i] = 0
	}

	switch currentState {
	case stateStarted:
		return handleStartedCommand(rx, tx)
	case stateLoading:
		return handleLoadingCommand(rx, tx)
	case stateSigning:
		return handleSigningCommand(rx, tx)
	case stateFailed:
	}

	return errors.ErrUnsupported
}
