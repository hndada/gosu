package audios

import (
	"io/fs"
	"math/rand"

	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
)

type SoundEffect interface {
	Play(vol float64)
}

// SoundBag is for playing randomly one of the sounds.
type SoundPod struct {
	SoundPlayer
}

func NewSoundPod(fsys fs.FS, volumeScale *float64) SoundPod {
	return SoundPod{NewSoundPlayer(fsys, volumeScale)}
}

// Play plays a random sound from the sound pod.
func (sp SoundPod) Play(vol float64) {
	name := sp.randomName()
	start := sp.startIndexMap[name]
	end := sp.endIndexMap[name]
	streamer := sp.buffer.Streamer(start, end)

	beepVol := beepVolume(vol * (*sp.volumeScale))
	volume := &effects.Volume{Streamer: streamer, Base: 2, Volume: beepVol}
	speaker.Play(volume)
}

func (sp SoundPod) randomName() string {
	var count int
	randomIndex := rand.Intn(len(sp.startIndexMap))
	for name := range sp.startIndexMap {
		if count == randomIndex {
			return name
		}
		count++
	}
	return ""
}
