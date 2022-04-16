package message

import (
	"fmt"
	"net"
	"time"
)

// A Bitfield represents the pieces that a peer has
type Bitfield []byte

// HasPiece tells if a bitfield has a particular index set
func (bf Bitfield) HasPiece(index int) bool {
	byteIndex := index / 8
	offset := index % 8
	if byteIndex < 0 || byteIndex >= len(bf) {
		return false
	}
	return bf[byteIndex]>>uint(7-offset)&1 != 0
}

// SetPiece sets a bit in the bitfield
func (bf Bitfield) SetPiece(index int) {
	byteIndex := index / 8
	offset := index % 8

	// silently discard invalid bounded index
	if byteIndex < 0 || byteIndex >= len(bf) {
		return
	}
	bf[byteIndex] |= 1 << uint(7-offset)
}

func Receive(conn net.Conn, timeout time.Duration) (Bitfield, error) {
	conn.SetDeadline(time.Now().Add(timeout))
	defer conn.SetDeadline(time.Time{})

	msg, err := Read(conn)
	if err != nil {
		return nil, err
	}
	if msg == nil {
		return nil, fmt.Errorf("expected bitfield but got %v", msg)
	}

	if msg.ID != MsgBitfield {
		return nil, fmt.Errorf("expected bitfield but got ID[%d]", msg.ID)
	}
	return msg.Payload, nil
}
