package piano

import (
	"image/color"
	"io/fs"
	"time"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/times"
)

type BacklightsRes struct {
	img draws.Image
}

func (br *BacklightsRes) Load(fsys fs.FS) {
	fname := "piano/key/backlight.png"
	br.img = draws.NewImageFromFile(fsys, fname)
}

type BacklightsOpts struct {
	keyCount int
	kw       []float64
	kx       []float64
	y        float64 // center bottom
	order    []KeyKind
	Colors   [4]color.NRGBA
}

func NewBacklightsOpts(keys KeysOpts) BacklightsOpts {
	return BacklightsOpts{
		keyCount: keys.keyCount,
		kw:       keys.kw,
		kx:       keys.kx,
		y:        keys.y,
		order:    keys.Order(),
		Colors: [4]color.NRGBA{
			{224, 224, 224, 64}, // One: white
			{255, 170, 204, 64}, // Two: pink
			{224, 224, 0, 64},   // Mid: yellow
			{255, 0, 0, 64},     // Tip: red
		},
	}
}

type BacklightsComp struct {
	keysPressed []bool
	sprites     []draws.Sprite
	startTimes  []time.Time
	minDuration time.Duration
}

func NewBacklightsComp(res BacklightsRes, opts BacklightsOpts) (comp BacklightsComp) {
	comp.sprites = make([]draws.Sprite, opts.keyCount)
	for k := range comp.sprites {
		s := draws.NewSprite(res.img)
		s.MultiplyScale(opts.kw[k] / s.W())
		s.Locate(opts.kx[k], opts.y, draws.CenterBottom)
		s.ColorScale.ScaleWithColor(opts.Colors[opts.order[k]])
		comp.sprites[k] = s
	}
	comp.startTimes = make([]time.Time, opts.keyCount)
	comp.minDuration = 30 * time.Millisecond
	return
}

func (comp *BacklightsComp) Update(kp []bool) {
	comp.keysPressed = kp
	for k, down := range kp {
		if down {
			comp.startTimes[k] = times.Now()
		}
	}
}

// Draw backlights for a while even if the press is brief.
func (comp BacklightsComp) Draw(dst draws.Image) {
	elapsed := times.Since(comp.startTimes[0])
	for k, p := range comp.keysPressed {
		if p || elapsed <= comp.minDuration {
			comp.sprites[k].Draw(dst)
		} else {
			comp.sprites[k].Draw(dst)
		}
	}
}
