package selects

import (
	"image/color"
	"time"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/plays"
	"github.com/hndada/gosu/tween"
	"github.com/hndada/gosu/ui"
)

// 2-depth is enough
const (
	depthFolder = iota
	depthChart
	depthPlay
)

// browser -> folder -> row
// Even though the list component seems too big, it is still a single component.
type chartBrowser struct {
	sprite draws.Sprite // Box sprite for a list item
	game.SearchQuery
	game.SearchResult
	i            int   // outer index
	js           []int // inner indexs, need to keep track of them
	depth        int
	indexHandler ui.KeyNumberHandler[int]
	depthHandler ui.KeyNumberHandler[int]
	tween        tween.Tween // Tween of the list's cursor position
}

const (
	chartRowWidth  = 550
	chartRowHeight = 50
	chartRowCount  = game.ScreenSizeY/chartRowHeight + 1
)

// setXxx and getXxx is not very idiomatic in Go.
// A.loadXxx: create new struct and assign to the member
// -> just use `scn.xxx = newXxx()`, as this looks the most succint way
// func (scn *Scene) addChartTree(boxSprite draws.Sprite, kbs *ui.KeyboardState, r game.SearchResult) (cb chartBrowser) {
func (scn Scene) newChartBrowser(sq game.SearchQuery, sr game.SearchResult) *chartBrowser {
	kbs := scn.KeyboardState
	cb := &chartBrowser{
		sprite:       scn.newChartBoxSprite(),
		SearchQuery:  sq,
		SearchResult: sr,
		// i:            0, // TODO: last chart
	}
	cb.js = make([]int, len(cb.Charts))
	cb.indexHandler = cb.newIndexHandler(kbs)
	cb.depthHandler = cb.newDepthHandler(kbs)
	cb.renewTween()
	return cb
}

func (scn Scene) newChartBoxSprite() draws.Sprite {
	// box realizes row.

	s := draws.NewSprite(scn.Resources.BoxMaskImage)
	s.SetSize(chartRowWidth, chartRowHeight)
	s.Locate(plays.ScreenSizeX/2, plays.ScreenSizeY/2, draws.CenterMiddle)
	return s
}

func (cb *chartBrowser) newIndexHandler(kbs *ui.KeyboardState) ui.KeyNumberHandler[int] {
	ctrls := game.DownUpControls
	ctrls[ui.Decrease].SoundFilename = game.SoundTransitionDown
	ctrls[ui.Increase].SoundFilename = game.SoundTransitionUp

	return ui.KeyNumberHandler[int]{
		NumberController: ui.NumberController[int]{
			Value: &cb.i,
			Min:   0,
			Max:   len(cb.Charts),
			Unit:  1,
		},
		KeyListener: *ui.NewKeyListener(
			kbs,
			[]input.Key{},
			ctrls[:],
		),
	}
}

func (cb *chartBrowser) newDepthHandler(kbs *ui.KeyboardState) ui.KeyNumberHandler[int] {
	ctrls := [2]ui.Control{
		{
			Key:           input.KeyEscape,
			Type:          ui.Decrease,
			SoundFilename: game.SoundTransitionDown,
		},
		{
			Key:           input.KeyEnter,
			Type:          ui.Increase,
			SoundFilename: game.SoundTransitionUp,
		},
	}

	return ui.KeyNumberHandler[int]{
		NumberController: ui.NumberController[int]{
			Value: &cb.depth,
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

func (cb chartBrowser) j() int                { return cb.js[cb.i] }
func (cb chartBrowser) chart() *game.ChartRow { return &cb.Charts[cb.i][cb.j()] }

func (cb *chartBrowser) update() (c *game.ChartRow, isPlay bool) {
	if c, ok := cb.depthHandler.Update(); ok {
		// TODO: refactor
		switch c.Type {
		case ui.Decrease:
			cb.depth--
			if cb.depth < 0 {
				cb.depth = 0
			}
		case ui.Increase:
			cb.depth++
			if cb.depth > depthPlay {
				cb.depth = depthPlay
			}
		}
		cb.renewIndexHandler()
	}

	if cb.depth == depthPlay {
		cb.depth = depthChart
		return cb.chart(), true
	}

	if c, ok := cb.indexHandler.Update(); ok {
		switch c.Type {
		case ui.Increase:
			cb.i--
			if cb.i < 0 {
				cb.i = 0
			}
		case ui.Decrease:
			cb.i++
			if cb.i >= cb.indexHandler.Max {
				cb.i = cb.indexHandler.Max - 1
			}
		}
	}
	cb.renewTween()

	return cb.chart(), false
}

func (cb *chartBrowser) renewIndexHandler() {
	var maxLen int
	var ptr *int
	switch cb.depth {
	case depthFolder:
		ptr = &cb.i
		maxLen = len(cb.Charts)
	case depthChart:
		ptr = &cb.js[cb.i] // cannot use j() here, as j() would return the new value
		maxLen = len(cb.Charts[cb.i])
	}
	cb.indexHandler.Value = ptr
	cb.indexHandler.Max = maxLen - 1
}

func (cb *chartBrowser) renewTween() {
	begin := cb.tween.Value()
	target := chartRowHeight * float64([]int{cb.i, cb.j()}[cb.depth])
	change := target - begin
	cb.tween = tween.Tween{MaxLoop: 1}
	cb.tween.Add(begin, change, 400*time.Millisecond, tween.EaseOutExponential)
	// List is persistent, so no need to check if it is finished.
	cb.tween.Update()
}

// TODO: smoothly enlarge focused item?
func (cb chartBrowser) Draw(dst draws.Image) {
	var (
		skyblue = color.NRGBA{R: 64, G: 255, B: 255, A: 128}
		pink    = color.NRGBA{R: 255, G: 128, B: 255, A: 128}
	)
	var list []draws.Text
	var focusedIndex int
	var maxLen int
	switch cb.depth {
	case depthFolder:
		list = cb.FolderNames
		focusedIndex = cb.i
		maxLen = len(cb.Charts)
	case depthChart:
		list = cb.ChartNames[cb.i]
		focusedIndex = cb.j()
		maxLen = len(cb.Charts[cb.i])
	}

	// List items' positions are fixed; Only the cursor of the list is changed.
	cursor := cb.tween.Value()
	top := cursor - float64(game.ScreenSizeY/2)
	first := int(top/chartRowHeight) - 1
	if first < 0 {
		first = 0
	}

	bottom := cursor + float64(game.ScreenSizeY/2)
	last := int(bottom/chartRowHeight) + 1
	if last >= maxLen {
		last = maxLen - 1
	}

	for i := first; i <= last; i++ {
		pos := chartRowHeight * float64(i)

		// draw box
		box := cb.sprite
		box.Move(0, pos-cursor)
		if i == focusedIndex {
			box.ColorScale.ScaleWithColor(pink)
		} else {
			box.ColorScale.ScaleWithColor(skyblue)
		}
		box.Draw(dst)

		// draw text
		txt := list[i]
		txt.Move(0, pos-cursor)
		txt.Draw(dst)
	}
}
