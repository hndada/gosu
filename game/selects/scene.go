package selects

import (
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/plays"
	"github.com/hndada/gosu/plays/piano"
)

// TODO: list key handler: double click left/right to open advanced options
// Component is basically EventHandler.
type Scene struct {
	*game.Game
	boxSprite          draws.Sprite
	list               ListComponent
	lastChart          *game.ChartRow
	background         game.BackgroundComponent
	volume             *float64
	previewMusicPlayer PreviewMusicPlayer
	// chartInfo  ChartInfoComponent
	// searchBox  SearchBoxComponent

	// Score box color: Gray128 with 50% transparent
	// Hovered Score box color: Gray96 with 50% transparent
	// leaderboard
}

const (
	listBoxWidth  = 400
	listBoxHeight = 100
	listBoxCount  = game.ScreenSizeY/listBoxHeight + 1
)

func NewScene(g *game.Game) (*Scene, error) {
	scn := &Scene{Game: g}
	s := draws.NewSprite(g.Resources.BoxMaskImage)
	s.SetSize(listBoxWidth, listBoxHeight)
	s.Locate(plays.ScreenSizeX/2, plays.ScreenSizeY/2, draws.CenterMiddle)
	scn.boxSprite = s

	// cmp.lastChart = charts[0][0]
	return scn, nil
}

// (+NumpadEnter)
// Left and Right arrows are for advanced options. (by double click)
func (s *Scene) Update() any {
	s.Handlers.MusicVolume.Handle()
	s.Handlers.SoundVolumeScale.Handle()
	s.Handlers.MusicOffset.Handle()
	s.Handlers.BackgroundBrightness.Handle()
	s.Handlers.DebugPrint.Handle()

	s.Handlers.Mode.Handle()
	s.Handlers.SubMode.Handle()
	s.Handlers.SpeedScales[s.mode()].Handle()

	c, isPlay := s.list.update()
	if c != nil && isPlay {
		return s.playChart(c)
	}

	lc := s.lastChart
	if lc == nil || lc.MusicName != c.MusicName {
		s.previewMusicPlayer.Close()
		pmp, err := NewPreviewMusicPlayer(c.FS, c.MusicName, s.volume)
		if err == nil { // music file may not exist
			s.previewMusicPlayer = pmp
		}
	}
	if lc == nil || lc.BackgroundFilename != c.BackgroundFilename {
		s.background = game.NewBackgroundComponent(s.Resources, s.Options)
	}
	s.lastChart = c
	return nil
}

func (s Scene) mode() int { return *s.Handlers.Mode.Value }

func (s *Scene) playChart(row *game.ChartRow) any {
	s.previewMusicPlayer.Close()
	mods := []plays.Mods{piano.Mods{}}[s.mode()]
	return game.PlayArgs{
		ChartFS:       row.FS,
		ChartFilename: row.Name,
		Mods:          mods,
	}
}

func (s Scene) Draw(dst draws.Image) {
	s.background.Draw(dst)
	s.list.Draw(dst)
}

func (s Scene) DebugString() string {
	return ""
}

// It is fine to call Close at blank MusicPlayer.

// It is safe to call len() at nil slice.
// https://go.dev/play/p/-1VWc9iDgMl

// Memo: 'name' is a officially used name as file path in io/fs.

// Memo: make([]T, len) and make([]T, 0, len) is prone to be erroneous.

// Avoid embedding game.Options directly.
// Pass options as pointers for syncing and saving memory.
