package client

import (
	"net"
	"time"

	"github.com/startdusk/go-torrent/torrent/handshake"
	"github.com/startdusk/go-torrent/torrent/message"
	"github.com/startdusk/go-torrent/torrent/peer"
	"github.com/startdusk/go-torrent/torrent/torrent"
)

const Timeout = 3 * time.Second

// Client a TCP connection with a peer
type Client struct {
	Conn     net.Conn
	Choked   bool
	Bitfield message.Bitfield
	peer     peer.PeerInfo
	infoHash torrent.InfoHash
	peerID   torrent.PeerID
}

// New connects with a peer, completes a handshake, and receives a handshake
// returns an err if any of those fail.
func New(peer peer.PeerInfo, peerID torrent.PeerID, infoHash torrent.InfoHash) (*Client, error) {
	// 1.create a tcp connection
	conn, err := net.DialTimeout("tcp", peer.String(), Timeout)
	if err != nil {
		return nil, err
	}

	// 2.complete handshake
	_, err = handshake.Connect(conn, peerID, infoHash, Timeout)
	if err != nil {
		conn.Close()
		return nil, err
	}

	// 3.receive bitfield
	bitfield, err := message.Receive(conn, Timeout)
	if err != nil {
		conn.Close()
		return nil, err
	}
	return &Client{
		Conn:     conn,
		Choked:   true,
		Bitfield: bitfield,
		peer:     peer,
		infoHash: infoHash,
		peerID:   peerID,
	}, nil
}

// Read reads and consumes a message from the connection
func (c *Client) Read() (*message.Message, error) {
	return message.Read(c.Conn)

}

// SendRequest sends a Request message to the peer
func (c *Client) SendRequest(index, begin, length int) error {
	msg := message.CreateReq(index, begin, length)
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

func (c *Client) SendChoke() error {
	msg := message.Message{ID: message.MsgChoke}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

func (c *Client) SendUnChoke() error {
	msg := message.Message{ID: message.MsgUnchoke}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

func (c *Client) SendInterested() error {
	msg := message.Message{ID: message.MsgInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

func (c *Client) SendNotInterested() error {
	msg := message.Message{ID: message.MsgNotInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

func (c *Client) SendHave(index int) error {
	msg := message.CreateHave(index)
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

func (c *Client) SentCancel() error {
	msg := message.Message{ID: message.MsgCancel}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}
