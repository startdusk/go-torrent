package l337x

import (
	_ "embed"
	"testing"
)

//go:embed l337x_index_test.html
var indexFile []byte

//go:embed l337x_detail_test.html
var detailFile []byte

func TestParse(t *testing.T) {
	findTorrents(indexFile)

	match := leechesRe.FindAllSubmatch(indexFile, -1)
	t.Errorf("%v", match)
	for _, m := range match {
		t.Errorf("%v", m)
	}
}
