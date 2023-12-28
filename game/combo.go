package game

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/draws"
)

type ComboRes struct {
	imgs [10]draws.Image
}

func (res *ComboRes) Load(fsys fs.FS) {
	for i := 0; i < 10; i++ {
		fname := fmt.Sprintf("combo/%d.png", i)
		res.imgs[i] = draws.NewImageFromFile(fsys, fname)
	}
}

type ComboOpts struct {
	Scale    float64
	X        float64
	Y        float64
	DigitGap float64
	Bounce   float64
	Persist  bool
}

type ComboComp struct {
	Combo     *int
	lastCombo int     // to reset tween
	w         float64 // fixed width
	sprites   [10]draws.Sprite
	tween     draws.Tween
}

func NewComboComp(res ComboRes, opts ComboOpts) (comp ComboComp) {
	for i := 0; i < 10; i++ {
		sprite := draws.NewSprite(res.imgs[i])
		sprite.MultiplyScale(opts.Scale)
		sprite.Locate(opts.X, opts.Y, draws.CenterMiddle)
		comp.sprites[i] = sprite
	}
	// Size of the whole image is 0.5w + (n-1)(w+gap) + 0.5w.
	// Since sprites are already at anchor, no need to care of two 0.5w.
	comp.w = comp.sprites[0].W() + opts.DigitGap

	tw := draws.Tween{}
	tw.Add(0, opts.Bounce, 200, draws.EaseLinear)
	tw.Add(opts.Bounce, -opts.Bounce, 100, draws.EaseLinear)
	tw.Add(0, 0, 1500, draws.EaseLinear)
	if !opts.Persist {
		tw.SetLoop(1)
	}
	return
}

func (comp *ComboComp) Update() {
	if comp.lastCombo != *comp.Combo {
		comp.lastCombo = *comp.Combo
		comp.tween.Reset()
	}
}

// Each number has different width. Number 0's width is used as standard.
// ComboDrawer's Draw draws each number at constant x regardless of their widths.
func (comp ComboComp) Draw(dst draws.Image) {
	if comp.tween.IsFinished() {
		return
	}
	if *comp.Combo == 0 {
		return
	}

	vs := make([]int, 0)
	for v := *comp.Combo; v > 0; v /= 10 {
		vs = append(vs, v%10) // Little endian.
	}

	tx := float64(len(vs)-1) * comp.w / 2
	for _, v := range vs {
		s := comp.sprites[v]
		ty := comp.tween.Current() * s.H()
		s.Move(tx, ty)
		s.Draw(dst)
		tx -= comp.w
	}
}
