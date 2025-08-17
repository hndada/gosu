package piano

import (
	"time"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/tween"
)

type JudgmentDrawer struct {
	anims []draws.Animation
	worst game.JudgmentKind
	tween tween.Tween
}

func NewJudgmentDrawer(res *Resources, opts *Options) JudgmentDrawer {
	jd := JudgmentDrawer{}
	jd.anims = make([]draws.Animation, 4)
	for i, frames := range res.JudgmentFramesList {
		a := draws.NewAnimation(frames, 40)
		a.Scale(opts.JudgmentImageScale)
		x := opts.StagePositionX
		y := opts.JudgmentPositionY
		a.Locate(x, y, draws.CenterMiddle)
		jd.anims[i] = a
	}

	tw := tween.Tween{MaxLoop: 1}
	tw.Add(1.00, +0.15, 25*time.Millisecond, tween.EaseLinear)
	tw.Add(1.15, -0.15, 25*time.Millisecond, tween.EaseLinear)
	tw.Add(1.00, +0.0, 200*time.Millisecond, tween.EaseLinear)
	tw.Add(1.00, -0.25, 25*time.Millisecond, tween.EaseLinear)
	tw.Stop() // To make sure jd is invisible at the beginning.
	jd.tween = tw
	return jd
}

func (jd *JudgmentDrawer) Update(keysJudgmentKind []game.JudgmentKind) {
	// worst is guaranteed not to be out of range.
	worst := blank
	for _, jk := range keysJudgmentKind {
		if jk == blank {
			continue
		}
		if worst == blank || worst < jk {
			worst = jk
		}
	}

	if worst <= miss {
		jd.worst = worst
		jd.anims[worst].Reset()
		jd.tween.Start()
	}
	if !jd.tween.IsFinished() {
		jd.tween.Update()
	}
}

func (jd JudgmentDrawer) Draw(dst draws.Image) {
	if jd.tween.IsFinished() {
		return
	}
	a := jd.anims[jd.worst]
	a.Scale(jd.tween.Value())
	a.Draw(dst)
}
