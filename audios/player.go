package audios

import (
	"io"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// It is embedded instead of aliasing, desiring adding new feature: e.g., time rate.
type Player struct {
	*audio.Player
}

func NewPlayer(fsys fs.FS, name string) (player Player, close func() error, err error) {
	var streamer io.ReadSeeker
	streamer, close, err = decode(fsys, name)
	if err != nil {
		return
	}
	p, err := Context.NewPlayer(streamer)
	if err != nil {
		return
	}
	player = Player{p}
	// var bufferSize = 100 * time.Millisecond
	// p.SetBufferSize(bufferSize)
	return
}
func (p Player) IsValid() bool { return p.Player != nil }
