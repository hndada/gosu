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
