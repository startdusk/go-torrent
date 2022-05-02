package main

import (
	"bufio"
	"log"
	"os"

	"github.com/startdusk/go-torrent/torrent"

	"github.com/startdusk/go-torrent/torrent/types"
)

func main() {
	if len(os.Args[1:]) != 2 {
		log.Fatal("please input the torrent file path and output file path eg: go-torrent ./demo.torrent ./demo.iso")
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
	var peerID types.PeerID = [20]byte{}
	peers, err := torrent.FindPeers(tf, peerID, 6881)
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
	target := os.Args[2] + tf.FileName
	err = torrent.MakeFile(tf, tempDir, target)
	if err != nil {
		log.Fatal(err)
	}
}
