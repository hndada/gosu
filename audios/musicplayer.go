package audios

import (
	"fmt"
	"io/fs"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
)

type MusicPlayer struct {
	streamer  beep.StreamSeekCloser
	ctrl      *beep.Ctrl
	resampler *beep.Resampler
	volume    *effects.Volume
	done      chan bool
	played    bool
}

const defaultSampleRate beep.SampleRate = 44100
const quality = 4

func init() {
	speaker.Init(defaultSampleRate, defaultSampleRate.N(time.Second/30))
}

func NewMusicPlayer(f beep.StreamSeekCloser, format beep.Format, ratio float64) MusicPlayer {
	done := make(chan bool)
	callback := beep.Callback(func() { done <- true })
	ctrl := &beep.Ctrl{Streamer: beep.Seq(f, callback)}

	// No ratio change. Is is for applying ctrl.
	resampler := beep.ResampleRatio(quality, 1, ctrl)
	// Do the actual resample here if sample rate is different.
	if format.SampleRate != defaultSampleRate {
		resampler = beep.Resample(quality, format.SampleRate, defaultSampleRate, ctrl)
	}
	// Change the ratio if it is not 1.
	if ratio != 1 {
		resampler = beep.ResampleRatio(quality, ratio, resampler)
	}

	volume := &effects.Volume{Streamer: resampler, Base: 2}
	return MusicPlayer{
		streamer:  f,
		ctrl:      ctrl,
		resampler: resampler,
		volume:    volume,
		done:      done,
	}
}

func NewMusicPlayerFromFile(fsys fs.FS, name string, ratio float64) (MusicPlayer, error) {
	f, format, err := DecodeFromFile(fsys, name)
	if err != nil {
		return MusicPlayer{}, fmt.Errorf("decode %s: %w", name, err)
	}
	return NewMusicPlayer(f, format, ratio), nil
}

func (mp *MusicPlayer) Play() {
	if mp.played {
		return
	}
	// speaker.Lock()
	speaker.Play(mp.volume)
	mp.played = true
	// speaker.Unlock()
}

func (mp *MusicPlayer) Rewind() {
	speaker.Lock()
	mp.streamer.Seek(0)
	speaker.Unlock()
}

func (mp MusicPlayer) IsEmpty() bool { return mp.streamer == nil }

func (mp MusicPlayer) IsPlayed() bool { return mp.played }

func (mp MusicPlayer) Time() time.Duration {
	return defaultSampleRate.D(mp.streamer.Position())
}
func (mp MusicPlayer) Duration() time.Duration {
	return defaultSampleRate.D(mp.streamer.Len())
}

func (mp MusicPlayer) PlaybackRate() float64 { return mp.resampler.Ratio() }

func (mp *MusicPlayer) SetPlaybackRate(ratio float64) {
	speaker.Lock()
	mp.resampler.SetRatio(ratio)
	speaker.Unlock()
}

func (mp *MusicPlayer) SetVolume(vol float64) {
	speaker.Lock()
	mp.volume.Volume = beepVolume(vol)
	if vol <= 0.001 { // 0.1%
		mp.volume.Silent = true
	} else {
		mp.volume.Silent = false
	}
	speaker.Unlock()
}

func (mp MusicPlayer) IsPaused() bool { return mp.ctrl.Paused }

func (mp *MusicPlayer) Pause() {
	speaker.Lock()
	mp.ctrl.Paused = true
	speaker.Unlock()
}

func (mp *MusicPlayer) Resume() {
	speaker.Lock()
	mp.ctrl.Paused = false
	speaker.Unlock()
}

func (mp *MusicPlayer) Close() {
	speaker.Clear()
	if mp != nil && mp.streamer != nil {
		mp.streamer.Close()
	}
}
