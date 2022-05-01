package main

import (
	"bufio"
	"log"
	"os"

	"github.com/startdusk/go-torrent/torrent"

	"github.com/startdusk/go-torrent/torrent/types"
)

var testFile = "../testfile/debian.iso.torrent"

func main() {
	// parse torrent file
	// file, err := os.Open(os.Args[1])
	file, err := os.Open(testFile)
	if err != nil {
		log.Fatalf("open file error %+v", err)
	}
	defer file.Close()
	tf, err := torrent.ParseFile(bufio.NewReader(file))
	if err != nil {
		log.Fatalf("parse file error %+v", err)
	}
	// connect tracker & find peers
	var peerID types.PeerID = [20]byte{}
	peers, err := torrent.FindPeers(tf, peerID, 6881)
	if err != nil {
		log.Fatalf("find peers error %+v", err)
	}
	if len(peers) == 0 {
		log.Fatalf("can not find peers")
	}
	// download from peers & make file
	torrent.Download(tf, peerID, peers)
	torrent.MakeFile(tf)
}
