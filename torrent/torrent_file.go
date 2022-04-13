package torrent

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/startdusk/go-torrent/bencode"
)

type rawInfo struct {
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
}

type rawFile struct {
	Announce string  `bencode:"announce"`
	Info     rawInfo `bencode:"info"`
}

const SHALEN int = 20

type TorrentFile struct {
	Announce string
	InfoSHA  [SHALEN]byte
	FileName string
	FileLen  int
	PieceLen int
	PieceSHA [][SHALEN]byte
}

func ParseFile(r io.Reader) (*TorrentFile, error) {
	raw := new(rawFile)
	err := bencode.Unmarshal(r, raw)
	if err != nil {
		return nil, fmt.Errorf("Fail to parse torrent file: %w", err)
	}
	ret := new(TorrentFile)
	ret.Announce = raw.Announce
	ret.FileName = raw.Info.Name
	ret.FileLen = raw.Info.Length
	ret.PieceLen = raw.Info.PieceLength
	// calculate info SHA
	buf := new(bytes.Buffer)
	wlen := bencode.Marshal(buf, raw.Info)
	if wlen == 0 {
		return nil, fmt.Errorf("raw file info error, content len = %d", wlen)
	}
	ret.InfoSHA = sha1.Sum(buf.Bytes())

	// calculate pieces SHA
	bys := []byte(raw.Info.Pieces)
	cnt := len(bys) / SHALEN
	hashes := make([][SHALEN]byte, cnt)
	for i := 0; i < cnt; i++ {
		copy(hashes[i][:], bys[i*SHALEN:(i+1)*SHALEN])
	}
	ret.PieceSHA = hashes
	return ret, nil
}

func (t *TorrentFile) buildTrackerURL(peerID [20]byte, port uint16) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}
	params := url.Values{
		"info_hash":  []string{string(t.InfoSHA[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.FileLen)},
	}
	base.RawQuery = params.Encode()
	return base.String(), nil
}
