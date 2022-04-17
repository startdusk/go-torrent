package peer

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	cases := []struct {
		name    string
		source  []byte
		target  []PeerInfo
		wantErr bool
	}{
		{
			name:   "parse success",
			source: []byte{127, 0, 0, 1, 0x1a, 0xe1},
			target: []PeerInfo{
				{
					IP:   net.IP{127, 0, 0, 1},
					Port: 6881,
				},
			},
			wantErr: false,
		},
		{
			name:   "parse fail",
			source: []byte{127, 0, 0, 1, 1, 0x1a, 0xe1},
			target: []PeerInfo{
				{
					IP:   net.IP{127, 0, 0, 1},
					Port: 6881,
				},
			},
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			peerInfo, err := Unmarshal(c.source)
			if c.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, peerInfo, c.target)
			}
		})
	}
}
