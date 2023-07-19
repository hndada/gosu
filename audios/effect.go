package audios

import (
	"time"

	"github.com/faiface/beep/speaker"
)

func (mp MusicPlayer) PlaybackRate() float64 { return mp.resampler.Ratio() }

func (mp *MusicPlayer) SetPlaybackRate(ratio float64) {
	speaker.Lock()
	mp.resampler.SetRatio(ratio)
	speaker.Unlock()
}

func (mp *MusicPlayer) FadeIn(duration time.Duration, volume *float64) {
	const stepCount = 100
	go func() {
		for i := 0; i < stepCount; i++ {
			size := float64(i) / stepCount
			vol := *volume * size
			mp.SetVolume(vol)
			time.Sleep(duration / stepCount)
		}
	}()
}

func (mp *MusicPlayer) FadeOut(duration time.Duration, volume *float64) {
	const stepCount = 100
	go func() {
		for i := 0; i < stepCount; i++ {
			size := 1 - float64(i)/stepCount
			vol := *volume * size
			mp.SetVolume(vol)
			time.Sleep(duration / stepCount)
		}
	}()
}
