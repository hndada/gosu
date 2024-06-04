package selects

import (
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/scene"
)

// Component is basically EventHandler.
type Scene struct {
	*scene.Resources
	*scene.Options
	*scene.Handlers
	*scene.Databases
	query              string
	musicList          []scene.MusicRow
	musicListIndex     int
	chartList          []scene.ChartRow
	chartListIndex     int
	previewMusicPlayer *PreviewMusicPlayer
	background         scene.BackgroundComponent
}

// It is fine to call Close at blank MusicPlayer.

// It is safe to call len() at nil slice.
// https://go.dev/play/p/-1VWc9iDgMl

// Avoid embedding scene.Options directly.
// Pass options as pointers for syncing and saving memory.
func NewScene(res *scene.Resources, opts *scene.Options, hds *scene.Handlers, dbs *scene.Databases) (*Scene, error) {
	return &Scene{
		Resources: res,
		Options:   opts,
		Handlers:  hds,
		Databases: dbs,
	}, nil
}

// music list, chart list, preview music player, background
// update preview and background
// list key handler

// Node is all in choose scene.
func (s *Scene) Update() any {
	oldChart := s.chart()

	key := s.UIKeyListener.Listen()
	if key != input.KeyNone {
		s.SwipeSoundPod.Play(s.SoundVolume)
	}
	switch key {
	case input.KeyEnter, input.KeyNumpadEnter:
		s.chartTreeNode = s.chartTreeNode.FirstChild
	case input.KeyEscape:
		// Todo: return to intro screen when
		// escape is pressed on root node.
		if s.chartTreeNode.Type != FolderNode {
			s.chartTreeNode = s.chartTreeNode.Parent
		}
	// case input.KeyArrowLeft:
	// 	// Arrow left has no effect on root node.
	// 	if s.chartTreeNode.Type == ChartNode {
	// 		s.chartTreeNode = s.chartTreeNode.Parent
	// 	}
	// case input.KeyArrowRight:
	// 	// Arrow right has no effect on leaf node.
	// 	if s.chartTreeNode.Type != LeafNode {
	// 		s.chartTreeNode = s.chartTreeNode.FirstChild
	// 	}
	case input.KeyArrowUp:
		if prev := s.chartTreeNode.Prev(); prev != nil && prev.Type != RootNode {
			s.chartTreeNode = prev
		}
	case input.KeyArrowDown:
		if next := s.chartTreeNode.Next(); next != nil {
			s.chartTreeNode = next
		}
	}

	if s.chartTreeNode.Type == LeafNode {
		return s.playChart()
	}

	newChart := s.chart()
	var (
		isMusicChanged      bool
		isBackgroundChanged bool
	)
	if newChart.Base != oldChart.Base {
		isMusicChanged = true
		isBackgroundChanged = true
	} else {
		if newChart.MusicFilename != oldChart.MusicFilename {
			isMusicChanged = true
		}
		if newChart.BackgroundFilename != oldChart.BackgroundFilename {
			isBackgroundChanged = true
		}
	}

	if isMusicChanged {
		s.updatePreviewMusic()
	}
	if isBackgroundChanged {
		s.updateBackground()
	}

	s.HandleEffect()
	if s.KeyHandleMusicVolume() {
		if !s.MusicPlayer.IsEmpty() {
			s.MusicPlayer.SetVolume(s.MusicVolume)
		}
	}
	s.KeyHandleSoundVolume()
	s.KeyHandleMusicOffset()
	s.KeyHandleBackgroundBrightness()
	s.KeyHandleDebugPrint()
	// s.KeyHandleMode()
	// s.KeyHandleSubMode()
	s.KeyHandleSpeedScale()
	return nil
}

func (s *Scene) updateBackground() {
	c := s.chart()
	fsys := c.MusicFS
	name := c.BackgroundFilename
	s.drawBackground = scene.NewBackgroundDrawer(s.Config, s.Asset, fsys, name)
}

func (s *Scene) playChart() any {
	s.previewMusicPlayer.Close()
	c := s.chartList[s.chartListIndex]

	return scene.PlayArgs{
		ChartFS:       c.FS,
		ChartFilename: c.Name,
		Mods:          nil,
	}
}

// Todo: s.drawSearchBox(screen)
// Todo: s.drawPanel(screen)
func (s *Scene) Draw(dst draws.Image) {
	s.Components.Draw(dst)
}
