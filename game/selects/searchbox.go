package selects

import (
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
)

type SearchBoxComponent struct {
	*game.Database

	// sprite draws.Sprite
	game.SearchQuery
	lastSearchResult game.SearchResult
}

func NewSearchBoxComponent(dbs *game.Database) (cmp SearchBoxComponent) {
	cmp.Database = dbs
	cmp.update()
	return cmp
}

func (s *SearchBoxComponent) update() game.SearchResult {
	// TODO: listen to enter
	r := s.Search(s.SearchQuery)
	s.lastSearchResult = r
	return r
}

func (s SearchBoxComponent) Draw(dst draws.Image) {
	// s.sprite.Draw()
}
