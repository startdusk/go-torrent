package client

import (
	"net"
	"sync"
	"testing"

	"github.com/startdusk/go-torrent/torrent/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	clientConn, serverConn := createClientAndServer(t)
	client := Client{Conn: clientConn}

	msgBytes := []byte{
		0x00, 0x00, 0x00, 0x05,
		4,
		0x00, 0x00, 0x05, 0x3c,
	}
	expected := &message.Message{
		ID:      message.MsgHave,
		Payload: []byte{0x00, 0x00, 0x05, 0x3c},
	}
	_, err := serverConn.Write(msgBytes)
	require.Nil(t, err)

	msg, err := client.Read()
	assert.Nil(t, err)
	assert.Equal(t, expected, msg)
}

func TestSendRequest(t *testing.T) {
	clientConn, serverConn := createClientAndServer(t)
	client := Client{Conn: clientConn}
	err := client.SendRequest(1, 2, 3)
	assert.Nil(t, err)
	expected := []byte{
		0x00, 0x00, 0x00, 0x0d,
		6,
		0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x02,
		0x00, 0x00, 0x00, 0x03,
	}
	buf := make([]byte, len(expected))
	_, err = serverConn.Read(buf)
	assert.Nil(t, err)
	assert.Equal(t, expected, buf)
}

func TestSendNotInterested(t *testing.T) {
	clientConn, serverConn := createClientAndServer(t)
	client := Client{Conn: clientConn}
	err := client.SendNotInterested()
	assert.Nil(t, err)
	expected := []byte{
		0x00, 0x00, 0x00, 0x01,
		3,
	}
	buf := make([]byte, len(expected))
	_, err = serverConn.Read(buf)
	assert.Nil(t, err)
	assert.Equal(t, expected, buf)
}

func TestSendUnchoke(t *testing.T) {
	clientConn, serverConn := createClientAndServer(t)
	client := Client{Conn: clientConn}
	err := client.SendUnchoke()
	assert.Nil(t, err)
	expected := []byte{
		0x00, 0x00, 0x00, 0x01,
		1,
	}
	buf := make([]byte, len(expected))
	_, err = serverConn.Read(buf)
	assert.Nil(t, err)
	assert.Equal(t, expected, buf)
}

func TestSendHave(t *testing.T) {
	clientConn, serverConn := createClientAndServer(t)
	client := Client{Conn: clientConn}
	err := client.SendHave(1340)
	assert.Nil(t, err)
	expected := []byte{
		0x00, 0x00, 0x00, 0x05,
		4,
		0x00, 0x00, 0x05, 0x3c,
	}
	buf := make([]byte, len(expected))
	_, err = serverConn.Read(buf)
	assert.Nil(t, err)
	assert.Equal(t, expected, buf)
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
