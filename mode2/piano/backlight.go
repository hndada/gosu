package piano

import (
	"image/color"
	"io/fs"

	"github.com/hndada/gosu/draws"
)

type BacklightsRes struct {
	img draws.Image
}

func (br *BacklightsRes) Load(fsys fs.FS) {
	br.img = draws.NewImageFromFile(fsys, "piano/key/backlight.png")
}

type BacklightsOpts struct {
	ws     []float64
	xs     []float64
	y      float64
	Colors [4]color.NRGBA
}

func NewBacklightsOpts(keys KeysOpts) BacklightsOpts {
	return BacklightsOpts{
		ws: keys.ws,
		xs: keys.xs,
		y:  keys.BaselineY,
		Colors: [4]color.NRGBA{
			{224, 224, 224, 64}, // One: white
			{255, 170, 204, 64}, // Two: pink
			{224, 224, 0, 64},   // Mid: yellow
			{255, 0, 0, 64},     // Tip: red
		},
	}
}

type BacklightsComp struct {
	sprites []draws.Sprite
}

func NewBacklightsComp(res BacklightsRes, opts BacklightsOpts) (comp BacklightsComp) {
	keyCount := len(opts.ws)
	comp.sprites = make([]draws.Sprite, keyCount)
	for k := range comp.sprites {
		s := draws.NewSprite(res.img)
		s.MultiplyScale(opts.ws[k] / s.W())
		s.Locate(opts.xs[k], opts.y, draws.CenterBottom) // -HintHeight
		comp.sprites[k] = s
	}
	return
}
