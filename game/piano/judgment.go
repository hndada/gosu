package piano

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
)

type JudgmentResources struct {
	framesList [4]draws.Frames
}

func (res *JudgmentResources) Load(fsys fs.FS) {
	for i, name := range []string{"kool", "cool", "good", "miss"} {
		fname := fmt.Sprintf("piano/judgment/%s.png", name)
		res.framesList[i] = draws.NewFramesFromFile(fsys, fname)
	}
	return
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
	worst int // index of worst judgment
	anims [4]draws.Animation
	tween draws.Tween
}

func NewJudgmentComponent(res JudgmentResources, opts JudgmentOptions) (cmp JudgmentComponent) {
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

func (cmp *JudgmentComponent) Update(jis []int) {
	worst := blank // -1
	for _, ji := range jis {
		if worst < ji {
			worst = ji
		}
	}
	if worst >= kool { // 0
		cmp.worst = worst
		cmp.anims[worst].Reset()
		cmp.tween.Reset()
	}
}

func (cmp JudgmentComponent) Draw(dst draws.Image) {
	if cmp.tween.IsFinished() {
		return
	}
	// worstJudgment is guaranteed not to be blank,
	// hence no panicked by index out of range.
	a := cmp.anims[cmp.worst]
	a.MultiplyScale(cmp.tween.Current())
	a.Draw(dst)
}
