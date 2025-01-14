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
	chartBrowser       *chartBrowser
	background         game.Background
	previewMusicPlayer PreviewMusicPlayer
	// chartInfo  chartInfo

	// Score box color: Gray128 with 50% transparent
	// Hovered Score box color: Gray96 with 50% transparent
	// leaderboard
}

func (Scene) New(g *game.Game, args game.Args) (game.Scene, error) {
	scn := &Scene{Game: g}
	sq := game.SearchQuery{}
	sr := g.Database.Search(sq)
	scn.chartBrowser = scn.newChartBrowser(sq, sr)
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

	lastChart := s.chartBrowser.chart()
	c, isPlay := s.chartBrowser.update()
	if c != nil && isPlay {
		return s.playChart(c)
	}

	if lastChart == nil || lastChart.MusicName != c.MusicName {
		vol := &s.Options.MusicVolume
		s.previewMusicPlayer.Close()
		pmp, err := NewPreviewMusicPlayer(c.FS, c.MusicName, vol)
		if err == nil { // music file may not exist
			s.previewMusicPlayer = pmp
		}
	}
	if lastChart == nil || lastChart.BackgroundFilename != c.BackgroundFilename {
		s.background = game.NewBackground(s.Resources, s.Options)
	}
	return nil
}

func (s Scene) mode() int { return *s.Handlers.Mode.Value }

func (s *Scene) playChart(row *game.ChartRow) any {
	// It is fine to call Close at blank MusicPlayer.
	s.previewMusicPlayer.Close()
	mods := []plays.Mods{piano.Mods{}}[s.mode()]
	return game.PlayArgs{
		ChartFS:       row.FS,
		ChartFilename: row.Name,
		Mods:          mods,
	}
}

func (s Scene) Draw(dst draws.Image) {
	// s.background.Draw(dst)
	s.chartBrowser.Draw(dst)
}

func (s Scene) WindowTitle() string { return "gosu" }
func (s Scene) DebugString() string { return "" }

// It is safe to call len() at nil slice.
// https://go.dev/play/p/-1VWc9iDgMl

// Memo: 'name' is a officially used name as file path in io/fs.

// Memo: make([]T, len) and make([]T, 0, len) is prone to be erroneous.
