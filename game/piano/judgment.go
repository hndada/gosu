package piano

import (
	"time"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/tween"
)

type JudgmentComponent struct {
	anims []draws.Animation
	worst game.JudgmentKind
	tween tween.Tween
}

func NewJudgmentComponent(res *Resources, opts *Options) (cmp JudgmentComponent) {
	cmp.anims = make([]draws.Animation, 4)
	for i, frames := range res.JudgmentFramesList {
		a := draws.NewAnimation(frames, 40)
		a.Scale(opts.JudgmentImageScale)
		x := opts.StagePositionX
		y := opts.JudgmentPositionY
		a.Locate(x, y, draws.CenterMiddle)
		cmp.anims[i] = a
	}

	tw := tween.Tween{}
	tw.Add(1.00, +0.15, 25*time.Millisecond, tween.EaseLinear)
	tw.Add(1.15, -0.15, 25*time.Millisecond, tween.EaseLinear)
	tw.Add(1.00, +0.0, 200*time.Millisecond, tween.EaseLinear)
	tw.Add(1.00, -0.25, 25*time.Millisecond, tween.EaseLinear)
	tw.Stop() // To make sure judgment is invisible at the beginning.
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
		cmp.tween.Start()
	}
}

func (cmp JudgmentComponent) Draw(dst draws.Image) {
	if cmp.tween.IsFinished() {
		return
	}
	a := cmp.anims[cmp.worst]
	a.Scale(cmp.tween.Value())
	a.Draw(dst)
}
