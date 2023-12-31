package audios

import (
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"path/filepath"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/speaker"
)

var defaultFormat = beep.Format{
	SampleRate:  defaultSampleRate,
	NumChannels: 2,
	Precision:   4,
}

type SoundPlayer struct {
	buffer       *beep.Buffer
	starts       map[string]int // start index
	ends         map[string]int // end index
	keys         []string
	PlaybackRate float64
}

func NewSoundPlayer() SoundPlayer {
	return SoundPlayer{
		buffer:       beep.NewBuffer(defaultFormat),
		starts:       make(map[string]int),
		ends:         make(map[string]int),
		PlaybackRate: 1,
	}
}

// It is possible for empty string to be a key of a map.
// https://go.dev/play/p/nn-peGAjawW
func (sp *SoundPlayer) Add(rc io.ReadCloser, name string) error {
	defer rc.Close()

	// Declaring streamer's type explicitly is for
	// assigning beep.Resampler to it.
	var streamer beep.Streamer
	streamer, format, err := Decode(rc, filepath.Ext(name))
	if err != nil {
		return fmt.Errorf("decode %s: %w", name, err)
	}
	if format.SampleRate != defaultSampleRate {
		old := format.SampleRate
		new := defaultSampleRate
		streamer = beep.Resample(quality, old, new, streamer)
	}

	sp.starts[name] = sp.buffer.Len()
	sp.buffer.Append(streamer)
	sp.ends[name] = sp.buffer.Len()
	sp.keys = append(sp.keys, name)
	return nil
}

func (sp *SoundPlayer) AddFromFile(fsys fs.FS, name string) error {
	f, err := fsys.Open(name)
	if err != nil {
		return fmt.Errorf("open %s: %w", name, err)
	}
	return sp.Add(f, name)
}

// Count returns the number of sounds in SoundPlayer.
func (sp SoundPlayer) Count() int { return len(sp.starts) }

func (sp SoundPlayer) Play(name string, vol float64) {
	var s beep.Streamer = sp.buffer.Streamer(sp.starts[name], sp.ends[name])
	if sp.PlaybackRate != 1 {
		s = beep.ResampleRatio(quality, sp.PlaybackRate, s)
	}
	volume := &effects.Volume{Streamer: s, Base: 2, Volume: beepVolume(vol)}
	speaker.Play(volume)
}

func (sp SoundPlayer) PlayRandom(vol float64) {
	i := rand.Intn(len(sp.keys))
	sp.Play(sp.keys[i], vol)
}

// func isFileSizeSmall(fsys fs.FS, name string) bool {
// 	// if filepath.Ext(path) == ".mp3" {
// 	// 	continue
// 	// }

// 	const maxSoundFileSize = 1 << 20 // 1MB

// 	f, err := fsys.Open(name)
// 	if err != nil {
// 		return false
// 	}

// 	info, err := f.Stat()
// 	if err != nil {
// 		return false
// 	}
// 	if info.Size() > maxSoundFileSize {
// 		return false
// 	}

// 	return true
// }

// func NewSoundMap(fsys fs.FS, format beep.Format) SoundMap {
// 	fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
// 		if err != nil {
// 			return err
// 		}
// 		if d.IsDir() || !isAudioFile(path) || !isFileSizeSmall(fsys, path) {
// 			return nil
// 		}

// 		// Skipping resampling then making sounds a bit slower or faster
// 		// wouldn't make a big difference.
// 		streamer, _, err := DecodeFromFile(fsys, path)
// 		if err != nil {
// 			return err
// 		}
// 		// var resampled beep.Resampler
// 		// if format.SampleRate != defaultSampleRate {
// 		// resampled = beep.Resample(quality, format.SampleRate, defaultSampleRate, f)
// 		// }

// 		sm.Add(path, streamer)
// 		return nil
// 	})
// 	return sm
// }
