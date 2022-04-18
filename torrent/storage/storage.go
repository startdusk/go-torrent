package storage

import (
	"github.com/startdusk/go-torrent/torrent/torrent"
)

type MetaData struct {
	InfoHash torrent.InfoHash
	Index    int
	Begin    int
	End      int
}

type Storage interface {
	Get(me MetaData) (Completion, error)
	Set(me MetaData, complete bool) error
	Close() error
}

type Completion struct {
	Complete bool
	Ok       bool
}
