package torrent

import (
	"testing"

	"gotest.tools/assert"
)

func TestHasPiece(t *testing.T) {
	bf := Bitfield{0b01010100, 0b01010100}
	outputs := []bool{false, true, false, true, false, true, false, false, false, true, false, true, false, true, false, false, false, false, false, false}
	for i := 0; i < len(outputs); i++ {
		assert.DeepEqual(t, outputs[i], bf.HasPiece(i))
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
		assert.DeepEqual(t, cc.outpt, bf)
	}
}
