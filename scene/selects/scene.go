package selects

import "github.com/hndada/gosu/scene"

// Component is basically EventHandler.
type Scene struct {
	// List consists of folders and items.
	folders []chartFolder
	index   int
}

// It is fine to call Close at blank MusicPlayer.

// Avoid embedding scene.Options directly.
// Pass options as pointers for syncing and saving memory.
func NewScene(res *scene.Resources, opts *scene.Options) *Scene {
	return &Scene{}
}
