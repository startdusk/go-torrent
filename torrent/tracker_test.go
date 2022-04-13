package torrent

import (
	"bufio"
	"os"
	"testing"

	"gotest.tools/assert"
)

func TestFindPeers(t *testing.T) {
	file, err := os.Open("../testfile/debian-iso.torrent")
	assert.Equal(t, nil, err)
	tf, err := ParseFile(bufio.NewReader(file))
	assert.Equal(t, nil, err)
	peerID := [20]byte{'c', 'b', 't', '-', '2', '0', '2', '2', '-', '0', '4', '-', '1', '3', '-', '0', '0', '0', '0', '0'}
	peers, err := FindPeers(tf, peerID, 6881)
	assert.Equal(t, nil, err)
	assert.Equal(t, len(peers), 50)
}
