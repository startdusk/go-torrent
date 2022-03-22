package search

type Provider interface {
	Search(term string) ([]*Torrent, error)
	Name() string
}
