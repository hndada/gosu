package choose

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/scene"
)

type Scene struct {
	*scene.Config
	*scene.Asset

	charts  map[string]*Chart // key: chart.Hash
	replays map[string]*osr.Format

	ctrl.UIKeyListener
	chartTreeNode *Node // focused chart
	PreviewMusicPlayer
	drawBackground func(draws.Image)

	KeyHandleMusicVolume          func() bool
	KeyHandleSoundVolume          func() bool
	KeyHandleMusicOffset          func() bool
	KeyHandleBackgroundBrightness func() bool
	KeyHandleDebugPrint           func() bool

	// Todo: key handle mode and sub mode.
	// Todo: add 'NoRepeat' option to KeyHandle.
	// KeyHandleMode       func() bool
	// KeyHandleSubMode    func() bool
	KeyHandleSpeedScale func() bool

	// queryTypeWriter input.TypeWriter
	// keySettings map[int][]string // todo: type aliasing for []string
}

func NewScene(cfg *scene.Config, asset *scene.Asset, root fs.FS) (s *Scene, err error) {
	s = &Scene{Config: cfg, Asset: asset}

	musicRoot, err := fs.Sub(root, cfg.MusicRoots[0])
	if err != nil {
		return
	}
	var errsParse []error
	s.charts, errsParse = newCharts(musicRoot)
	for _, err := range errsParse {
		fmt.Println(err)
	}

	replayRoot, err := fs.Sub(root, "replays")
	if err != nil {
		return
	}
	s.replays = newReplays(replayRoot, s.charts)

	uiKeys := []input.Key{
		input.KeyEnter, input.KeyNumpadEnter, input.KeyEscape,
		input.KeyArrowUp, input.KeyArrowDown,
		input.KeyArrowLeft, input.KeyArrowRight,
	}
	s.UIKeyListener = ctrl.NewUIKeyListener(uiKeys)
	s.chartTreeNode = newChartTree(s.charts).FirstChild
	s.updatePreviewMusic()
	s.updateBackground()

	s.KeyHandleMusicVolume = scene.NewMusicVolumeKeyHandler(s.Config, s.Asset)
	s.KeyHandleSoundVolume = scene.NewSoundVolumeKeyHandler(s.Config, s.Asset)
	s.KeyHandleMusicOffset = scene.NewMusicOffsetKeyHandler(s.Config, s.Asset)
	s.KeyHandleBackgroundBrightness = scene.NewBackgroundBrightnessKeyHandler(s.Config, s.Asset)
	s.KeyHandleDebugPrint = scene.NewDebugPrintKeyHandler(s.Config, s.Asset)

	// s.KeyHandleMode = scene.NewModeKeyHandler(s.Config, s.Asset)
	// s.KeyHandleSubMode = scene.NewSubModeKeyHandler(s.Config, s.Asset, s.Mode)
	s.KeyHandleSpeedScale = scene.NewSpeedScaleKeyHandler(s.Config, s.Asset, s.Mode)
	return
}

// Node is all in choose scene.
func (s *Scene) Update() any {
	oldChart := s.chart()

	key := s.UIKeyListener.Listen()
	switch key {
	case input.KeyEnter, input.KeyNumpadEnter:
		s.chartTreeNode = s.chartTreeNode.FirstChild
		fmt.Printf("%+v\n", s.chartTreeNode)
	case input.KeyEscape:
		// Todo: return to intro screen when
		// escape is pressed on root node.
		if s.chartTreeNode.Type != FolderNode {
			s.chartTreeNode = s.chartTreeNode.Parent
		}
	case input.KeyArrowLeft:
		// Arrow left has no effect on root node.
		if s.chartTreeNode.Type == ChartNode {
			s.chartTreeNode = s.chartTreeNode.Parent
		}
	case input.KeyArrowRight:
		// Arrow right has no effect on leaf node.
		if s.chartTreeNode.Type != LeafNode {
			s.chartTreeNode = s.chartTreeNode.FirstChild
		}
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
		s.MusicPlayer.SetVolume(s.MusicVolume)
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

func (s Scene) chart() *Chart { return s.charts[s.chartTreeNode.LeafData()] }

func (s *Scene) playChart() any {
	s.MusicPlayer.Close()
	c := s.chart()
	s.chartTreeNode = s.chartTreeNode.Parent
	return scene.PlayArgs{
		MusicFS:       c.MusicFS,
		ChartFilename: c.Filename,
		Replay:        s.replays[c.Hash],
	}
}

func (s *Scene) updateBackground() {
	c := s.chart()
	fsys := c.MusicFS
	name := c.BackgroundFilename
	s.drawBackground = scene.NewBackgroundDrawer(s.Config, s.Asset, fsys, name)
}

func (s Scene) DebugString() string {
	return ""
}
