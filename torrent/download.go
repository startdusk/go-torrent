package torrent

import (
	"github.com/startdusk/go-torrent/torrent/peer"
)

func Download(tf *TorrentFile, peers []peer.PeerInfo) error {
	//TODO: check local tmp file
	// tmpdir, err := os.MkdirTemp("", "torrent-tmp-")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//TODO: download piceces and check
	//TODO: write picece bytes into local tmp file
	return nil
}

func MakeFile(tf *TorrentFile) {
	//TODO: assemble tmp to file

	// 使用内嵌的kv数据库(boltdb)存储临时文件的分片信息
	// https://github1s.com/anacrolix/torrent/blob/master/storage/bolt-piece-completion.go

}
