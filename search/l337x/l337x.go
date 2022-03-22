package l337x

import (
	"fmt"
	"regexp"

	"github.com/startdusk/go-torrent/search"
)

// Torrent regex
var (
	magnetRe  = regexp.MustCompile(`(stratum-|)magnet:?xt=urn:(sha1|btih|ed2k|aich|kzhash|md5|tree:tiger):([A-Fa-f0-9]+|[A-Za-z2-7]+)&[A-Za-z0-9!@#$%^&*=+.\-_()]*(announce|[A-Fa-f0-9]{40}|[A-Za-z2-7]+)`)
	seedsRe   = regexp.MustCompile(`<td class="coll-2 seeds">([0-9])+</td>`)
	leechesRe = regexp.MustCompile(`<td class="coll-3 leeches">([0-9])+</td>`)
	
	trRe = regexp.MustCompile(`<tr>(.*)+</tr>`)
	torrentRe = regexp.MustCompile(`<td class="coll-1 name"><a href="/sub/[0-9]*/[0-9]*/" class="icon"><i class="flaticon-[a-zA-Z0-9]*"></i></a><a href="(/torrent/[0-9]*/.*?/)"`)
)

type L337x struct {
}

func New() *L337x {
	return &L337x{}
}

func (l *L337x) Search(term string) ([]*search.Torrent, error) {
	_, err := get1337x(term)
	if err != nil {
		return nil, err
	}

	var torrents []*search.Torrent
	// blocks, err := findTorrents(page)
	// if err != nil {
	// 	return nil, err
	// }
	// for _, t := range blocks {
	// 	seeder, leecher, err := findPeer(t)
	// 	if err != nil {
	// 		continue
	// 	}
	// 	magnet, err := findMagnet(t)
	// 	if err != nil {
	// 		continue
	// 	}
	// 	torrents = append(torrents, &search.Torrent{
	// 		Name:     "",
	// 		Magnet:   magnet,
	// 		Seeders:  uint32(seeder),
	// 		Leechers: uint32(leecher),
	// 	})
	// }
	return torrents, nil
}

func (l *L337x) Name() string {
	return "1337x"
}

func get1337x(term string) ([]byte, error) {
	url := fmt.Sprintf("https://1337x.to/search/%v/1/", term)
	bytes, err := search.Fetch(url)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func findTorrents(content []byte) {
	match := torrentRe.FindAllSubmatch(content, -1)
	for _, m := range match {
		fmt.Println(m)
	}
}
