package main

import (
	"bufio"
	"log"
	"os"

	"github.com/startdusk/go-torrent/torrent"
)

func main() {
	// parse torrent file
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Printf("open file error %+v", err)
	}
	defer file.Close()
	tf, err := torrent.ParseFile(bufio.NewReader(file))
	if err != nil {
		log.Printf("parse file error %+v", err)
	}
	// connect tracker & find peers
	peerID := [20]byte{'c', 'b', 't', '-', '2', '0', '2', '2', '-', '0', '4', '-', '1', '3', '-', '0', '0', '0', '0', '0'}
	peers, err := torrent.FindPeers(tf, peerID, 6881)
	if err != nil {
		log.Printf("find peers error %+v", err)
	}
	if len(peers) == 0 {
		log.Printf("can not find peers")
	}
	// download from peers & make file
	torrent.Download(tf, peerID, peers)
	torrent.MakeFile(tf)
}
