package audios

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/vorbis"
	"github.com/gopxl/beep/wav"
)

const (
	defaultSampleRate beep.SampleRate = 44100 // 48000
	quality           int             = 4
)

func init() {
	speaker.Init(defaultSampleRate, defaultSampleRate.N(time.Second/20))
}

// streamer as StreamerSeekCloser will close it.
func Decode(rc io.ReadCloser, ext string) (beep.StreamSeekCloser, beep.Format, error) {
	switch ext {
	case ".mp3":
		return mp3.Decode(rc)
	case ".wav":
		return wav.Decode(rc)
	case ".ogg":
		return vorbis.Decode(rc)
	}
	err := fmt.Errorf("decode %s: %s", ext, "unsupported extension")
	return nil, beep.Format{}, err
}

func DecodeFromFile(fsys fs.FS, name string) (beep.StreamCloser, beep.Format, error) {
	f, err := fsys.Open(name)
	if err != nil {
		// err = &fs.PathError{Op: "open", Path: name, Err: err}
		return nil, beep.Format{}, err
	}
	return Decode(f, filepath.Ext(name))
}

// func isAudioFile(name string) bool {
// 	ext := filepath.Ext(name)
// 	switch strings.ToLower(ext) {
// 	case ".mp3", ".wav", ".ogg":
// 		return true
// 	}
// 	return false
// }

// // FormatFromFS returns the format of the first audio file in the file system.
// // It is possible that there is no audio file in file system.
// func FormatFromFS(fsys fs.FS) (Format, error) {
// 	var (
// 		format beep.Format
// 		err    error
// 	)
// 	fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
// 		if err != nil {
// 			return err
// 		}
// 		if d.IsDir() || !isAudioFile(path) || !isFileSizeSmall(fsys, path) {
// 			return nil
// 		}

// 		_, format, err = DecodeFromFile(fsys, path)
// 		if err != nil {
// 			return err
// 		}
// 		return filepath.SkipDir // Skip further processing of directories
// 	})
// 	if err != nil {
// 		err = fmt.Errorf("format from fs: %w", err)
// 	}
// 	return format, err
// }
