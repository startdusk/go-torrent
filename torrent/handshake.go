package torrent

import (
	"fmt"
	"io"
)

const PROTOCOL = "BitTorrent protocol"

type Handshake struct {
	Pstrlen  int
	Pstr     string
	Reserved [8]byte
	InfoHash [SHALEN]byte
	PeerID   [SHALEN]byte
}

func NewHandshake(infoHash, peerID [SHALEN]byte) *Handshake {
	return &Handshake{
		Pstrlen:  len(PROTOCOL),
		Pstr:     PROTOCOL,
		Reserved: [8]byte{0, 0, 0, 0, 0, 0, 0, 0},
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

func ReadHandshake(r io.Reader) (*Handshake, error) {
	lenBuf := make([]byte, 1)
	_, err := io.ReadFull(r, lenBuf)
	if err != nil {
		return nil, err
	}
	pstrlen := int(lenBuf[0])
	if pstrlen == 0 {
		return nil, fmt.Errorf("pstrlen cannot be 0")
	}

	contentBuf := make([]byte, 49+pstrlen-1)
	_, err = io.ReadFull(r, contentBuf)
	if err != nil {
		return nil, err
	}

	var infoHash, peerID [SHALEN]byte
	copy(infoHash[:], contentBuf[pstrlen+8:pstrlen+8+20])
	copy(peerID[:], contentBuf[pstrlen+8+20:])
	return &Handshake{
		Pstrlen:  pstrlen,
		Pstr:     string(contentBuf[0:pstrlen]),
		Reserved: [8]byte{0, 0, 0, 0, 0, 0, 0, 0},
		InfoHash: infoHash,
		PeerID:   peerID,
	}, nil
}
