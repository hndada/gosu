package piano

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
)

type JudgmentRes struct {
	framesList [4]draws.Frames
}

func (res *JudgmentRes) Load(fsys fs.FS) {
	for i, name := range []string{"kool", "cool", "good", "miss"} {
		fname := fmt.Sprintf("piano/judgment/%s.png", name)
		res.framesList[i] = draws.NewFramesFromFile(fsys, fname)
	}
	return
}

type JudgmentOpts struct {
	Scale float64
	x     float64
	Y     float64
}

func NewJudgmentOpts(keys KeysOpts) JudgmentOpts {
	return JudgmentOpts{
		Scale: 0.33,
		x:     keys.stageW,
		Y:     0.66 * game.ScreenH,
	}
}

// Order of fields: logic -> drawing.
// Try to keep the order of initializing consistent with the order of fields.
type JudgmentComp struct {
	// Judgments [4]game.Judgment
	worst int // index of worst judgment
	anims [4]draws.Animation
	tween draws.Tween
}

// Passing args instead of opts is preferred
// to avoid using opts' fields directly.
func NewJudgmentComp(res JudgmentRes, opts JudgmentOpts) (comp JudgmentComp) {
	for i, frames := range res.framesList {
		a := draws.NewAnimation(frames, 40)
		a.MultiplyScale(opts.Scale)
		a.Locate(opts.x, opts.Y, draws.CenterMiddle)
		comp.anims[i] = a
	}

	tw := draws.Tween{}
	tw.Add(1.00, +0.15, 25, draws.EaseLinear)
	tw.Add(1.15, -0.15, 25, draws.EaseLinear)
	tw.Add(1.00, +0.0, 200, draws.EaseLinear)
	tw.Add(1.00, -0.25, 25, draws.EaseLinear)
	comp.tween = tw
	return
}

func (comp *JudgmentComp) Update(jis []int) {
	worst := blank // -1
	for _, ji := range jis {
		if worst < ji {
			worst = ji
		}
	}
	if worst >= kool {
		comp.worst = worst
		comp.anims[worst].Reset()
		comp.tween.Reset()
	}
}

// worstJudgment is guaranteed not to be blank,
// hence no panicked by index out of range.
func (comp JudgmentComp) Draw(dst draws.Image) {
	if comp.tween.IsFinished() {
		return
	}
	a := comp.anims[comp.worst]
	a.MultiplyScale(comp.tween.Current())
	a.Draw(dst)
}
