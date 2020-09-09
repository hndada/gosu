package gosu

import (
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
	return ap.sampleRate.D(ap.streamer.Position())
}

// NewSceneMania에서 쓰던거
// {
// f, err := os.Open(s.chart.AbsPath(s.chart.AudioFilename))
// if err != nil {
// panic(err)
// }
// var streamer beep.StreamSeekCloser
// var format beep.Format
// switch strings.ToLower(filepath.Ext(s.chart.AudioFilename)) {
// case ".mp3":
// streamer, format, err = mp3.Decode(f)
// if err != nil {
// panic(err)
// }
// }
// _ = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
// s.audioPlayer = NewAudioPlayer(streamer, format.SampleRate, s.mods.TimeRate)
// }