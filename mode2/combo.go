package mode

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/draws"
)

type ComboRes struct {
	imgs [10]draws.Image
}

func (cr *ComboRes) Load(fsys fs.FS) {
	for i := 0; i < 10; i++ {
		fname := fmt.Sprintf("combo/%d.png", i)
		cr.imgs[i] = draws.NewImageFromFile(fsys, fname)
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
	Combo     int
	lastCombo int     // to reset tween
	w         float64 // fixed width
	sprites   [10]draws.Sprite
	tween     draws.Tween
}

func NewComboComp(res ComboRes, opts ComboOpts) (cc ComboComp) {
	for i := 0; i < 10; i++ {
		sprite := draws.NewSprite(res.imgs[i])
		sprite.MultiplyScale(opts.Scale)
		sprite.Locate(opts.X, opts.Y, draws.CenterMiddle)
		cc.sprites[i] = sprite
	}
	// Size of the whole image is 0.5w + (n-1)(w+gap) + 0.5w.
	// Since sprites are already at anchor, no need to care of two 0.5w.
	cc.w = cc.sprites[0].W() + opts.DigitGap

	tw := draws.Tween{}
	tw.Add(0, opts.Bounce, 200, draws.EaseLinear)
	tw.Add(opts.Bounce, -opts.Bounce, 100, draws.EaseLinear)
	tw.Add(0, 0, 1500, draws.EaseLinear)
	if !opts.Persist {
		tw.SetLoop(1)
	}
	return
}

func (cc *ComboComp) Update() {
	if cc.lastCombo != cc.Combo {
		cc.lastCombo = cc.Combo
		cc.tween.Reset()
	}
}

// Each number has different width. Number 0's width is used as standard.
// ComboDrawer's Draw draws each number at constant x regardless of their widths.
func (cc ComboComp) Draw(dst draws.Image) {
	if cc.tween.IsFinished() {
		return
	}
	if cc.Combo == 0 {
		return
	}

	vs := make([]int, 0)
	for v := cc.Combo; v > 0; v /= 10 {
		vs = append(vs, v%10) // Little endian.
	}

	tx := float64(len(vs)-1) * cc.w / 2
	for _, v := range vs {
		s := cc.sprites[v]
		ty := cc.tween.Current() * s.H()
		s.Move(tx, ty)
		s.Draw(dst)
		tx -= cc.w
	}
}
