package client

import (
	"net"

	"github.com/startdusk/go-torrent/torrent/message"
	"github.com/startdusk/go-torrent/torrent/peer"
	"github.com/startdusk/go-torrent/torrent/torrent"
)

// Client a TCP connection with a peer
type Client struct {
	Conn     net.Conn
	Choked   bool
	Bitfield message.Bitfield
	peer     peer.PeerInfo
	infoHash torrent.InfoHash
	peerID   torrent.PeerID
}
