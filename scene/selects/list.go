package selects

import (
	"image/color"
	"time"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/scene"
	"github.com/hndada/gosu/tween"
	"github.com/hndada/gosu/ui"
)

const (
	depthFolder = iota
	depthChart
	depthPlay
)

// Even though the list component seems too big, it is still a single component.
type ListComponent struct {
	sprite draws.Sprite // Box sprite for a list item
	// h      float64      // Height of a list item
	tween tween.Tween // Tween of the list's cursor position

	charts       [][]scene.ChartRow
	folderTexts  []draws.Text
	chartTexts   [][]draws.Text
	i            int   // outer index
	js           []int // inner indexs, need to keep track of them
	depth        int
	indexHandler ui.KeyNumberHandler[int]
	depthHandler ui.KeyNumberHandler[int]
}

func NewListComponent(boxSprite draws.Sprite, kbs *ui.KeyboardState, t1 []draws.Text, t2 [][]draws.Text, cs [][]scene.ChartRow) (cmp ListComponent) {
	cmp.sprite = boxSprite
	// cmp.h = boxSprite.H()

	cmp.charts = cs
	cmp.folderTexts = t1
	cmp.chartTexts = t2
	cmp.js = make([]int, len(cs))
	cmp.indexHandler = cmp.newIndexHandler(&cmp.i, len(cs), kbs)
	cmp.depthHandler = cmp.newDepthHandler(&cmp.depth, kbs)
	return cmp
}

func (ListComponent) newIndexHandler(i *int, maxLen int, kbs *ui.KeyboardState) ui.KeyNumberHandler[int] {
	ctrls := scene.DownUpControls
	ctrls[ui.Decrease].SoundFilename = scene.SoundTransitionDown
	ctrls[ui.Increase].SoundFilename = scene.SoundTransitionUp

	return ui.KeyNumberHandler[int]{
		NumberController: ui.NumberController[int]{
			Value: i,
			Min:   0,
			Max:   maxLen,
			Unit:  1,
		},
		KeyListener: *ui.NewKeyListener(
			kbs,
			[]input.Key{},
			ctrls[:],
		),
	}
}

func (ListComponent) newDepthHandler(depth *int, kbs *ui.KeyboardState) ui.KeyNumberHandler[int] {
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
			Max:   depthPlay,
			Unit:  1,
		},
		KeyListener: *ui.NewKeyListener(
			kbs,
			[]input.Key{},
			ctrls[:],
		),
	}
}

func (cmp ListComponent) j() int                 { return cmp.js[cmp.i] }
func (cmp ListComponent) chart() *scene.ChartRow { return &cmp.charts[cmp.i][cmp.j()] }

func (cmp *ListComponent) update() (r *scene.ChartRow, isPlay bool) {
	if _, ok := cmp.depthHandler.Update(); ok {
		cmp.updateIndexHandler()
		cmp.updateTween()
	}
	if cmp.depth == depthPlay {
		cmp.depth = depthChart
		return cmp.chart(), true
	}
	if _, ok := cmp.indexHandler.Update(); ok {
		cmp.updateTween()
	}

	return nil, false
}

func (cmp *ListComponent) updateIndexHandler() {
	var maxLen int
	var ptr *int
	switch cmp.depth {
	case depthFolder:
		ptr = &cmp.i
		maxLen = len(cmp.charts)
	case depthChart:
		ptr = &cmp.js[cmp.i] // cannot use j() here, as j() would return the new value
		maxLen = len(cmp.charts[cmp.i])
	}
	cmp.indexHandler.Value = ptr
	cmp.indexHandler.Max = maxLen - 1
}

func (cmp *ListComponent) updateTween() {
	begin := cmp.tween.Value()
	target := listBoxHeight * float64([]int{cmp.i, cmp.j()}[cmp.depth])
	change := target - begin
	cmp.tween = tween.Tween{MaxLoop: 1}
	cmp.tween.Add(begin, change, 400*time.Millisecond, tween.EaseOutExponential)
	// List is persistent, so no need to check if it is finished.
	cmp.tween.Update()
}

// TODO: smoothly enlarge focused item?
func (cmp ListComponent) Draw(dst draws.Image) {
	var list []draws.Text
	var idx int
	var maxLen int
	var clr = color.NRGBA{R: 64, G: 255, B: 255, A: 128} // skyblue
	switch cmp.depth {
	case depthFolder:
		list = cmp.folderTexts
		idx = cmp.i
		maxLen = len(cmp.charts)
	case depthChart:
		list = cmp.chartTexts[cmp.i]
		idx = cmp.j()
		maxLen = len(cmp.charts[cmp.i])
	}

	// List items' positions are fixed; Only the cursor of the list is changed.
	cursor := cmp.tween.Value()
	top := cursor - float64(scene.ScreenSizeY/2)
	first := int(top/listBoxHeight) - 1
	if first < 0 {
		first = 0
	}

	bottom := cursor + float64(scene.ScreenSizeY/2)
	last := int(bottom/listBoxHeight) + 1
	if last >= maxLen {
		last = maxLen - 1
	}

	for i := first; i <= last; i++ {
		pos := listBoxHeight * float64(i)

		// draw box
		box := cmp.sprite
		box.Move(0, pos-cursor)
		if i == idx {
			clr = color.NRGBA{R: 255, G: 128, B: 255, A: 128} // pink
		}
		box.ColorScale.ScaleWithColor(clr)
		box.Draw(dst)

		// draw text
		txt := list[i]
		txt.Move(0, pos-cursor)
		txt.Draw(dst)
	}
}
