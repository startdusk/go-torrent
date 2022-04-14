package torrent

import (
	"fmt"
	"io"
)

const PROTOCOL = "BitTorrent protocol"

type Handshake struct {
	pstrlen  int
	pstr     string
	reserved [8]byte
	infoHash [SHALEN]byte
	peerID   [SHALEN]byte
}

func NewHandshake(infoHash, peerID [SHALEN]byte) *Handshake {
	return &Handshake{
		pstrlen:  len(PROTOCOL),
		pstr:     PROTOCOL,
		reserved: [8]byte{},
		infoHash: infoHash,
		peerID:   peerID,
	}
}

func (h Handshake) MsgLen() int {
	// handshake msg: <pstrlen><pstr><reserved><info_hash><peer_id>
	return 1 + len(PROTOCOL) + 8 + 20 + 20
}

// handshake: <pstrlen><pstr><reserved><info_hash><peer_id>
func (h *Handshake) Serialize() []byte {
	buf := make([]byte, h.MsgLen())
	buf[0] = byte(h.pstrlen)
	cur := 1
	cur += copy(buf[cur:], h.pstr)
	cur += copy(buf[cur:], h.reserved[:])
	cur += copy(buf[cur:], h.infoHash[:])
	cur += copy(buf[cur:], h.peerID[:])
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

	var handshake Handshake
	contentBuf := make([]byte, handshake.MsgLen()-1)
	_, err = io.ReadFull(r, contentBuf)
	if err != nil {
		return nil, err
	}

	var infoHash, peerID [SHALEN]byte
	copy(infoHash[:], contentBuf[pstrlen+8:pstrlen+8+20])
	copy(peerID[:], contentBuf[pstrlen+8+20:])
	handshake.infoHash = infoHash
	handshake.peerID = peerID
	handshake.pstr = string(contentBuf[0:pstrlen])
	return &handshake, nil
}
