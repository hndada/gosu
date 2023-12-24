package piano

import (
	"image/color"
	"io/fs"
	"time"

	"github.com/hndada/gosu/draws"
)

type BacklightsRes struct {
	img draws.Image
}

func (br *BacklightsRes) Load(fsys fs.FS) {
	fname := "piano/key/backlight.png"
	br.img = draws.NewImageFromFile(fsys, fname)
}

type BacklightsOpts struct {
	ws          []float64
	xs          []float64
	y           float64
	order       []KeyKind
	Colors      [4]color.NRGBA
	minDuration int32 // milliseconds
}

func NewBacklightsOpts(keys KeysOpts) BacklightsOpts {
	return BacklightsOpts{
		ws:    keys.ws,
		xs:    keys.xs,
		y:     keys.BaselineY,
		order: keys.Order(),
		Colors: [4]color.NRGBA{
			{224, 224, 224, 64}, // One: white
			{255, 170, 204, 64}, // Two: pink
			{224, 224, 0, 64},   // Mid: yellow
			{255, 0, 0, 64},     // Tip: red
		},
		minDuration: 30,
	}
}

type BacklightsComp struct {
	keyDowns    []bool
	sprites     []draws.Sprite
	startTimes  []time.Time
	minDuration int32
}

func NewBacklightsComp(res BacklightsRes, opts BacklightsOpts) (comp BacklightsComp) {
	keyCount := len(opts.ws)
	comp.sprites = make([]draws.Sprite, keyCount)
	for k := range comp.sprites {
		s := draws.NewSprite(res.img)
		s.MultiplyScale(opts.ws[k] / s.W())
		s.Locate(opts.xs[k], opts.y, draws.CenterBottom) // -HintHeight
		s.ColorScale.ScaleWithColor(opts.Colors[opts.order[k]])
		comp.sprites[k] = s
	}
	comp.startTimes = make([]time.Time, keyCount)
	comp.minDuration = opts.minDuration
	return
}

func (comp *BacklightsComp) Update(keyDowns []bool) {
	comp.keyDowns = keyDowns
	for k, down := range keyDowns {
		if down {
			comp.startTimes[k] = time.Now()
		}
	}
}

// Draw backlights for a while even if the press is brief.
func (comp BacklightsComp) Draw(dst draws.Image) {
	elapsed := time.Since(comp.startTimes[0]).Milliseconds()
	for k, keyDown := range comp.keyDowns {
		if keyDown || int32(elapsed) <= comp.minDuration {
			comp.sprites[k].Draw(dst)
		} else {
			comp.sprites[k].Draw(dst)
		}
	}
}
