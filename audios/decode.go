package audios

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
)

type StreamSeekCloser = beep.StreamSeekCloser
type Format = beep.Format

func DecodeFromFile(fsys fs.FS, name string) (streamer beep.StreamSeekCloser, format beep.Format, err error) {
	f, err := fsys.Open(name)
	if err != nil {
		return
	}
	// No close file. Streamer will close it.

	ext := filepath.Ext(name)
	switch strings.ToLower(ext) {
	case ".mp3":
		streamer, format, err = mp3.Decode(f)
	case ".wav":
		streamer, format, err = wav.Decode(f)
	case ".ogg":
		streamer, format, err = vorbis.Decode(f)
	}
	return
}

func isAudioFile(name string) bool {
	ext := filepath.Ext(name)
	switch strings.ToLower(ext) {
	case ".mp3", ".wav", ".ogg":
		return true
	}
	return false
}

// FormatFromFS returns the format of the first audio file in the file system.
// It is possible that there is no audio file in file system.
func FormatFromFS(fsys fs.FS) (format beep.Format, err error) {
	fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !isAudioFile(path) || !isFileSizeSmall(fsys, path) {
			return nil
		}

		_, format, err = DecodeFromFile(fsys, path)
		if err != nil {
			return err
		}
		return filepath.SkipDir // Skip further processing of directories
	})
	return
}

// vol: [0, 1] -> Volume: [-5, 0] => [1/32, 1]
func beepVolume(vol float64) float64 { return vol*5 - 5 }
