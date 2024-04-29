package audios

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"

	"github.com/gopxl/beep"
)

var defaultFormat = beep.Format{
	SampleRate:  defaultSampleRate,
	NumChannels: 2,
	Precision:   4,
}

type SoundBuffer struct {
	buffer *beep.Buffer
	starts map[string]int // start index
	ends   map[string]int // end index
	keys   []string
}

func newSoundBuffer() SoundBuffer {
	return SoundBuffer{
		buffer: beep.NewBuffer(defaultFormat),
		starts: make(map[string]int),
		ends:   make(map[string]int),
	}
}

// It is possible for empty string to be a key of a map.
// https://go.dev/play/p/nn-peGAjawW
func (sb *SoundBuffer) add(data []byte, name string) error {
	r := bytes.NewReader(data)
	rc := io.NopCloser(r)

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

	sb.starts[name] = sb.buffer.Len()
	sb.buffer.Append(streamer)
	sb.ends[name] = sb.buffer.Len()
	sb.keys = append(sb.keys, name)
	return nil
}
