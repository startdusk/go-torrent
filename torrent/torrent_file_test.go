package torrent

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSingleFile(t *testing.T) {
	file, err := os.Open("../testfile/debian.iso.torrent")
	assert.Equal(t, nil, err)
	tf, err := ParseFile(bufio.NewReader(file))
	assert.Equal(t, nil, err)
	assert.Equal(t, "http://bttracker.debian.org:6969/announce", tf.Announce)
	assert.Equal(t, "debian-11.3.0-amd64-netinst.iso", tf.FileName)
	assert.Equal(t, int64(396361728), tf.FileLen)
	assert.Equal(t, 30240, len(tf.SingleFile.Pieces))
	assert.Equal(t, 262144, tf.PieceLen)
	assert.Equal(t, 1512, len(tf.PieceHashes))
	var expectHASH = [20]byte{177, 17, 129, 60, 230, 15, 66, 145, 151, 52, 130, 61, 245, 236, 32, 189, 30, 4, 231, 247}
	assert.Equal(t, expectHASH, tf.InfoHash)
}

func TestParseMultipleFile(t *testing.T) {
	file, err := os.Open("../testfile/MP3-daily-2022-April-02-Electronic-[rarbg.to].torrent")
	assert.Equal(t, nil, err)
	tf, err := ParseFile(bufio.NewReader(file))
	assert.Equal(t, nil, err)

	assert.Equal(t, tf.Announce, "http://tracker.trackerfix.com:80/announce")
	assert.Equal(t, tf.CreationDate, int64(1648983368))
	assert.Equal(t, tf.Comment, "Torrent downloaded from https://rarbg.to")
	assert.Equal(t, tf.CreatedBy, "RARBG")
	assert.Equal(t, tf.AnnounceList, [][]string{
		{"http://tracker.trackerfix.com:80/announce"},
		{"udp://9.rarbg.me:2770/announce"},
		{"udp://9.rarbg.to:2800/announce"},
		{"udp://tracker.fatkhoala.org:13760/announce"},
		{"udp://tracker.thinelephant.org:12750/announce"},
	})

	var expectHASH = [20]byte{0xe9, 0xce, 0x30, 0x88, 0x1b, 0x90, 0x5b, 0x80, 0x9e, 0xc0, 0x94, 0xbc, 0xb2, 0xf1, 0x7d, 0x20, 0x26, 0x27, 0xf0, 0x5f}
	assert.Equal(t, tf.InfoHash, expectHASH)
}
