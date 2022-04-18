package storage

import (
	"github.com/startdusk/go-torrent/torrent/torrent"
)

type Storage interface {
	Get(me torrent.MetaInfo) (Completion, error)
	Set(me torrent.MetaInfo, complete bool) error
	Close() error
}

type Completion struct {
	Complete bool
	Ok       bool
}
