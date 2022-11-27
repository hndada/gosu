package audios

import (
	"io"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hndada/gosu/input"
)

type MusicPlayer struct {
	*audio.Player
	*Timer
	Volume   *float64
	volume   float64
	pause    bool
	pauseKey input.Key
	Closer   func() error
}

func NewMusicPlayer(fsys fs.FS, name string, t *Timer, vol *float64, k input.Key) (MusicPlayer, error) {
	p, close, err := newAudioPlayer(fsys, name)
	if err != nil {
		return MusicPlayer{}, err
	}
	p.SetVolume(*vol)
	// player.SetBufferSize(100 * time.Millisecond)
	return MusicPlayer{
		Player:   p,
		Timer:    t,
		Volume:   vol,
		volume:   *vol,
		pause:    false,
		pauseKey: k,
		Closer:   close,
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
func (p MusicPlayer) IsValid() bool { return p.Player != nil }
func (p *MusicPlayer) Update() {
	if !p.IsValid() {
		return
	}
	if inpututil.IsKeyJustPressed(p.pauseKey) {
		if p.pause {
			p.Player.Play()
		} else {
			p.Player.Pause()
		}
		p.pause = !p.pause
	}
	if p.pause {
		return
	}
	if vol := *p.Volume; p.volume != vol {
		p.volume = vol
		p.Player.SetVolume(vol)
	}
	if p.Now == 0+p.Offset {
		p.Player.Play()
	}
	// if p.Now == 150+p.Offset {
	// 	p.Player.Seek(time.Duration(150) * time.Millisecond)
	// }
	if p.IsDone() {
		p.Close()
	}
}
func (p MusicPlayer) Close() {
	if p.IsValid() {
		p.Player.Close()
		p.Closer()
	}
}
