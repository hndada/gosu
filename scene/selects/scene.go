package selects

import (
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/scene"
)

// TODO: list key handler: double click left/right to open advanced options
// Component is basically EventHandler.
type Scene struct {
	*scene.Resources
	*scene.Options
	*scene.States
	*scene.Handlers
	*scene.Databases

	list               ListComponent
	background         scene.BackgroundComponent
	previewMusicPlayer *PreviewMusicPlayer
	// chartInfo  ChartInfoComponent
	// searchBox  SearchBoxComponent
}

func NewScene(res *scene.Resources, opts *scene.Options,
	states *scene.States, hds *scene.Handlers, dbs *scene.Databases) (*Scene, error) {
	return &Scene{
		Resources: res,
		Options:   opts,
		States:    states,
		Handlers:  hds,
		Databases: dbs,
	}, nil
}

func (s *Scene) Update() any {
	// 1. Listen key input, then update list
	// : ESC, Enter (+NumpadEnter), ArrowUp, ArrowDown, ArrowLeft, ArrowRight
	// Left and Right is for advanced options. (by double click)
	// 2. Update preview music and background
	// 3. render list

	s.Handlers.MusicVolume.Handle()
	s.Handlers.SoundVolumeScale.Handle()
	s.Handlers.MusicOffset.Handle()
	s.Handlers.BackgroundBrightness.Handle()
	s.Handlers.DebugPrint.Handle()

	s.Handlers.Mode.Handle()
	s.Handlers.SubMode.Handle()
	mode := *s.Handlers.Mode.Value
	s.Handlers.SpeedScales[mode].Handle()
	return nil
}

func (s *Scene) playChart() any {
	s.previewMusicPlayer.Close()
	c := s.list.charts[s.list.i][s.list.j]
	return scene.PlayArgs{
		ChartFS:       c.FS,
		ChartFilename: c.Name,
		Mods:          nil,
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

// Avoid embedding scene.Options directly.
// Pass options as pointers for syncing and saving memory.
