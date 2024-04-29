package audios

import (
	"fmt"
	"io/fs"
	"math/rand"
	"path/filepath"
	"strings"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/speaker"
	"github.com/hndada/gosu/util"
)

type SoundPlayer struct {
	buffers          map[string]SoundBuffer
	bufferNames      []string
	soundVolumeScale *float64
	PlaybackRate     float64
}

func NewSoundPlayer(scale *float64) SoundPlayer {
	buffers := map[string]SoundBuffer{"": newSoundBuffer()}
	bufferNames := []string{"default"}
	return SoundPlayer{
		buffers:          buffers,
		bufferNames:      bufferNames,
		soundVolumeScale: scale,
		PlaybackRate:     1,
	}
}

func (sp *SoundPlayer) AddFile(fsys fs.FS, name string) error {
	data, err := fs.ReadFile(fsys, name)
	if err != nil {
		return fmt.Errorf("read file %s: %w", name, err)
	}
	sb := sp.buffers["default"]
	return sb.add(data, name)
}

// AddDir adds audio files in the directory to SoundPlayer.
func (sp *SoundPlayer) AddDir(fsys fs.FS, name string) error {
	base := strings.TrimSuffix(name, filepath.Ext(name))
	paths := util.DirElements(fsys, base)
	if len(paths) == 0 {
		return sp.AddFile(fsys, name)
	}

	sb := newSoundBuffer()
	for _, path := range paths {
		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return fmt.Errorf("read file %s: %w", path, err)
		}
		_ = sb.add(data, path)
	}
	sp.buffers[base] = sb
	sp.bufferNames = append(sp.bufferNames, base)
	return nil
}

// Count returns the number of kinds of sounds in SoundPlayer.
func (sp SoundPlayer) Count() int {
	defBuf := sp.buffers["default"]
	return len(defBuf.keys) + len(sp.bufferNames) - 1
}

func (sp SoundPlayer) containsName(name string) bool {
	for _, bn := range sp.bufferNames {
		if bn == name {
			return true
		}
	}
	return false
}

func (sp SoundPlayer) Play(name string) { sp.PlayWithVolume(name, 1) }
func (sp SoundPlayer) PlayWithVolume(name string, vol float64) {
	var s beep.Streamer
	if sp.containsName(name) {
		sb := sp.buffers[name]
		i := rand.Intn(len(sb.keys))
		key := sb.keys[i]
		s = sb.buffer.Streamer(sb.starts[key], sb.ends[key])
	} else {
		sb := sp.buffers["default"]
		s = sb.buffer.Streamer(sb.starts[name], sb.ends[name])
	}

	if sp.PlaybackRate != 1 {
		s = beep.ResampleRatio(quality, sp.PlaybackRate, s)
	}
	vol *= *sp.soundVolumeScale
	volume := &effects.Volume{Streamer: s, Base: 2, Volume: beepVolume(vol)}
	speaker.Play(volume)
}

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
