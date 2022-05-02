package torrent

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"

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
	fmt.Println(tempDir)
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
	// assemble tmp to file
	// TODO: support other platform
	cmd := fmt.Sprintf("cd %s && ls | sort -n | xargs cat > %s", sourceDir, targetDir)
	_, err := execLinuxShell(cmd)
	return err
}

func execLinuxShell(s string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", s)

	var result bytes.Buffer
	cmd.Stdout = &result

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return result.String(), err
}
