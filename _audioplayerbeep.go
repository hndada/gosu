package gosu

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"time"
)

type audioPlayer struct {
	sampleRate beep.SampleRate
	streamer   beep.StreamSeeker
	ctrl       *beep.Ctrl
	resampler  *beep.Resampler
	volume     *effects.Volume
}

func newAudioPlayer(streamer beep.StreamSeeker, sampleRate beep.SampleRate, timeRate float64) *audioPlayer {
	ctrl := &beep.Ctrl{Streamer: beep.Seq(streamer)}
	resampler := beep.ResampleRatio(4, timeRate, ctrl)
	volume := &effects.Volume{Streamer: resampler, Base: 2}
	return &audioPlayer{sampleRate, streamer, ctrl, resampler, volume}
}

func (ap *audioPlayer) play() {
	speaker.Play(ap.volume)
}
func (ap *audioPlayer) time() time.Duration {
	fmt.Println(ap.streamer.Position())
	return ap.sampleRate.D(ap.streamer.Position())
}

