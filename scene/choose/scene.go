package choose

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/mode/drum"
)

const (
	cursorDepthSearch = -1
	cursorDepthMusic  = iota
	cursorDepthChart
	// FocusKeySettings
)

// keySettings: choose from finite selections.
const modeAll = -1

// 'name' is a officially used name as file path in io/fs.

// Todo: implement non-playing score simulator
// osu! seems fading music out when changing music.

// save the scene to gosu.scenes

// Background brightness at Song select: 60% (153 / 255), confirmed.
// Score box color: Gray128 with 50% transparent
// Hovered Score box color: Gray96 with 50% transparent

// Todo: fetch Score with Replay
// Group1, Group2, Sort, Filter int
// defer sort
type Scene struct {
	musics          []Music
	charts          map[[16]byte]*Chart
	cursorDepth     int // Focus
	mode            int
	subMode         int
	queryTypeWriter input.TypeWriter
	cursors         []int
	audios.MusicPlayer
	keySettings map[int][]string // todo: type aliasing for []string
}

func NewScene() *Scene {
	s := &Scene{}
	s.musics = newMusics()
	s.cursorDepth = cursorDepthMusic
	// s.handleEnter()
	// ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)
	// debug.SetGCPercent(100)
	ebiten.SetWindowTitle("gosu")
	return s
}

func (s *Scene) Update() any {
	// sort by. // music name, level folder (+time?)
	modeKeyHandler() // IsKeyJustPressed? // then I cannot play music.
	if cursorKeyHandler() {
		// play preview music if possible
	}
	isEnter() // cursorDepth++;
	isBack()  // cursorDepth--; 0
	// select music
	// select chart
	return nil
}

func (s Scene) newLists() {}

func (s *Scene) setModeKeyHandler() {
	// s.Mode = ctrl.KeyHandler{
	// 	Handler: ctrl.IntHandler{
	// 		Value: &s.mode,
	// 		Min:   0,
	// 		Max:   3 - 1, // There are 3 modes.
	// 		Loop:  true,
	// 	},
	// 	Modifiers: []input.Key{},
	// 	Keys:      [2]input.Key{-1, input.KeyF1},
	// 	Sounds:    [2]audios.Sounder{scene.UserSkin.Swipe, scene.UserSkin.Swipe},
	// 	Volume:    &mode.S.SoundVolume,
	// }

	// if inpututil.IsKeyJustPressed(input.KeyF1) {
	// 	s.mode++
	// 	s.mode %= 3
	// 	scene.UserSkin.Enter.Play(*s.volumeSound)
	// 	s.Focus = FocusChartSet
	// 	// s.Focus = FocusSearch
	// 	// err := s.handleEnter()
	// 	// if err != nil {
	// 	// 	fmt.Println(err)
	// 	// }
	// }
}

//	if inpututil.IsKeyJustPressed(input.KeyF5) {
//		if s.Focus != FocusKeySettings {
//			s.lastFocus = s.Focus
//		}
//		s.keySettings = make([]string, 0)
//		s.Focus = FocusKeySettings
//		scene.UserSkin.Swipe.Play(*s.volumeSound)
//	}

func (s *Scene) changeMusic() {

	s.MusicPlayer, err := audios.NewMusicPlayerFromFile(fsys, s.MusicFilename, ratio)
	if err != nil {
		return
	}
}

// fade in/out effect
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

// Todo: KeySettings -> KeyBinding?
func setKeySettings() {
	for k := input.Key(0); k < input.KeyReserved0; k++ {
		if inpututil.IsKeyJustPressed(k) {
			name := input.KeyToName(k)
			if name[0] == 'F' && name != "F" {
				continue
			}
			s.keySettings = append(s.keySettings, name)
		}
	}
	switch s.mode {
	case 0:
		if len(s.keySettings) >= 4 {
			s.keySettings = s.keySettings[:4]
			s.keySettings = mode.NormalizeKeys(s.keySettings)
			piano.S.KeySettings[4] = s.keySettings
			s.Focus = s.lastFocus
		}
	case 1:
		if len(s.keySettings) >= 7 {
			s.keySettings = s.keySettings[:7]
			s.keySettings = mode.NormalizeKeys(s.keySettings)
			piano.S.KeySettings[7] = s.keySettings
			s.Focus = s.lastFocus
		}
	case 2:
		if len(s.keySettings) >= 4 {
			s.keySettings = s.keySettings[:4]
			s.keySettings = mode.NormalizeKeys(s.keySettings)
			drum.S.KeySettings[4] = s.keySettings
			s.Focus = s.lastFocus
		}
	}
}

func (s Scene) DebugString() string {

}
