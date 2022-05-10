package torrent

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindPeers(t *testing.T) {
	file, err := os.Open("../testfile/debian.iso.torrent")
	assert.Equal(t, nil, err)
	tf, err := ParseFile(bufio.NewReader(file))
	assert.Equal(t, nil, err)
	peerID := [20]byte{}
	peers, err := FindPeers(tf, peerID, 6881)
	assert.Equal(t, nil, err)
	assert.Equal(t, len(peers), 50)
}
