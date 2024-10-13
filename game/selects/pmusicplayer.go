package selects

import (
	"io/fs"
	"time"

	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/times"
	"github.com/hndada/gosu/tween"
)

// musicPlayer loops music with fade in/out.
// (old memo: MusicPlayer should be pointer so that it plays only once.)
type PreviewMusicPlayer struct {
	*audios.MusicPlayer
	volume    *float64
	startTime time.Time
}

const waitDuration = 150 * time.Millisecond

var volumeTween tween.Tween

func init() {
	tw := tween.Tween{}
	tw.Add(0, 0, waitDuration, tween.EaseLinear)   // wait
	tw.Add(0, 1, 1*time.Second, tween.EaseLinear)  // fade in
	tw.Add(1, 1, 12*time.Second, tween.EaseLinear) // keep
	tw.Add(1, 0, 2*time.Second, tween.EaseLinear)  // fade out
	volumeTween = tw
}

func NewPreviewMusicPlayer(fsys fs.FS, name string, volume *float64) (PreviewMusicPlayer, error) {
	mp, err := audios.NewMusicPlayerFromFile(fsys, name)
	if err != nil {
		return PreviewMusicPlayer{}, err
	}

	return PreviewMusicPlayer{
		MusicPlayer: mp,
		volume:      volume,
		startTime:   times.Now(),
	}, nil
}

func (mp *PreviewMusicPlayer) Update() {
	volumeTween.Update()
	vol := *mp.volume * volumeTween.Value()
	mp.SetVolume(vol)

	t := times.Since(mp.startTime)
	switch {
	case t > waitDuration:
		mp.Play()
	case t > mp.Duration()+waitDuration:
		mp.Rewind()
		mp.startTime = times.Now()
	}
}

func (pmp *PreviewMusicPlayer) Close() {
	if pmp.MusicPlayer != nil {
		pmp.MusicPlayer.Close()
	}
}
