package choose

import (
	"time"

	"github.com/hndada/gosu/audios"
)

const waitDuration = 500 * time.Millisecond

const (
	// EffectModeSilence = iota
	EffectModeNormal = iota
	EffectModeFadeIn
	EffectModeFadeOut
)

type PreviewMusicPlayer struct {
	audios.MusicPlayer
	// MusicVolume *float64
	StartTime  time.Time
	EffectMode int
}

func (s *Scene) updatePreviewMusic() *PreviewMusicPlayer {
	// It is fine to call Close at blank MusicPlayer.
	s.MusicPlayer.Close()

	c := s.chart()
	fsys := c.MusicFS
	name := c.MusicFilename
	mp, _ := audios.NewMusicPlayerFromFile(fsys, name, 1)

	return &PreviewMusicPlayer{
		MusicPlayer: mp,
		StartTime:   time.Now().Add(waitDuration),
		EffectMode:  EffectModeFadeIn,
	}
}

// Memo: osu! seems fading music out when changing music.
func (s *Scene) HandleEffect() {
	const fadeDuration = time.Second
	mp := s.PreviewMusicPlayer
	if mp.IsEmpty() {
		return
	}

	t := time.Since(mp.StartTime)
	switch mp.EffectMode {
	case EffectModeFadeIn:
		if t > fadeDuration {
			mp.EffectMode = EffectModeNormal
		}
	case EffectModeNormal:
		if t > mp.Duration()-fadeDuration {
			mp.FadeOut(fadeDuration, &s.MusicVolume)
			mp.EffectMode = EffectModeFadeOut
		}
	case EffectModeFadeOut:
		if t > mp.Duration()+waitDuration {
			mp.Rewind()
			mp.StartTime = time.Now()
			mp.FadeIn(fadeDuration, &s.MusicVolume)
			mp.EffectMode = EffectModeFadeIn
		}
	}
}
