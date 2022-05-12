package torrent

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/startdusk/go-torrent/torrent/peer"
	"github.com/startdusk/go-torrent/torrent/types"
	"golang.org/x/sync/errgroup"
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
	if !tf.Info.IsMultiple() {
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

	var g errgroup.Group
	for _, file := range tf.MultipleFile.Files {
		file := file
		g.Go(func() error {
			p := path.Join(targetDir, tf.FileName, tf.MultipleFile.Name, strings.Join(file.Path, "/"))
			log.Printf("assemble file: %s\n", p)
			f, err := os.Create(p)
			if err != nil {
				return err
			}
			defer f.Close()
			w := bufio.NewWriter(f)

			for i := 0; i < len(tf.PieceHashes); i++ {
				err := func() error {
					// p := path.Join(sourceDir, )
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

			log.Printf("completed assemble file: %s !!!\n", p)
			return w.Flush()
		})
	}

	return g.Wait()
}
