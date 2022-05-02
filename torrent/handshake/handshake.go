package handshake

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/startdusk/go-torrent/torrent/types"
)

const PROTOCOL = "BitTorrent protocol"

var Extestions = [8]byte{0, 0, 0, 0, 0, 0, 0, 0}

type Handshake struct {
	Pstrlen  int
	Pstr     string
	Reserved [8]byte
	InfoHash types.InfoHash
	PeerID   types.PeerID
}

func New(infoHash, peerID types.PeerID) *Handshake {
	return &Handshake{
		Pstrlen:  len(PROTOCOL),
		Pstr:     PROTOCOL,
		Reserved: Extestions,
		InfoHash: infoHash,
		PeerID:   peerID,
	}
}

func (h Handshake) MsgLen() int {
	// handshake msg: <pstrlen><pstr><reserved><info_hash><peer_id>
	return 1 + len(PROTOCOL) + 8 + 20 + 20
}

// handshake: <pstrlen><pstr><reserved><info_hash><peer_id>
func (h *Handshake) Serialize() []byte {
	buf := make([]byte, h.MsgLen())
	buf[0] = byte(h.Pstrlen)
	cur := 1
	cur += copy(buf[cur:], h.Pstr)
	cur += copy(buf[cur:], h.Reserved[:])
	cur += copy(buf[cur:], h.InfoHash[:])
	cur += copy(buf[cur:], h.PeerID[:])
	return buf
}

func Read(r io.Reader) (*Handshake, error) {
	var buf [68]byte
	_, err := io.ReadFull(r, buf[:68])
	if err != nil {
		return nil, err
	}
	pstrlen := int(buf[0])
	if pstrlen == 0 {
		return nil, fmt.Errorf("pstrlen cannot be 0")
	}

	var infoHash types.InfoHash
	var peerID types.PeerID
	copy(infoHash[:], buf[1+pstrlen+8:1+pstrlen+8+20])
	copy(peerID[:], buf[1+pstrlen+8+20:])
	return &Handshake{
		Pstrlen:  pstrlen,
		Pstr:     string(buf[1 : pstrlen+1]),
		Reserved: [8]byte{0, 0, 0, 0, 0, 0, 0, 0},
		InfoHash: infoHash,
		PeerID:   peerID,
	}, nil
}

// Connect net handshake
func Connect(conn net.Conn, peerID types.PeerID, infoHash types.InfoHash, deadline time.Time) (*Handshake, error) {
	conn.SetDeadline(deadline)
	defer conn.SetDeadline(time.Time{}) // Disable the deadline

	_, err := conn.Write([]byte("\x13" + PROTOCOL))
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(Extestions[:])
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(infoHash[:])
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(peerID[:])
	if err != nil {
		return nil, err
	}

	res, err := Read(conn)
	if err != nil {
		return nil, fmt.Errorf("cannot read conn err: %w", err)
	}
	if !bytes.Equal(res.InfoHash[:], infoHash[:]) {
		return nil, fmt.Errorf("expected info hash %x but got %x", res.InfoHash, infoHash)
	}

	return res, nil
}
