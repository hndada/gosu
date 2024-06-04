package selects

import (
	"io/fs"
	"time"

	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/times"
)

// musicPlayer loops music with fade in/out.
// (old memo: MusicPlayer should be pointer so that it plays only once.)
type PreviewMusicPlayer struct {
	audios.MusicPlayer
	volume       *float64
	waitDuration time.Duration
	tween        times.Tween
	startTime    time.Time
}

func NewPreviewMusicPlayer(fsys fs.FS, name string, volume *float64) (*PreviewMusicPlayer, error) {
	const waitDuration = 150 * time.Millisecond
	tw := times.Tween{}
	tw.Add(0, 0, waitDuration, times.EaseLinear)   // wait
	tw.Add(0, 1, 1*time.Second, times.EaseLinear)  // fade in
	tw.Add(1, 1, 12*time.Second, times.EaseLinear) // keep
	tw.Add(1, 0, 2*time.Second, times.EaseLinear)  // fade out

	mp, err := audios.NewMusicPlayerFromFile(fsys, name)
	if err != nil {
		return nil, err
	}

	return &PreviewMusicPlayer{
		MusicPlayer:  mp,
		volume:       volume,
		waitDuration: waitDuration,
		tween:        tw,
		startTime:    times.Now(),
	}, nil
}

func (mp *PreviewMusicPlayer) Update() {
	vol := *mp.volume * mp.tween.Current()
	mp.SetVolume(vol)

	t := times.Since(mp.startTime)
	switch {
	case t > mp.waitDuration:
		mp.Play()
	case t > mp.Duration()+mp.waitDuration:
		mp.Rewind()
		mp.startTime = times.Now()
	}
}
