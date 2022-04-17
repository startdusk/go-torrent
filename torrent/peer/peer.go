package peer

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

type PeerInfo struct {
	IP   net.IP `bencode:"ip"`
	Port uint16 `bencode:"port"`
}

func (p PeerInfo) String() string {
	return net.JoinHostPort(p.IP.String(), strconv.Itoa(int(p.Port)))
}

// Unmarshal parses peer IP addresses and ports from a buffer
func Unmarshal(bytes []byte) ([]PeerInfo, error) {
	const peerSize = 6 // 4 for ip, 2 for port
	if len(bytes)%peerSize != 0 {
		return nil, fmt.Errorf("received malformed peers")
	}

	peersLen := len(bytes) / peerSize
	peers := make([]PeerInfo, peersLen)
	for i := 0; i < peersLen; i++ {
		offset := i * peerSize
		peers[i].IP = net.IP(bytes[offset : offset+4])
		peers[i].Port = binary.BigEndian.Uint16([]byte(bytes[offset+4 : offset+6]))
	}
	return peers, nil
}
