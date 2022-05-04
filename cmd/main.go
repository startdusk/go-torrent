package main

import (
	"bufio"
	"log"
	"os"

	"github.com/startdusk/go-torrent/torrent"

	"github.com/startdusk/go-torrent/torrent/types"
)

var peerID types.PeerID = [20]byte{}
var port uint16 = 6881

func main() {
	if len(os.Args[1:]) != 2 {
		log.Fatal("please input the torrent file path and output file path eg: go-torrent ./demo.torrent ./")
	}
	// parse torrent file
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("open file error %+v", err)
	}
	defer file.Close()
	tf, err := torrent.ParseFile(bufio.NewReader(file))
	if err != nil {
		log.Fatalf("parse file error %+v", err)
	}
	// connect tracker & find peers
	peers, err := torrent.FindPeers(tf, peerID, port)
	if err != nil {
		log.Fatalf("find peers error %+v", err)
	}
	if len(peers) == 0 {
		log.Fatalf("can not find peers")
	}
	// download from peers & make file
	tempDir, err := torrent.Download(tf, peerID, peers)
	if err != nil {
		log.Fatal(err)
	}
	err = torrent.MakeFile(tf, tempDir, os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
}
