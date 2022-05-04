package torrent

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFile(t *testing.T) {

	file, err := os.Open("../testfile/debian.iso.torrent")
	assert.Equal(t, nil, err)
	tf, err := ParseFile(bufio.NewReader(file))
	assert.Equal(t, nil, err)
	assert.Equal(t, "http://bttracker.debian.org:6969/announce", tf.Announce)
	assert.Equal(t, "debian-11.3.0-amd64-netinst.iso", tf.FileName)
	assert.Equal(t, 396361728, tf.FileLen)
	assert.Equal(t, 262144, tf.PieceLen)
	assert.Equal(t, 1512, len(tf.PieceHashes))
	var expectHASH = [20]byte{177, 17, 129, 60, 230, 15, 66, 145, 151, 52, 130, 61, 245, 236, 32, 189, 30, 4, 231, 247}
	assert.Equal(t, expectHASH, tf.InfoHash)
}
