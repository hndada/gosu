package choose

import (
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
	// s.handleVolume()
	// scene's handlers
	return nil
}

// fade in/out effect
// osu! seems fading music out when changing music.
// Todo: audios.MusicPlayer.FadeOut(duration time.Duration)
// Todo: audios.MusicPlayer.Duration() time.Duration
func (s *Scene) handleVolume() {
	const wait = -500

	s.Tick++
	if s.Tick == 0 {
		s.Play()
	}
	if s.Tick > 0 && s.Tick <= 1000 {
		age := float64(s.Tick) / 1000
		s.SetVolume(s.Volume * age)
	}

	// when playing osu!'s 10 seconds preview music.
	// need to handle rewinding music anyway.
	if s.Tick > 9000 && s.Tick <= 10000 {
		age := float64(s.Tick-9000) / 1000
		s.SetVolume(s.Volume * (1 - age))
	}

	if s.Tick >= 10000 {
		s.Tick = wait
		s.MusicPlayer.Rewind()
		s.MusicPlayer.Pause()
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
