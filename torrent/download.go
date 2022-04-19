package torrent

import (
	"os"

	"github.com/startdusk/go-torrent/torrent/peer"
	"github.com/startdusk/go-torrent/torrent/torrent"
)

const TempPrefix = "torrent-temp-"

func Download(tf *TorrentFile, peerID torrent.PeerID, peers []peer.PeerInfo) error {
	// check local tmp file
	tempDir, err := os.MkdirTemp("", TempPrefix+string(tf.InfoHash[:]))
	if err != nil {
		return err
	}
	// download piceces and check
	t := &Torrent{
		Peers:       peers,
		PeerID:      peerID,
		InfoHash:    tf.InfoHash,
		PieceHashes: tf.PieceHashes,
		PieceLen:    tf.PieceLen,
		Length:      tf.FileLen,
		Name:        tf.FileName,
	}
	// write picece bytes into local tmp file
	return t.Download(tempDir)
}

func MakeFile(tf *TorrentFile) {
	//TODO: assemble tmp to file
	// cmd := fmt.Sprintf("cd %s && ls | sort -n | xargs cat > %s", srcPath, destPath)
}
