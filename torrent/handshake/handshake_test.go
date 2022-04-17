package handshake

import (
	"bytes"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		assert.Equal(t, cc.input.Serialize(), cc.output)
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
			h, err := Read(reader)
			if cc.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, cc.output, h)
		})
	}
}

func TestConnect(t *testing.T) {
	tests := map[string]struct {
		clientInfohash  [20]byte
		clientPeerID    [20]byte
		serverHandshake []byte
		output          *Handshake
		fails           bool
	}{
		"successful handshake": {
			clientInfohash:  [20]byte{134, 212, 200, 0, 36, 164, 105, 190, 76, 80, 188, 90, 16, 44, 247, 23, 128, 49, 0, 116},
			clientPeerID:    [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			serverHandshake: []byte{19, 66, 105, 116, 84, 111, 114, 114, 101, 110, 116, 32, 112, 114, 111, 116, 111, 99, 111, 108, 0, 0, 0, 0, 0, 0, 0, 0, 134, 212, 200, 0, 36, 164, 105, 190, 76, 80, 188, 90, 16, 44, 247, 23, 128, 49, 0, 116, 45, 83, 89, 48, 48, 49, 48, 45, 192, 125, 147, 203, 136, 32, 59, 180, 253, 168, 193, 19},
			output: &Handshake{
				Pstrlen:  len(PROTOCOL),
				Pstr:     PROTOCOL,
				InfoHash: [20]byte{134, 212, 200, 0, 36, 164, 105, 190, 76, 80, 188, 90, 16, 44, 247, 23, 128, 49, 0, 116},
				PeerID:   [20]byte{45, 83, 89, 48, 48, 49, 48, 45, 192, 125, 147, 203, 136, 32, 59, 180, 253, 168, 193, 19},
			},
			fails: false,
		},
		"wrong infohash": {
			clientInfohash:  [20]byte{134, 212, 200, 0, 36, 164, 105, 190, 76, 80, 188, 90, 16, 44, 247, 23, 128, 49, 0, 116},
			clientPeerID:    [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			serverHandshake: []byte{19, 66, 105, 116, 84, 111, 114, 114, 101, 110, 116, 32, 112, 114, 111, 116, 111, 99, 111, 108, 0, 0, 0, 0, 0, 0, 0, 0, 0xde, 0xe8, 0x6a, 0x7f, 0xa6, 0xf2, 0x86, 0xa9, 0xd7, 0x4c, 0x36, 0x20, 0x14, 0x61, 0x6a, 0x0f, 0xf5, 0xe4, 0x84, 0x3d, 45, 83, 89, 48, 48, 49, 48, 45, 192, 125, 147, 203, 136, 32, 59, 180, 253, 168, 193, 19},
			output:          nil,
			fails:           true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			clientConn, serverConn := createClientAndServer(t)
			serverConn.Write(test.serverHandshake)

			h, err := Connect(clientConn, test.clientPeerID, test.clientInfohash, 3*time.Second)

			if test.fails {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, h, test.output)
			}
		})
	}
}

func createClientAndServer(t *testing.T) (client, server net.Conn) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.Nil(t, err)

	// net.Dial does not block, so we need this waitgroup to make sure
	// we don't return before serverConn is ready
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer ln.Close()
		server, err = ln.Accept()
		require.Nil(t, err)
		wg.Done()
	}()
	client, err = net.Dial("tcp", ln.Addr().String())
	require.Nil(t, err)
	wg.Wait()
	return
}
