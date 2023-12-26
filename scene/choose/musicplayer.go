package choose

import (
	"time"

	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/times"
)

const waitDuration = 150 * time.Millisecond

const (
	EffectModeSilence = iota
	EffectModeFadeIn
	EffectModeNormal
	EffectModeFadeOut
)

type PreviewMusicPlayer struct {
	*audios.MusicPlayer
	// MusicVolume *float64
	StartTime  time.Time
	EffectMode int
}

func (s *Scene) updatePreviewMusic() {
	// It is fine to call Close at blank MusicPlayer.
	if s.MusicPlayer != nil {
		s.MusicPlayer.Close()
	}

	c := s.chart()
	fsys := c.MusicFS
	name := c.MusicFilename
	mp, _ := audios.NewMusicPlayerFromFile(fsys, name, 1)
	mp.SetVolume(s.MusicVolume)
	// MusicPlayer should be pointer so that it plays only once.
	s.PreviewMusicPlayer = PreviewMusicPlayer{
		MusicPlayer: &mp,
		StartTime:   times.Now(),
		EffectMode:  EffectModeSilence,
	}
}

// Memo: osu! seems fading music out when changing music.
func (s *Scene) HandleEffect() {
	const fadeInDuration = 1 * time.Second
	const fadeOutDuration = 3 * time.Second

	mp := s.PreviewMusicPlayer
	if mp.IsEmpty() {
		return
	}

	t := times.Since(mp.StartTime)
	switch mp.EffectMode {
	case EffectModeSilence:
		if t > waitDuration {
			mp.Play()
			mp.StartTime = times.Now()
			mp.EffectMode = EffectModeFadeIn
		}
	case EffectModeFadeIn:
		size := float64(t) / float64(fadeInDuration)
		vol := s.MusicVolume * size
		mp.SetVolume(vol)
		if t > fadeInDuration {
			mp.EffectMode = EffectModeNormal
		}
	case EffectModeNormal:
		if t > mp.Duration()-fadeOutDuration {
			mp.EffectMode = EffectModeFadeOut
		}
	case EffectModeFadeOut:
		size := float64(mp.Duration()-t) / float64(fadeOutDuration)
		vol := s.MusicVolume * size
		mp.SetVolume(vol)
		if t > mp.Duration() {
			mp.StartTime = times.Now()
			mp.Rewind()
			mp.EffectMode = EffectModeSilence
		}
	}
}
