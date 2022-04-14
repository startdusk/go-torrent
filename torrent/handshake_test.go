package torrent

import (
	"bytes"
	"testing"

	"gotest.tools/assert"
)

func TestHandshakeSerialize(t *testing.T) {
	cases := []struct {
		name   string
		input  *Handshake
		output []byte
	}{
		{
			name: "serialize message",
			input: &Handshake{
				Pstrlen:  len(PROTOCOL),
				Pstr:     PROTOCOL,
				Reserved: [8]byte{0, 0, 0, 0, 0, 0, 0, 0},
				InfoHash: [20]byte{134, 212, 200, 0, 36, 164, 105, 190, 76, 80, 188, 90, 16, 44, 247, 23, 128, 49, 0, 116},
				PeerID:   [20]byte{'c', 'b', 't', '-', '2', '0', '2', '2', '-', '0', '4', '-', '1', '3', '-', '0', '0', '0', '0', '0'},
			},
			output: []byte{19, 66, 105, 116, 84, 111, 114, 114, 101, 110, 116, 32, 112, 114, 111, 116, 111, 99, 111, 108, 0, 0, 0, 0, 0, 0, 0, 0, 134, 212, 200, 0, 36, 164, 105, 190, 76, 80, 188, 90, 16, 44, 247, 23, 128, 49, 0, 116, 99, 98, 116, 45, 50, 48, 50, 50, 45, 48, 52, 45, 49, 51, 45, 48, 48, 48, 48, 48},
		},
	}

	for _, cc := range cases {
		assert.DeepEqual(t, cc.input.Serialize(), cc.output)
	}
}

func TestHandshakeDeserialize(t *testing.T) {
	cases := []struct {
		name    string
		input   []byte
		output  *Handshake
		wantErr bool
	}{
		{
			name:  "correct to struct",
			input: []byte{19, 66, 105, 116, 84, 111, 114, 114, 101, 110, 116, 32, 112, 114, 111, 116, 111, 99, 111, 108, 0, 0, 0, 0, 0, 0, 0, 0, 134, 212, 200, 0, 36, 164, 105, 190, 76, 80, 188, 90, 16, 44, 247, 23, 128, 49, 0, 116, 99, 98, 116, 45, 50, 48, 50, 50, 45, 48, 52, 45, 49, 51, 45, 48, 48, 48, 48, 48},
			output: &Handshake{
				Pstrlen:  len(PROTOCOL),
				Pstr:     PROTOCOL,
				Reserved: [8]byte{0, 0, 0, 0, 0, 0, 0, 0},
				InfoHash: [20]byte{134, 212, 200, 0, 36, 164, 105, 190, 76, 80, 188, 90, 16, 44, 247, 23, 128, 49, 0, 116},
				PeerID:   [20]byte{'c', 'b', 't', '-', '2', '0', '2', '2', '-', '0', '4', '-', '1', '3', '-', '0', '0', '0', '0', '0'},
			},
			wantErr: false,
		},
		{
			name:    "empty",
			input:   []byte{},
			output:  nil,
			wantErr: true,
		},
		{
			name:    "pstrlen is 0",
			input:   []byte{0, 0, 0},
			output:  nil,
			wantErr: true,
		},
		{
			name:    "invalid message",
			input:   []byte{1, 2, 3},
			output:  nil,
			wantErr: true,
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			reader := bytes.NewBuffer(cc.input)
			h, err := ReadHandshake(reader)
			if cc.wantErr && err == nil {
				t.Errorf("%s test unexpect want err but got %v", cc.name, err)
			}

			if !cc.wantErr && err != nil {
				t.Errorf("%s test unexpect want err but got %v", cc.name, err)
			}

			assert.DeepEqual(t, cc.output, h)
		})
	}
}
