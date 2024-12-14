package main

import (
	"testing"

	"github.com/hybridgroup/tinygo-tkey/pkg/proto"
)

func TestHandleCommand(t *testing.T) {
	tests := []struct {
		name     string
		cmd      proto.AppCmd
		id       int
		data     []byte
		hasError bool
	}{
		{
			name:     "cmdSetLED",
			cmd:      cmdSetLED,
			id:       2,
			data:     []byte{0x00, 0x00, 0x00, 0x00},
			hasError: false,
		},
		{
			name:     "cmdSetTiming",
			cmd:      cmdSetTiming,
			id:       2,
			data:     []byte{0x00, 0x00, 0x00, 0x00},
			hasError: false,
		},
		{
			name:     "cmdBlinking",
			cmd:      cmdBlinking,
			id:       2,
			data:     []byte{0x00, 0x00, 0x00, 0x00},
			hasError: false,
		},
	}

	for _, tt := range tests {
		rx := make([]byte, 256)
		tx := make([]byte, 256)

		t.Run(tt.name, func(t *testing.T) {
			frame, _ := proto.NewFrame(tt.cmd, tt.id, tt.data)
			frame.Read(rx)

			err := handleCommand(rx, tx)
			if err != nil && !tt.hasError {
				t.Errorf("handleCommand(%v) = %v, want %v", tt.cmd, err, tt.hasError)
			}
		})
	}
}
