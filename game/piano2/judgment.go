package piano

import (
	"fmt"
	"io/fs"

	draws "github.com/hndada/gosu/draws5"
	"github.com/hndada/gosu/game"
)

type JudgmentResources struct {
	framesList []draws.Frames
}

func (res *JudgmentResources) Load(fsys fs.FS) {
	res.framesList = make([]draws.Frames, 4)
	for i, name := range []string{"kool", "cool", "good", "miss"} {
		fname := fmt.Sprintf("piano/judgment/%s.png", name)
		res.framesList[i] = draws.NewFramesFromFile(fsys, fname)
	}
}

type JudgmentOptions struct {
	Scale float64
	x     float64
	Y     float64
}

func NewJudgmentOptions(stage StageOptions) JudgmentOptions {
	return JudgmentOptions{
		Scale: 0.33,
		x:     stage.w,
		Y:     0.66 * game.ScreenH,
	}
}

type JudgmentComponent struct {
	anims []draws.Animation
	worst game.JudgmentKind
	tween draws.Tween
}

func NewJudgmentComponent(res JudgmentResources, opts JudgmentOptions) (cmp JudgmentComponent) {
	cmp.anims = make([]draws.Animation, 4)
	for i, frames := range res.framesList {
		a := draws.NewAnimation(frames, 40)
		a.MultiplyScale(opts.Scale)
		a.Locate(opts.x, opts.Y, draws.CenterMiddle)
		cmp.anims[i] = a
	}

	tw := draws.Tween{}
	tw.Add(1.00, +0.15, 25, draws.EaseLinear)
	tw.Add(1.15, -0.15, 25, draws.EaseLinear)
	tw.Add(1.00, +0.0, 200, draws.EaseLinear)
	tw.Add(1.00, -0.25, 25, draws.EaseLinear)
	tw.Finish() // To avoid drawing at the beginning.
	cmp.tween = tw
	return
}

func (cmp *JudgmentComponent) Update(kjk []game.JudgmentKind) {
	// worst is guaranteed not to be out of range.
	worst := blank
	for _, jk := range kjk {
		if worst == blank || worst < jk {
			worst = jk
		}
	}
	if worst <= miss {
		cmp.worst = worst
		cmp.anims[worst].Reset()
		cmp.tween.Reset()
	}
}

func (cmp JudgmentComponent) Draw(dst draws.Image) {
	if cmp.tween.IsFinished() {
		return
	}
	a := cmp.anims[cmp.worst]
	a.MultiplyScale(cmp.tween.Current())
	a.Draw(dst)
}
