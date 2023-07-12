package audios

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
)

// type Sound []byte
type Sound beep.Streamer

type SoundPlayer struct {
	// sounds        map[string]Sound
	// fsys          fs.FS
	format        beep.Format
	buffer        *beep.Buffer
	startIndexMap map[string]int
	endIndexMap   map[string]int

	volumeScale   *float64
	resampleRatio float64
}

func NewSoundPlayer(fsys fs.FS, volumeScale *float64) SoundPlayer {
	sp := SoundPlayer{
		// fsys:          fsys,
		startIndexMap: make(map[string]int),
		endIndexMap:   make(map[string]int),

		volumeScale:   volumeScale,
		resampleRatio: 1,
	}
	// sp.format.SampleRate = 44100
	// streamers = make([]beep.StreamSeekCloser, 0)
	sp.walkAndLoad(fsys, ".")
	return sp
}
func (sp *SoundPlayer) walkAndLoad(root fs.FS, dir string) error {
	entries, err := fs.ReadDir(root, dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			if err := sp.walkAndLoad(root, path); err != nil {
				return fmt.Errorf("error in walkDir at %q: %v", path, err)
			}
		} else {
			if !isAudioFile(path) {
				continue
			}

			streamer, format, err := DecodeFromFile(root, path)
			if err != nil {
				return err
			}

			// Resample if the streamer has different sample rate
			// var resampled beep.Resampler
			// if format.SampleRate != defaultSampleRate {
			// resampled = beep.Resample(quality, format.SampleRate, defaultSampleRate, f)
			// }

			// Skipping resampling then making sounds a bit slower or faster
			// wouldn't make a big difference.

			if sp.buffer == nil {
				sp.format = format
				sp.buffer = beep.NewBuffer(format)
			}
			sp.AppendSound(path, streamer)
		}
	}
	return nil
}

func (sp *SoundPlayer) AppendSound(name string, streamer beep.StreamSeekCloser) {
	sp.startIndexMap[name] = sp.buffer.Len()
	sp.buffer.Append(streamer)
	streamer.Close()
	sp.endIndexMap[name] = sp.buffer.Len()
}

func (sp *SoundPlayer) AppendSoundFromFile(fsys fs.FS, name string) error {
	streamer, _, err := DecodeFromFile(fsys, name)
	if err != nil {
		return err
	}
	sp.AppendSound(name, streamer)
	return nil
}

func (sp SoundPlayer) Play(name string, vol float64) {
	start := sp.startIndexMap[name]
	end := sp.endIndexMap[name]
	streamer := sp.buffer.Streamer(start, end)

	resampler := beep.ResampleRatio(quality, sp.resampleRatio, streamer)
	beepVol := beepVolume(vol * (*sp.volumeScale))
	volume := &effects.Volume{Streamer: resampler, Base: 2, Volume: beepVol}
	speaker.Play(volume)
}

func (sp *SoundPlayer) SetResampleRatio(ratio float64) {
	sp.resampleRatio = ratio
}
