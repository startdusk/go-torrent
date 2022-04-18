package storage

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"time"

	"github.com/startdusk/go-torrent/torrent/torrent"
	"go.etcd.io/bbolt"
)

const (
	boltDbCompleteValue   = "c"
	boltDbIncompleteValue = "i"
)

var completionBucketKey = []byte("completion")

type boltdb struct {
	db *bbolt.DB
}

var _ Storage = (*boltdb)(nil)

func NewBoltDB(dir string) (ret Storage, err error) {
	os.MkdirAll(dir, 0o750)
	p := filepath.Join(dir, ".torrent.bolt.db")
	db, err := bbolt.Open(p, 0o660, &bbolt.Options{
		Timeout: time.Second,
	})
	if err != nil {
		return
	}
	db.NoSync = true
	ret = &boltdb{db}
	return
}

func (me boltdb) Get(pk torrent.MetaInfo) (cn Completion, err error) {
	err = me.db.View(func(tx *bbolt.Tx) error {
		cb := tx.Bucket(completionBucketKey)
		if cb == nil {
			return nil
		}
		ih := cb.Bucket(pk.InfoHash[:])
		if ih == nil {
			return nil
		}
		var key [4]byte
		binary.BigEndian.PutUint32(key[:], uint32(pk.Index))
		cn.Ok = true
		switch string(ih.Get(key[:])) {
		case boltDbCompleteValue:
			cn.Complete = true
		case boltDbIncompleteValue:
			cn.Complete = false
		default:
			cn.Ok = false
		}
		return nil
	})
	return
}

func (me boltdb) Set(pk torrent.MetaInfo, b bool) error {
	if c, err := me.Get(pk); err == nil && c.Ok && c.Complete == b {
		return nil
	}
	return me.db.Update(func(tx *bbolt.Tx) error {
		c, err := tx.CreateBucketIfNotExists(completionBucketKey)
		if err != nil {
			return err
		}
		ih, err := c.CreateBucketIfNotExists(pk.InfoHash[:])
		if err != nil {
			return err
		}
		var key [4]byte
		binary.BigEndian.PutUint32(key[:], uint32(pk.Index))
		return ih.Put(key[:], []byte(func() string {
			if b {
				return boltDbCompleteValue
			} else {
				return boltDbIncompleteValue
			}
		}()))
	})
}

func (me *boltdb) Close() error {
	return me.db.Close()
}
