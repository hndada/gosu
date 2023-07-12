package audios

import (
	"fmt"
	"io/fs"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
)

type SoundMap struct {
	// sounds        map[string]Sound
	// fsys          fs.FS
	format        beep.Format
	buffer        *beep.Buffer
	startIndexMap map[string]int
	endIndexMap   map[string]int

	volumeScale   *float64
	resampleRatio float64
}

func NewSoundMap(fsys fs.FS, volumeScale *float64) SoundMap {
	sm := SoundMap{
		// fsys:          fsys,
		startIndexMap: make(map[string]int),
		endIndexMap:   make(map[string]int),

		volumeScale:   volumeScale,
		resampleRatio: 1,
	}
	// sm.format.SampleRate = 44100
	// streamers = make([]beep.StreamSeekCloser, 0)
	// sm.walkAndLoad(fsys, ".")
	fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("%+v\n", err)
		}
		if d.IsDir() || !isAudioFile(path) || !isSoundFileSize(fsys, path) {
			return nil
		}

		streamer, format, err := DecodeFromFile(fsys, path)
		if err != nil {
			return err
		}

		// Skipping resampling then making sounds a bit slower or faster
		// wouldn't make a big difference.

		// var resampled beep.Resampler
		// if format.SampleRate != defaultSampleRate {
		// resampled = beep.Resample(quality, format.SampleRate, defaultSampleRate, f)
		// }

		if sm.buffer == nil {
			sm.format = format
			sm.buffer = beep.NewBuffer(format)
		}
		sm.AppendSound(path, streamer)

		return nil
	})
	return sm
}

func isSoundFileSize(fsys fs.FS, name string) bool {
	// if filepath.Ext(path) == ".mp3" {
	// 	continue
	// }

	const maxSoundFileSize = 1 << 20 // 1MB

	f, err := fsys.Open(name)
	if err != nil {
		return false
	}

	info, err := f.Stat()
	if err != nil {
		return false
	}
	if info.Size() > maxSoundFileSize {
		return false
	}

	return true
}

// Len returns the number of sounds in SoundMap.
func (sm SoundMap) Len() int { return len(sm.startIndexMap) }

func (sm *SoundMap) AppendSound(name string, streamer beep.StreamSeekCloser) {
	sm.startIndexMap[name] = sm.buffer.Len()
	sm.buffer.Append(streamer)
	streamer.Close()
	sm.endIndexMap[name] = sm.buffer.Len()
}

func (sm *SoundMap) AppendSoundFromFile(fsys fs.FS, name string) error {
	streamer, _, err := DecodeFromFile(fsys, name)
	if err != nil {
		return err
	}
	sm.AppendSound(name, streamer)
	return nil
}

func (sm *SoundMap) SetResampleRatio(ratio float64) {
	sm.resampleRatio = ratio
}

func (sm SoundMap) Play(name string, vol float64) {
	start := sm.startIndexMap[name]
	end := sm.endIndexMap[name]
	streamer := sm.buffer.Streamer(start, end)

	resampler := beep.ResampleRatio(quality, sm.resampleRatio, streamer)
	beepVol := beepVolume(vol * (*sm.volumeScale))
	volume := &effects.Volume{Streamer: resampler, Base: 2, Volume: beepVol}
	speaker.Play(volume)
}
