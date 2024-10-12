package game

import (
	"fmt"
	"io/fs"
	"time"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/tween"
)

func LoadComboImages(fsys fs.FS) []draws.Image {
	imgs := make([]draws.Image, 10)
	for i := 0; i < 10; i++ {
		fname := fmt.Sprintf("combo/%d.png", i)
		imgs[i] = draws.NewImageFromFile(fsys, fname)
	}
	return imgs
}

type ComboOptions struct {
	ImageScale float64
	PositionX  float64
	DigitGap   float64
	PositionY  float64
	IsPersist  bool
	Bounce     float64
}

type ComboComponent struct {
	sprites []draws.Sprite
	combo   int
	w       float64
	tween   tween.Tween
}

func NewComboComponent(imgs []draws.Image, opts *ComboOptions) (cmp ComboComponent) {
	cmp.sprites = make([]draws.Sprite, 10)
	for i := 0; i < 10; i++ {
		sprite := draws.NewSprite(imgs[i])
		sprite.Scale(opts.ImageScale)
		sprite.Locate(opts.PositionX, opts.PositionY, draws.CenterMiddle)
		cmp.sprites[i] = sprite
	}
	// Size of the whole image is 0.5w + (n-1)(w+gap) + 0.5w.
	// Since sprites are already at anchor, no need to care of two 0.5w.
	cmp.w = cmp.sprites[0].W() + opts.DigitGap
	tw := tween.Tween{}
	b := opts.Bounce
	tw.Add(0, b, 150*time.Millisecond, tween.EaseLinear)
	tw.Add(b, -b, 100*time.Millisecond, tween.EaseLinear)
	tw.Add(0, 0, 1500*time.Millisecond, tween.EaseLinear)
	if !opts.IsPersist {
		tw.MaxLoop = 1
	}
	cmp.tween = tw
	return
}

func (cmp *ComboComponent) Update(newCombo int) {
	if old := cmp.combo; old != newCombo {
		cmp.combo = newCombo
		cmp.tween.Start()
	}

	if !cmp.tween.IsFinished() {
		cmp.tween.Update()
	}
}

// Each number has different width. Number 0's width is used as standard.
// ComboDrawer's Draw draws each number at constant x regardless of their widths.
func (cmp ComboComponent) Draw(dst draws.Image) {
	if cmp.tween.IsFinished() {
		return
	}
	if cmp.combo == 0 {
		return
	}

	vs := make([]int, 0)
	for v := cmp.combo; v > 0; v /= 10 {
		vs = append(vs, v%10) // Little endian.
	}

	tx := float64(len(vs)-1) * cmp.w / 2
	for _, v := range vs {
		s := cmp.sprites[v]
		ty := cmp.tween.Value() * s.H()
		s.Move(tx, ty)
		s.Draw(dst)
		tx -= cmp.w
	}
}
