package proto

import (
	"testing"
)

func TestParseFramingHdr(t *testing.T) {
	tests := []struct {
		name    string
		b       byte
		want    FramingHdr
		wantErr bool
	}{
		{
			name:    "bad reserved bit",
			b:       0b1000_0000,
			want:    FramingHdr{},
			wantErr: true,
		},
		{
			name: "response not OK",
			b:    0b0000_0100,
			want: FramingHdr{
				ResponseNotOK: true,
			},
		},
		{
			name: "ID",
			b:    0b0110_0000,
			want: FramingHdr{
				ID: 0b011,
			},
		},
		{
			name: "endpoint",
			b:    0b0001_1000,
			want: FramingHdr{
				Endpoint: DestApp,
			},
		},
		{
			name: "cmdlen",
			b:    0b0000_0011,
			want: FramingHdr{
				CmdLen: CmdLen128,
			},
		},
		{
			name: "all",
			b:    0b0111_1011,
			want: FramingHdr{
				ID:       0b011,
				Endpoint: DestApp,
				CmdLen:   CmdLen128,
			},
		},
		{
			name: "all response not OK",
			b:    0b0111_1011 | 0b0000_0100,
			want: FramingHdr{
				ID:            0b011,
				Endpoint:      DestApp,
				CmdLen:        CmdLen128,
				ResponseNotOK: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFramingHdr(tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFramingHdr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseFramingHdr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFrame_Read(t *testing.T) {
	tests := []struct {
		name    string
		f       Frame
		s       []byte
		want    int
		wantErr bool
	}{
		{
			name: "bad frame ID",
			f: Frame{
				cmd:  CmdGetNameVersion,
				id:   4,
				data: []byte{0x00},
			},
			s:       make([]byte, 4),
			want:    0,
			wantErr: true,
		},
		{
			name: "bad buffer len",
			f: Frame{
				cmd:  CmdLoadApp,
				id:   0,
				data: []byte{0x00},
			},
			s:       make([]byte, 4),
			want:    0,
			wantErr: true,
		},
		{
			name: "ok",
			f: Frame{
				cmd:  CmdGetNameVersion,
				id:   0,
				data: []byte{0x00},
			},
			s:    make([]byte, 4),
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.Read(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Frame.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Frame.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}
