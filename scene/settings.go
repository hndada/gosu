package scene

import (
	"strings"
)

type Settings struct {
	MusicVolume          float64
	SoundVolume          float64
	BackgroundBrightness float64
	Offset               int64
	DebugPrint           bool

	MusicRoots  []string
	CursorScale float64
	ClearScale  float64
}

var TheSettings = Settings{
	MusicVolume:          0.50,
	SoundVolume:          0.50,
	BackgroundBrightness: 0.6,
	Offset:               -20,
	DebugPrint:           true,

	MusicRoots:  []string{"music"},
	CursorScale: 0.1,
	ClearScale:  0.5,
}

func NormalizeMusicRoots(roots []string) []string {
	if len(roots) == 0 {
		return []string{"music"}
	}

	// Leading dot and slash is not allowed in fs.
	for i, name := range roots {
		name = strings.TrimPrefix(name, "..")
		name = strings.TrimPrefix(name, ".")
		name = strings.TrimPrefix(name, "/")
		name = strings.TrimPrefix(name, "\\")
		roots[i] = name
	}
	return roots
}
