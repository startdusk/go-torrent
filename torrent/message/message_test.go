package message

import (
	"testing"

	"gotest.tools/assert"
)

func TestCreateRequest(t *testing.T) {
	msg := CreateReq(4, 567, 4321)
	expected := &Message{
		ID: MsgRequest,
		Payload: []byte{
			0x00, 0x00, 0x00, 0x04, // index
			0x00, 0x00, 0x02, 0x37, // begin
			0x00, 0x00, 0x10, 0xe1, // length
		},
	}
	assert.DeepEqual(t, msg, expected)
}

func TestCreateHave(t *testing.T) {
	msg := CreateHave(4)
	expected := &Message{
		ID: 0,
		Payload: []byte{
			0x00, 0x00, 0x00, 0x04,
		},
	}
	assert.DeepEqual(t, msg, expected)
}
