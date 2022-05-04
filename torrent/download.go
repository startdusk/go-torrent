package torrent

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/startdusk/go-torrent/torrent/peer"
	"github.com/startdusk/go-torrent/torrent/types"
)

const TempPrefix = "torrent-temp-"

func Download(tf *TorrentFile, peerID types.PeerID, peers []peer.PeerInfo) (string, error) {
	// check local tmp file
	tempDir, err := os.MkdirTemp("", TempPrefix+hex.EncodeToString(tf.InfoHash[:]))
	if err != nil {
		return "", err
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
	return tempDir, t.Download(tempDir)
}

func MakeFile(tf *TorrentFile, sourceDir, targetDir string) error {
	log.Printf("assemble file: %s\n", tf.FileName)
	// assemble tmp to file
	f, err := os.Create(targetDir + "/" + tf.FileName)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	for i := 0; i < len(tf.PieceHashes); i++ {
		err := func() error {
			file, err := os.Open(sourceDir + "/" + fmt.Sprintf("%d", i))
			if err != nil {
				return fmt.Errorf("cannot find the #%d piece: %w", i, err)
			}
			defer file.Close()
			io.Copy(w, file)
			return nil
		}()
		if err != nil {
			return err
		}
	}

	log.Printf("completed assemble file: %s !!!\n", tf.FileName)
	return w.Flush()
}
