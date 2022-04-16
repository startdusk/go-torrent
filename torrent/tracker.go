package torrent

import (
	"net/http"

	"github.com/startdusk/go-torrent/bencode"
	"github.com/startdusk/go-torrent/torrent/peer"
	"github.com/startdusk/go-torrent/torrent/torrent"
)

func FindPeers(tf *TorrentFile, peerID torrent.PeerID, port uint16) ([]peer.PeerInfo, error) {
	url, err := tf.BuildTrackerURL(peerID, port)
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
	if err != nil {
		return nil, err
	}
	return peer.Unmarshal([]byte(res.Peers))
}
