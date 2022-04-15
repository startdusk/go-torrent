package message

import "encoding/binary"

type MessageID int

const (
	// MsgChoke chokes the receiver
	MsgChoke MessageID = 0
	// MsgUnchoke unchokes the receiver
	MsgUnchoke MessageID = 1
	// MsgInterested expresses interest in receiving data
	MsgInterested MessageID = 2
	// MsgNotInterested expresses disinterest in receiving data
	MsgNotInterested MessageID = 3
	// MsgHave alerts the receiver that the sender has downloaded a piece
	MsgHave MessageID = 4
	// MsgBitfield encodes which pieces that the sender has downloaded
	MsgBitfield MessageID = 5
	// MsgRequest requests a block of data from the receiver
	MsgRequest MessageID = 6
	// MsgPiece delivers a block of data to fulfill a request
	MsgPiece MessageID = 7
	// MsgCancel cancels a request
	MsgCancel MessageID = 8
)

// Message stores ID and payload of a message
type Message struct {
	ID      MessageID
	Payload []byte
}

// creates a request message
func CreateReq(index, begin, length int) *Message {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], uint32(index))
	binary.BigEndian.PutUint32(payload[4:8], uint32(begin))
	binary.BigEndian.PutUint32(payload[8:12], uint32(length))
	return &Message{
		ID:      MsgRequest,
		Payload: payload,
	}
}
