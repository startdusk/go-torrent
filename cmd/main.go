package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/startdusk/go-torrent/torrent"
)

func main() {
	//parse torrent file
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("open file error")
	}
	defer file.Close()
	tf, err := torrent.ParseFile(bufio.NewReader(file))
	if err != nil {
		fmt.Println("parse file error")
	}
	//connect tracker & find peers
	peerID := [20]byte{'c', 'b', 't', '-', '2', '0', '2', '2', '-', '0', '4', '-', '1', '3', '-', '0', '0', '0', '0', '0'}
	peers, err := torrent.FindPeers(tf, peerID, 6881)
	if err != nil {
		fmt.Println("find peers error")
	}
	if len(peers) == 0 {
		fmt.Println("can not find peers")
	}
	//download from peers & make file
	torrent.Download(tf, peers)
	torrent.MakeFile(tf)
}
