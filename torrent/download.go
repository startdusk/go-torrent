package torrent

import "github.com/startdusk/go-torrent/torrent/peer"

func Download(tf *TorrentFile, peers []peer.PeerInfo) error {
	//TODO: check local tmp file
	// tmpdir, err := ioutil.TempDir("", "torrent-pick-")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//TODO: download piceces and check
	//TODO: write picece bytes into local tmp file
	return nil
}

func MakeFile(tf *TorrentFile) {
	//TODO: assemble tmp to file
}
