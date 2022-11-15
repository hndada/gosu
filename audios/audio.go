package audios

import "io/fs"

type MusicPlayer struct {
	*Timer
	Volume *float64
	volume float64
	Player
	Closer func() error
	pause  bool
}

func NewMusicPlayer(fsys fs.FS, name string, timer *Timer, volume *float64) (MusicPlayer, error) {
	player, close, err := NewPlayer(fsys, name)
	if err != nil {
		return MusicPlayer{}, err
	}
	player.SetVolume(*volume)
	// player.SetBufferSize(100 * time.Millisecond)
	return MusicPlayer{
		Timer:  timer,
		Volume: volume,
		volume: *volume,
		Player: player,
		Closer: close,
	}, nil
}

func (p *MusicPlayer) Update() {
	if !p.Player.IsValid() {
		return
	}
	// Calling SetVolume in every Update is fine, confirmed by the author, by the way.
	if vol := *p.Volume; p.volume != vol {
		p.volume = vol
		p.Player.SetVolume(vol)
	}
	if p.Timer.Pause {
		if !p.pause {
			p.Player.Pause()
			p.pause = true
		}
	} else {
		if p.pause {
			p.Player.Play()
			p.pause = false
		}
	}
	if p.pause {
		return
	}
	if p.Now == 0+(*p.Offset) {
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
	if p.Player.IsValid() {
		p.Player.Close()
		p.Closer()
	}
}

// A player for sample sound is generated at a place.
// Todo: implement
type SoundPlayer struct {
	Volume *float64
	// Player *audio.Player
}
