package torrent

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/startdusk/go-torrent/bencode"
	"github.com/startdusk/go-torrent/torrent/types"
)

type rawInfo struct {
	Length      int64   `bencode:"length"`
	Name        string  `bencode:"name"`
	PieceLength int     `bencode:"piece length"`
	Pieces      string  `bencode:"pieces"`
	Private     *int    `bencode:"private"`
	Files       []File  `bencode:"files"`
	MD5Sum      *string `bencode:"md5sum"`
}

func (r rawInfo) IsMultiple() bool {
	return len(r.Files) > 0
}

type rawFile struct {
	Announce     string     `bencode:"announce"`
	AnnounceList [][]string `bencode:"announce-list"`
	CreationDate int64      `bencode:"creation date"`
	Comment      string     `bencdoe:"comment"`
	CreatedBy    string     `bencode:"created by"`
	Encoding     string     `bencode:"encoding"`
	Info         rawInfo    `bencode:"info"`
}

func ParseFile(r io.Reader) (*TorrentFile, error) {
	raw := new(rawFile)
	err := bencode.Unmarshal(r, raw)
	if err != nil {
		return nil, fmt.Errorf("fail to parse torrent file: %w", err)
	}
	ret := new(TorrentFile)
	ret.Announce = raw.Announce
	ret.AnnounceList = raw.AnnounceList
	ret.CreationDate = raw.CreationDate
	ret.CreatedBy = raw.CreatedBy
	ret.Comment = raw.Comment
	ret.Encoding = raw.Encoding
	ret.FileName = raw.Info.Name
	ret.FileLen = raw.Info.Length
	ret.PieceLen = raw.Info.PieceLength
	// calculate info SHA
	buf := new(bytes.Buffer)
	var wlen int
	if raw.Info.IsMultiple() {
		ret.Info = Multiple
		ret.MultipleFile = MultipleFile{
			PieceLength: raw.Info.PieceLength,
			Pieces:      raw.Info.Pieces,
			Private:     raw.Info.Private,
			Name:        raw.Info.Name,
			Files:       raw.Info.Files,
		}
		wlen = bencode.Marshal(buf, ret.MultipleFile)
	} else {
		ret.Info = Single
		ret.SingleFile = SingleFile{
			PieceLength: raw.Info.PieceLength,
			Pieces:      raw.Info.Pieces,
			Private:     raw.Info.Private,
			Name:        raw.Info.Name,
			Length:      raw.Info.Length,
			MD5Sum:      raw.Info.MD5Sum,
		}
		wlen = bencode.Marshal(buf, ret.SingleFile)
	}
	if wlen == 0 {
		return nil, fmt.Errorf("raw file info error, content len = %d", wlen)
	}
	ret.InfoHash = sha1.Sum(buf.Bytes())

	// calculate pieces SHA
	bys := []byte(raw.Info.Pieces)
	cnt := len(bys) / types.SHALEN
	hashes := make(types.PieceHashes, cnt)
	for i := 0; i < cnt; i++ {
		copy(hashes[i][:], bys[i*types.SHALEN:(i+1)*types.SHALEN])
	}
	ret.PieceHashes = hashes
	return ret, nil
}

type TorrentFile struct {
	Info         Info
	Announce     string
	AnnounceList [][]string
	InfoHash     types.InfoHash
	CreationDate int64
	Comment      string
	CreatedBy    string
	Encoding     string
	FileName     string
	FileLen      int64
	PieceLen     int
	PieceHashes  types.PieceHashes

	SingleFile   SingleFile
	MultipleFile MultipleFile
}

type File struct {
	Length int64    `bencode:"length"`
	Path   []string `bencode:"path"`
	MD5Sum *string  `bencode:"md5sum"`
}

type Info int

const (
	Single Info = 1 << iota
	Multiple
)

func (i Info) IsMultiple() bool {
	return i == Multiple
}

type SingleFile struct {
	Length      int64   `bencode:"length"`
	Name        string  `bencode:"name"`
	PieceLength int     `bencode:"piece length"`
	Pieces      string  `bencode:"pieces"`
	Private     *int    `bencode:"private"`
	MD5Sum      *string `bencode:"md5sum"`
}

type MultipleFile struct {
	Files       []File `bencode:"files"`
	Name        string `bencode:"name"`
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
	Private     *int   `bencode:"private"`
}

func (tf *TorrentFile) BuildURL(peerID types.PeerID, port uint16) (string, error) {
	base, err := url.Parse(tf.Announce)
	if err != nil {
		return "", err
	}
	params := url.Values{
		"info_hash":  []string{string(tf.InfoHash[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{fmt.Sprintf("%d", tf.FileLen)},
	}
	base.RawQuery = params.Encode()
	return base.String(), nil
}
