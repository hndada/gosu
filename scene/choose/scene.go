package choose

import (
	"fmt"
	"io/fs"
	"time"

	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/scene"
)

// Todo: KeyHandler.OneAtATime option

// List is all in choose scene.
type Scene struct {
	*scene.Config
	*scene.Asset

	charts []*Chart
	audios.MusicPlayer
	currentNode *Node

	// queryTypeWriter input.TypeWriter
	// keySettings map[int][]string // todo: type aliasing for []string
	// list        *scene.List
	// listCursors []int
	// listDepth   int
}

func NewScene(root fs.FS) *Scene {
	s := &Scene{}
	var errs []error
	s.charts, errs = newCharts(root)
	for _, err := range errs {
		fmt.Println(err)
	}
	// s.listDepth = listDepthMusic
	return s
}

// sort by. // music name, level folder (+time?)
func (s *Scene) Update() any {
	// up down: move cursor
	// enter: listDepth++ select list item or play the chart.
	// back: listDepth--; 0

	// play preview music if music changes.
	// osu! seems fading music out when changing music.

	s.handleMusicPlayer()
	// scene's handlers
	return nil
}

// type FromChooseToPlay struct {
// 	cfg     *Config
// 	asset   *Asset
// 	fsys    fs.FS
// 	name    string
// 	rf      *osr.Format
// }

func (s *Scene) setMusicPlayer(fsys fs.FS, name string) {
	// Loop: wait + streamer

}
func (s *Scene) setBackground(fsys fs.FS, name string) {
	scene.NewBackgroundDrawer()
}

// handleMusicPlayer handles fade in/out effect.
func (s *Scene) handleMusicPlayer() {
	const waitDuration = 500 * time.Millisecond
	const fadeDuration = time.Second

	if s.MusicPlayer.Time() == waitDuration {
		s.MusicPlayer.FadeIn(fadeDuration, &s.MusicVolume)
	}
	if s.MusicPlayer.Time() == s.MusicPlayer.Duration()-fadeDuration {
		s.MusicPlayer.FadeOut(fadeDuration, &s.MusicVolume)
	}
}

func (s Scene) DebugString() string {
	return ""
}
