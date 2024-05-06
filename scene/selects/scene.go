package selects

import (
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/scene"
)

// Component is basically EventHandler.
type Scene struct {
	*scene.Resources
	*scene.Options
	query          string
	musicList      []scene.MusicRow
	musicListIndex int
	chartList      []scene.ChartRow
	chartListIndex int
}

// It is fine to call Close at blank MusicPlayer.

// Avoid embedding scene.Options directly.
// Pass options as pointers for syncing and saving memory.
func NewScene(res *scene.Resources, opts *scene.Options) (*Scene, error) {
	return &Scene{
		Resources: res,
		Options:   opts,
	}, nil
}

func (sc *Scene) Update() any {
	return nil
}

func (sc *Scene) Draw(dst draws.Image) {
}
