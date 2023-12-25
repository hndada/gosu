package audios

import (
	"io/fs"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
)

type SoundMap struct {
	format        beep.Format
	buffer        *beep.Buffer
	startIndexMap map[string]int
	endIndexMap   map[string]int

	volumeScale   *float64
	resampleRatio float64
}

// It is possible that there is no sounds in file system.
// Hence, selecting format from the first met sound in file system
// then do beep.NewBuffer(format) may cause error.
func NewSoundMap(fsys fs.FS, format beep.Format, volumeScale *float64) SoundMap {
	sm := SoundMap{
		format:        format,
		buffer:        beep.NewBuffer(format),
		startIndexMap: make(map[string]int),
		endIndexMap:   make(map[string]int),

		volumeScale:   volumeScale,
		resampleRatio: 1,
	}

	fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !isAudioFile(path) || !isFileSizeSmall(fsys, path) {
			return nil
		}

		// Skipping resampling then making sounds a bit slower or faster
		// wouldn't make a big difference.
		streamer, _, err := DecodeFromFile(fsys, path)
		if err != nil {
			return err
		}
		// var resampled beep.Resampler
		// if format.SampleRate != defaultSampleRate {
		// resampled = beep.Resample(quality, format.SampleRate, defaultSampleRate, f)
		// }

		sm.Add(path, streamer)
		return nil
	})
	return sm
}

func isFileSizeSmall(fsys fs.FS, name string) bool {
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

func (sm *SoundMap) Add(name string, streamer beep.StreamSeekCloser) {
	sm.startIndexMap[name] = sm.buffer.Len()
	sm.buffer.Append(streamer)
	streamer.Close()
	sm.endIndexMap[name] = sm.buffer.Len()
}

func (sm *SoundMap) AddFromFile(fsys fs.FS, name string) error {
	streamer, _, err := DecodeFromFile(fsys, name)
	if err != nil {
		return err
	}
	sm.Add(name, streamer)
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
