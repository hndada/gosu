package piano

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/draws"
	mode "github.com/hndada/gosu/mode2"
)

type JudgmentRes struct {
	framesList [4]draws.Frames
}

func (res *JudgmentRes) Load(fsys fs.FS) {
	for i, name := range []string{"kool", "cool", "good", "miss"} {
		fname := fmt.Sprintf("piano/judgment/%s.png", name)
		res.framesList[i] = draws.NewFramesFromFilename(fsys, fname)
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
		Y:     0.66 * mode.ScreenH,
	}
}

// Order of fields: logic -> drawing.
// Try to keep the order of initializing consistent with the order of fields.
type JudgmentComp struct {
	Judgments [4]mode.Judgment
	Counts    [4]int
	worst     mode.Judgment // worst judgment
	anims     [4]draws.Animation
	tween     draws.Tween
}

// Passing args instead of opts is preferred
// to avoid using opts' fields directly.
func NewJudgmentComp(res JudgmentRes, opts JudgmentOpts) (comp JudgmentComp) {
	js := DefaultJudgments()
	comp.Judgments = [4]mode.Judgment(js)

	for i, frames := range res.framesList {
		anim := draws.NewAnimation(frames, mode.ToTick(40))
		anim.MultiplyScale(opts.Scale)
		anim.Locate(opts.x, opts.Y, draws.CenterMiddle)
		comp.anims[i] = anim
	}

	tween := draws.Tween{}
	tween.AppendTween(1, 0.15, mode.ToTick(25), draws.EaseLinear)
	tween.AppendTween(1.15, -0.15, mode.ToTick(25), draws.EaseLinear)
	tween.AppendTween(1, 0, mode.ToTick(200), draws.EaseLinear)
	tween.AppendTween(1, -0.25, mode.ToTick(25), draws.EaseLinear)
	comp.tween = tween
	return
}

const (
	Kool = iota
	Cool
	Good
	Miss
)

func DefaultJudgments() []mode.Judgment {
	return []mode.Judgment{
		{Window: 20, Weight: 1},
		{Window: 40, Weight: 1},
		{Window: 80, Weight: 0.5},
		{Window: 120, Weight: 0},
	}
}

func (comp JudgmentComp) Kool() mode.Judgment { return comp.Judgments[Kool] }
func (comp JudgmentComp) Cool() mode.Judgment { return comp.Judgments[Cool] }
func (comp JudgmentComp) Good() mode.Judgment { return comp.Judgments[Good] }
func (comp JudgmentComp) Miss() mode.Judgment { return comp.Judgments[Miss] }

func (comp *JudgmentComp) Update(js []mode.Judgment) {
	comp.worst = mode.Judgment{}
	for _, j := range js {
		if comp.worst.Window < j.Window { // j is worse
			comp.worst = j
		}
	}
	if !comp.worst.IsBlank() {
		comp.anims[comp.index()].Reset()
		comp.tween.Reset()
	}
}

func (comp JudgmentComp) index() int {
	for i, j := range comp.Judgments {
		if j.Is(comp.worst) {
			return i
		}
	}
	return len(comp.Judgments) // blank judgment
}

func (comp *JudgmentComp) Increment(j mode.Judgment) {
	for i, j2 := range comp.Judgments {
		if j.Is(j2) {
			comp.Counts[i]++
			break
		}
	}
}

// worstJudgment is guaranteed not to be blank,
// hence no panicked by index out of range.
func (comp JudgmentComp) Draw(dst draws.Image) {
	if comp.tween.IsFinished() {
		return
	}
	anim := comp.anims[comp.index()]
	anim.MultiplyScale(comp.tween.Current())
	anim.Draw(dst, draws.Op{})
}
