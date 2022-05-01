package storage

import (
	"github.com/startdusk/go-torrent/torrent/types"
)

type MetaData struct {
	InfoHash types.InfoHash
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
