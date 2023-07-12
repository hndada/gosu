package audios

import (
	"io"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

type MusicPlayer struct {
	*audio.Player
	Closer func() error
}

func NewMusicPlayer(fsys fs.FS, name string) (MusicPlayer, error) {
	p, close, err := newAudioPlayer(fsys, name)
	if err != nil {
		return MusicPlayer{}, err
	}
	// player.SetBufferSize(100 * time.Millisecond)
	return MusicPlayer{
		Player: p,
		Closer: close,
	}, nil
}
func newAudioPlayer(fsys fs.FS, name string) (p *audio.Player, close func() error, err error) {
	var streamer io.ReadSeeker
	streamer, close, err = decode(fsys, name)
	if err != nil {
		return
	}
	p, err = Context.NewPlayer(streamer)
	if err != nil {
		return
	}
	return
}
func (p MusicPlayer) Close() {
	p.Player.Close()
	p.Closer()
}
