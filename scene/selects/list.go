package selects

import (
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/scene"
	"github.com/hndada/gosu/ui"
)

type ListComponent struct {
	sprite      draws.Sprite // Box sprite
	folderTexts []draws.Text
	charts      [][]scene.ChartRow
	i           int // outer index
	j           int // inner index
	depth       int
	lastChart   scene.ChartRow

	indexHandler ui.KeyNumberHandler[int]
	depthHandler ui.KeyNumberHandler[int]
}

func NewListComponent(res *scene.Resources, states *scene.States, folderTexts []draws.Text, charts [][]scene.ChartRow) *ListComponent {
	// TODO: sprite: draws.NewSprite(),
	cmp := &ListComponent{
		folderTexts: folderTexts,
		charts:      charts,
	}
	cmp.indexHandler = newIndexHandler(&cmp.i, len(charts), states)
	cmp.depthHandler = newDepthHandler(&cmp.depth, states)
	return cmp
}

func newIndexHandler(i *int, length int, states *scene.States) ui.KeyNumberHandler[int] {
	ctrls := scene.DownUpControls
	ctrls[ui.Decrease].SoundFilename = scene.SoundTransitionDown
	ctrls[ui.Increase].SoundFilename = scene.SoundTransitionUp

	return ui.KeyNumberHandler[int]{
		NumberController: ui.NumberController[int]{
			Value: i,
			Min:   0,
			Max:   length,
			Unit:  1,
		},
		KeyListener: *ui.NewKeyListener(
			states.Keyboard,
			[]input.Key{},
			ctrls[:],
		),
	}
}

func (cmp *ListComponent) updateIndexHandler() {
	var length int
	var pi *int
	switch cmp.depth {
	case 0:
		pi = &cmp.i
		length = len(cmp.charts)
	case 1:
		pi = &cmp.j
		length = len(cmp.charts[cmp.i])
	}
	cmp.indexHandler.NumberController.Value = pi
	cmp.indexHandler.NumberController.Max = length - 1
}

func newDepthHandler(depth *int, states *scene.States) ui.KeyNumberHandler[int] {
	ctrls := [2]ui.Control{
		{
			Key:           input.KeyEscape,
			Type:          ui.Decrease,
			SoundFilename: scene.SoundTransitionDown,
		},
		{
			Key:           input.KeyEnter,
			Type:          ui.Increase,
			SoundFilename: scene.SoundTransitionUp,
		},
	}

	return ui.KeyNumberHandler[int]{
		NumberController: ui.NumberController[int]{
			Value: depth,
			Min:   0,
			Max:   1,
			Unit:  1,
		},
		KeyListener: *ui.NewKeyListener(
			states.Keyboard,
			[]input.Key{},
			ctrls[:],
		),
	}
}

func (cmp *ListComponent) Update() {
	if _, ok := cmp.depthHandler.Update(); ok {
		cmp.updateIndexHandler()
	}

	lc := cmp.lastChart
	c := cmp.charts[cmp.i][cmp.j]
	if lc.MusicName != c.MusicName {
		// update previewPlayer
	}
	if lc.BackgroundFilename != c.BackgroundFilename {
		// update background
	}
	cmp.lastChart = c
}

// Score box color: Gray128 with 50% transparent
// Hovered Score box color: Gray96 with 50% transparent
func (cmp ListComponent) Draw(dst draws.Image) {
	cmp.sprite.Draw(dst)
}
