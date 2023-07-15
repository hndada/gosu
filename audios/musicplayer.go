package audios

import (
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

	// format    beep.Format
	// volume    float64
	// resampleRatio float64
	// pauseChannel  chan struct{}
	// resumeChannel chan struct{}
	// paused        bool
	played bool
}

// var isSpeakerInit bool

const defaultSampleRate beep.SampleRate = 44100
const quality = 4

func init() {
	speaker.Init(defaultSampleRate, defaultSampleRate.N(time.Second/30))
}

func NewMusicPlayer(f beep.StreamSeekCloser, format beep.Format, ratio float64) MusicPlayer {
	done := make(chan bool)
	callback := beep.Callback(func() { done <- true })
	ctrl := &beep.Ctrl{Streamer: beep.Seq(f, callback)}

	resampler := beep.ResampleRatio(quality, 1, ctrl) // for applying ctrl
	if format.SampleRate != defaultSampleRate {
		resampler = beep.Resample(quality, format.SampleRate, defaultSampleRate, ctrl)
	}
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

		// format:    format,
		// volume:    1,
		// resampleRatio: ratio,
		// pauseChannel:  make(chan struct{}),
		// resumeChannel: make(chan struct{}),
		// paused:        false,
	}
}

func NewMusicPlayerFromFile(fsys fs.FS, name string, ratio float64) (MusicPlayer, error) {
	f, format, err := DecodeFromFile(fsys, name)
	if err != nil {
		return MusicPlayer{}, err
	}
	return NewMusicPlayer(f, format, ratio), nil
}

func (mp *MusicPlayer) Play() {
	// if !isSpeakerInit {
	// 	speaker.Init(mp.format.SampleRate, mp.format.SampleRate.N(time.Second/30))
	// 	isSpeakerInit = true
	// }
	if mp.played {
		return
	}
	speaker.Play(mp.volume)
	for range mp.done {
		mp.Close()
		return
	}
	// for {
	// 	select {
	// 	case <-done:
	// 		mp.Close()
	// 		return
	// 		// case <-mp.pauseChannel:
	// 		// case <-mp.resumeChannel:
	// 	}
	// }
}

func (mp MusicPlayer) CurrentTime() time.Duration {
	return defaultSampleRate.D(mp.streamer.Position())
}

func (mp *MusicPlayer) SetVolume(vol float64) {
	mp.volume.Volume = beepVolume(vol)
	if vol <= 0.001 { // 0.1%
		mp.volume.Silent = true
	} else {
		mp.volume.Silent = false
	}
}

// vol: [0, 1] -> Volume: [-5, 0] => [1/32, 1]
func beepVolume(vol float64) float64 { return vol*5 - 5 }

func (mp MusicPlayer) IsPaused() bool { return mp.ctrl.Paused }
func (mp *MusicPlayer) Pause() {
	// if !mp.paused {
	// 	mp.pauseChannel <- struct{}{}
	// }
	// mp.paused = true
	speaker.Lock()
	mp.ctrl.Paused = true
	speaker.Unlock()
}

func (mp *MusicPlayer) Resume() {
	// if mp.paused {
	// 	mp.resumeChannel <- struct{}{}
	// }
	// mp.paused = false
	speaker.Lock()
	mp.ctrl.Paused = false
	speaker.Unlock()
}

func (mp MusicPlayer) Speed() float64 { return mp.resampler.Ratio() }

func (mp *MusicPlayer) SetResampleRatio(ratio float64) {
	speaker.Lock()
	mp.resampler.SetRatio(ratio)
	speaker.Unlock()
	// mp.resampleRatio = ratio
}

func (mp *MusicPlayer) Close() {
	speaker.Lock()
	speaker.Clear()
	speaker.Unlock()
	mp.streamer.Close()
}

// Always call beep.ResampleRatio to set the ratio even if ratio is 1,
// because return type of beep.ResampleRatio is *beep.Resampler
// whereas type of f is beep.StreamSeekCloser.
