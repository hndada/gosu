package choose

import (
	"time"

	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/scene"
)

// Todo: KeyHandler.OneAtATime option

// List is all in choose scene.
type Scene struct {
	*scene.Config
	*scene.Asset

	// queryTypeWriter input.TypeWriter
	// keySettings map[int][]string // todo: type aliasing for []string

	charts      []*Chart
	list        *scene.List
	listCursors []int
	listDepth   int

	audios.MusicPlayer
}

func NewScene() *Scene {
	s := &Scene{}
	s.charts = newMusics()
	s.listDepth = listDepthMusic
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
func (s *Scene) setMusicPlayer() {
	// Loop: wait + streamer

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

// type FromChooseToPlay struct {
// 	cfg     *Config
// 	asset   *Asset
// 	fsys    fs.FS
// 	name    string
// 	rf      *osr.Format
// }

// choose key bindings from finite selections.
// Todo: KeySettings -> KeyBinding?
//	if inpututil.IsKeyJustPressed(input.KeyF5) {
//		if s.Focus != FocusKeySettings {
//			s.lastFocus = s.Focus
//		}
//		s.keySettings = make([]string, 0)
//		s.Focus = FocusKeySettings
//		scene.UserSkin.Swipe.Play(*s.volumeSound)
//	}
// func setKeySettings() {
// 	for k := input.Key(0); k < input.KeyReserved0; k++ {
// 		if input.IsKeyJustPressed(k) {
// 			name := input.KeyToName(k)
// 			if name[0] == 'F' && name != "F" {
// 				continue
// 			}
// 			s.keySettings = append(s.keySettings, name)
// 		}
// 	}
// 	switch s.mode {
// 	case 0:
// 		if len(s.keySettings) >= 4 {
// 			s.keySettings = s.keySettings[:4]
// 			s.keySettings = mode.NormalizeKeys(s.keySettings)
// 			s.Focus = s.lastFocus
// 		}
// 	}
// }
