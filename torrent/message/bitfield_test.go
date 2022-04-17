package message

import (
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestHasPiece(t *testing.T) {
	bf := Bitfield{0b01010100, 0b01010100}
	outputs := []bool{false, true, false, true, false, true, false, false, false, true, false, true, false, true, false, false, false, false, false, false}
	for i := 0; i < len(outputs); i++ {
		assert.Equal(t, outputs[i], bf.HasPiece(i))
	}
}

func TestSetPiece(t *testing.T) {
	cases := []struct {
		input Bitfield
		index int
		outpt Bitfield
	}{
		{
			input: Bitfield{0b01010100, 0b01010100},
			index: 4, //          v (set)
			outpt: Bitfield{0b01011100, 0b01010100},
		},
		{
			input: Bitfield{0b01010100, 0b01010100},
			index: 9, //                   v (noop)
			outpt: Bitfield{0b01010100, 0b01010100},
		},
		{
			input: Bitfield{0b01010100, 0b01010100},
			index: 15, //                        v (set)
			outpt: Bitfield{0b01010100, 0b01010101},
		},
		{
			input: Bitfield{0b01010100, 0b01010100},
			index: 19, //                            v (noop)
			outpt: Bitfield{0b01010100, 0b01010100},
		},
	}
	for _, cc := range cases {
		bf := cc.input
		bf.SetPiece(cc.index)
		assert.Equal(t, cc.outpt, bf)
	}
}

func TestRecvBitfield(t *testing.T) {
	tests := map[string]struct {
		msg    []byte
		output Bitfield
		fails  bool
	}{
		"successful bitfield": {
			msg:    []byte{0x00, 0x00, 0x00, 0x06, 5, 1, 2, 3, 4, 5},
			output: Bitfield{1, 2, 3, 4, 5},
			fails:  false,
		},
		"message is not a bitfield": {
			msg:    []byte{0x00, 0x00, 0x00, 0x06, 99, 1, 2, 3, 4, 5},
			output: nil,
			fails:  true,
		},
		"message is keep-alive": {
			msg:    []byte{0x00, 0x00, 0x00, 0x00},
			output: nil,
			fails:  true,
		},
	}

	for _, test := range tests {
		clientConn, serverConn := createClientAndServer(t)
		serverConn.Write(test.msg)

		bf, err := Receive(clientConn, 3*time.Second)

		if test.fails {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, bf, test.output)
		}
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
