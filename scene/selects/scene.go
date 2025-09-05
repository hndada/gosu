package selects

import "github.com/hndada/gosu/draws"

type Scene struct {
}

func NewScene() (*Scene, error) {
	return &Scene{}, nil
}

func (scn *Scene) Update() any {
	return nil
}

func (scn Scene) Draw(screen draws.Image) {
}

func (scn Scene) WindowTitle() string {
	return "Select a song - gosu"
}

func (scn Scene) DebugString() string {
	return "Selects Scene"
}
