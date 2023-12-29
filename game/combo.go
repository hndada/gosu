package game

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/draws"
)

type ComboResources struct {
	imgs [10]draws.Image
}

func (res *ComboResources) Load(fsys fs.FS) {
	for i := 0; i < 10; i++ {
		fname := fmt.Sprintf("combo/%d.png", i)
		res.imgs[i] = draws.NewImageFromFile(fsys, fname)
	}
}

type ComboOptions struct {
	Scale    float64
	X        float64
	Y        float64
	DigitGap float64
	Bounce   float64
	Persist  bool
}

type ComboComponent struct {
	Combo     *int
	lastCombo int     // to reset tween
	w         float64 // fixed width
	sprites   [10]draws.Sprite
	tween     draws.Tween
}

func NewComboComponent(res ComboResources, opts ComboOptions) (cmp ComboComponent) {
	for i := 0; i < 10; i++ {
		sprite := draws.NewSprite(res.imgs[i])
		sprite.MultiplyScale(opts.Scale)
		sprite.Locate(opts.X, opts.Y, draws.CenterMiddle)
		cmp.sprites[i] = sprite
	}
	// Size of the whole image is 0.5w + (n-1)(w+gap) + 0.5w.
	// Since sprites are already at anchor, no need to care of two 0.5w.
	cmp.w = cmp.sprites[0].W() + opts.DigitGap

	tw := draws.Tween{}
	tw.Add(0, opts.Bounce, 200, draws.EaseLinear)
	tw.Add(opts.Bounce, -opts.Bounce, 100, draws.EaseLinear)
	tw.Add(0, 0, 1500, draws.EaseLinear)
	if !opts.Persist {
		tw.SetLoop(1)
	}
	return
}

func (cmp *ComboComponent) Update() {
	if cmp.lastCombo != *cmp.Combo {
		cmp.lastCombo = *cmp.Combo
		cmp.tween.Reset()
	}
}

// Each number has different width. Number 0's width is used as standard.
// ComboDrawer's Draw draws each number at constant x regardless of their widths.
func (cmp ComboComponent) Draw(dst draws.Image) {
	if cmp.tween.IsFinished() {
		return
	}
	if *cmp.Combo == 0 {
		return
	}

	vs := make([]int, 0)
	for v := *cmp.Combo; v > 0; v /= 10 {
		vs = append(vs, v%10) // Little endian.
	}

	tx := float64(len(vs)-1) * cmp.w / 2
	for _, v := range vs {
		s := cmp.sprites[v]
		ty := cmp.tween.Current() * s.H()
		s.Move(tx, ty)
		s.Draw(dst)
		tx -= cmp.w
	}
}
