package audios

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/generators"
	"github.com/gopxl/beep/speaker"
)

type MusicPlayer struct {
	seekCloser beep.StreamSeekCloser // for seek and close
	format     beep.Format           // for duration
	ctrl       *beep.Ctrl            // for pause
	streamer   *beep.Resampler       // main streamer
	volume     *effects.Volume       // for volume
}

// I guess NewMusicPlayer should return pointer, so that
// it can be controlled by pointer receiver.
func NewMusicPlayer(rc io.ReadCloser, ext string) (*MusicPlayer, error) {
	seekCloser, format, err := Decode(rc, ext)
	if err != nil {
		return nil, fmt.Errorf("decode %s: %w", ext, err)
	}
	// done := make(chan bool)
	// callback := beep.Callback(func() { done <- true })
	// ctrl := &beep.Ctrl{Streamer: beep.Seq(seekCloser, callback)}
	ctrl := &beep.Ctrl{Streamer: seekCloser}
	streamer := beep.Resample(quality, format.SampleRate, defaultSampleRate, ctrl)
	volume := &effects.Volume{Streamer: streamer, Base: 2}
	return &MusicPlayer{
		seekCloser: seekCloser,
		format:     format,
		ctrl:       ctrl,
		streamer:   streamer,
		volume:     volume,
	}, nil
}

func NewMusicPlayerFromFile(fsys fs.FS, name string) (*MusicPlayer, error) {
	ext := filepath.Ext(name)
	f, err := fsys.Open(name)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", name, err)
	}
	return NewMusicPlayer(f, ext)
}

func (mp MusicPlayer) IsEmpty() bool { return mp.seekCloser == nil }

func (mp MusicPlayer) Play() {
	if mp.IsEmpty() {
		return
	}
	speaker.Play(mp.volume)
}

func (mp *MusicPlayer) Rewind() {
	if mp.IsEmpty() {
		return
	}
	mp.seekCloser.Seek(0)
}

func (mp MusicPlayer) Current() time.Duration {
	if mp.IsEmpty() {
		return 0
	}
	sr := mp.format.SampleRate
	return sr.D(mp.seekCloser.Position())
}

func (mp MusicPlayer) Duration() time.Duration {
	if mp.IsEmpty() {
		return 0
	}
	sr := mp.format.SampleRate
	return sr.D(mp.seekCloser.Len())
}

func (mp MusicPlayer) PlaybackRate() float64 {
	if mp.IsEmpty() {
		return 1
	}
	return mp.streamer.Ratio()
}

func (mp *MusicPlayer) SetPlaybackRate(rate float64) {
	if mp.IsEmpty() {
		return
	}
	mp.streamer.SetRatio(rate)
}

// beepVolume converts volume from [0, 1] to [-5, 0].
// [-5, 0] is log scale.
func beepVolume(vol float64) float64 { return vol*5 - 5 }

func (mp *MusicPlayer) SetVolume(vol float64) {
	if mp.IsEmpty() {
		return
	}
	mp.volume.Volume = beepVolume(vol)
	if vol <= 0.001 { // 0.1%
		mp.volume.Silent = true
	} else {
		mp.volume.Silent = false
	}
}

func (mp MusicPlayer) IsPaused() bool {
	if mp.IsEmpty() {
		return false
	}
	return mp.ctrl.Paused
}

// Lock is required when modifying beep.Ctrl.
func (mp *MusicPlayer) Pause() {
	if mp.IsEmpty() {
		return
	}
	speaker.Lock()
	mp.ctrl.Paused = true
	speaker.Unlock()
}

// Lock is required when modifying beep.Ctrl.
func (mp *MusicPlayer) Resume() {
	if mp.IsEmpty() {
		return
	}
	speaker.Lock()
	mp.ctrl.Paused = false
	speaker.Unlock()
}

func (mp *MusicPlayer) Close() {
	if mp.IsEmpty() {
		return
	}
	speaker.Clear()
	mp.seekCloser.Close()
}

func NewSilence(duration time.Duration) beep.Streamer {
	num := defaultSampleRate.N(duration)
	return generators.Silence(num)
}
