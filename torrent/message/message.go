package message

import (
	"encoding/binary"
	"fmt"
	"io"
)

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

func (m *Message) Name() string {
	if m == nil {
		return "KeepAlive"
	}
	switch m.ID {
	case MsgChoke:
		return "Choke"
	case MsgUnchoke:
		return "Unchoke"
	case MsgInterested:
		return "Interested"
	case MsgNotInterested:
		return "NotInterested"
	case MsgHave:
		return "Have"
	case MsgBitfield:
		return "Bitfield"
	case MsgRequest:
		return "Request"
	case MsgPiece:
		return "Piece"
	case MsgCancel:
		return "Cancel"
	default:
		return fmt.Sprintf("Unknown MessageID#%d", m.ID)
	}
}

func (m *Message) String() string {
	if m == nil {
		return m.Name()
	}
	return fmt.Sprintf("%s [%d]", m.Name(), len(m.Payload))
}

// Serialize serializes a message into a buffer of the form
// <length prefix><message ID><payload>
// Interprets `nil` as a keep-alive message
func (m *Message) Serialize() []byte {
	if m == nil {
		return make([]byte, 4)
	}
	length := uint32(len(m.Payload) + 1) // +1 for id
	buf := make([]byte, length+4)
	binary.BigEndian.PutUint32(buf[0:4], length)
	buf[4] = byte(m.ID)
	copy(buf[5:], m.Payload)
	return buf
}

// Read parses a message from a stream. Returns `nil` on keep-alive message
func Read(r io.Reader) (*Message, error) {
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lengthBuf)

	// keep-alive message
	if length == 0 {
		return nil, nil
	}

	messageBuf := make([]byte, length)
	_, err = io.ReadFull(r, messageBuf)
	if err != nil {
		return nil, err
	}

	m := Message{
		ID:      MessageID(messageBuf[0]),
		Payload: messageBuf[1:],
	}

	return &m, nil
}

// CreateReq creates a REQUEST message
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

// CreateHave creates a HAVE message
func CreateHave(index int) *Message {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, uint32(index))
	return &Message{
		ID:      MsgHave,
		Payload: payload,
	}
}

// ParsePiece parses a PIECE message and copies its payload into a buffer
func ParsePiece(index int, buf []byte, msg *Message) (int, error) {
	if msg.ID != MsgPiece {
		return 0, fmt.Errorf("expected PIECE ID[%d], but got ID[%d]", MsgPiece, msg.ID)
	}

	if len(msg.Payload) < 8 {
		return 0, fmt.Errorf("payload too short. %d < 8", len(msg.Payload))
	}
	parsedIndex := int(binary.BigEndian.Uint32(msg.Payload[0:4]))
	if parsedIndex != index {
		return 0, fmt.Errorf("expected index %d, got %d", index, parsedIndex)
	}
	begin := int(binary.BigEndian.Uint32(msg.Payload[4:8]))
	if begin >= len(buf) {
		return 0, fmt.Errorf("begin offset too high. %d >= %d", begin, len(buf))
	}
	data := msg.Payload[8:]
	if begin+len(data) > len(buf) {
		return 0, fmt.Errorf("data too long [%d] for offset %d with length %d", len(data), begin, len(buf))
	}
	copy(buf[begin:], data)
	return len(data), nil
}

// ParseHave parses a HAVE message
func ParseHave(msg *Message) (int, error) {
	if msg.ID != MsgHave {
		return 0, fmt.Errorf("expected HAVE ID[%d], but got ID[%d]", MsgHave, msg.ID)
	}

	if len(msg.Payload) != 4 {
		return 0, fmt.Errorf("expected payload length 4, but got length %d", len(msg.Payload))
	}
	index := int(binary.BigEndian.Uint32(msg.Payload))
	return index, nil
}
