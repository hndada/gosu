package selects

import "github.com/hndada/gosu/scene"

// Component is basically EventHandler.
type Scene struct {
}

// Avoid embedding scene.Options directly.
// Pass options as pointers for syncing and saving memory.
func NewScene(res *scene.Resources, opts *scene.Options) *Scene {
	return &Scene{}
}
