package torrent

import (
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"

	"github.com/startdusk/go-torrent/bencode"
)

type PeerInfo struct {
	IP   net.IP `bencode:"ip"`
	Port uint16 `bencode:"port"`
}

func (p PeerInfo) String() string {
	return net.JoinHostPort(p.IP.String(), strconv.Itoa(int(p.Port)))
}

func FindPeers(tf *TorrentFile, peerID [20]byte, port uint16) ([]PeerInfo, error) {
	url, err := buildTrackerURL(tf, peerID, port)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// TODO: process another json format
	res := new(struct {
		Peers string `bencode:"peers"`
	})
	err = bencode.Unmarshal(resp.Body, res)
	return unmarshal([]byte(res.Peers))
}

// unmarshal parses peer IP addresses and ports from a buffer
func unmarshal(bytes []byte) ([]PeerInfo, error) {
	const peerSize = 6 // 4 for ip, 2 for port
	peersLen := len(bytes) / peerSize
	if len(bytes)%peerSize != 0 {
		return nil, fmt.Errorf("received malformed peers")
	}

	peers := make([]PeerInfo, peersLen)
	for i := 0; i < peersLen; i++ {
		offset := i * peerSize
		peers[i].IP = net.IP(bytes[offset : offset+4])
		peers[i].Port = binary.BigEndian.Uint16([]byte(bytes[offset+4 : offset+6]))
	}
	return peers, nil
}

func buildTrackerURL(tf *TorrentFile, peerID [20]byte, port uint16) (string, error) {
	base, err := url.Parse(tf.Announce)
	if err != nil {
		return "", err
	}
	params := url.Values{
		"info_hash":  []string{string(tf.InfoSHA[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(tf.FileLen)},
	}
	base.RawQuery = params.Encode()
	return base.String(), nil
}
