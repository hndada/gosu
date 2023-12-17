package piano

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/mode"
)

const (
	Kool = iota
	Cool
	Good
	Miss
)

type JudgmentComponent struct {
	Judgments  [4]mode.Judgment
	Counts     [4]int
	worst      mode.Judgment // worst judgment
	animations [4]draws.Animation
	tween      draws.Tween
}

type JudgmentConfig struct {
	FieldPosition *float64
	Position      float64
	Scale         float64
}

func DefaultJudgments() []mode.Judgment {
	return []mode.Judgment{
		{Window: 20, Weight: 1},
		{Window: 40, Weight: 1},
		{Window: 80, Weight: 0.5},
		{Window: 120, Weight: 0},
	}
}

func LoadJudgmentImages(fsys fs.FS) (framesList [4]draws.Frames) {
	for i, name := range []string{"kool", "cool", "good", "miss"} {
		fname := fmt.Sprintf("piano/judgment/%s.png", name)
		framesList[i] = draws.NewFramesFromFilename(fsys, fname)
	}
	return
}

func NewJudgmentComponent(framesList [4]draws.Frames, cfg JudgmentConfig) (jc JudgmentComponent) {
	jc.Judgments = [4]mode.Judgment(DefaultJudgments())
	for i, frames := range framesList {
		a := draws.NewAnimation(frames, mode.ToTick(40))
		a.MultiplyScale(cfg.Scale)
		a.Locate(*cfg.FieldPosition, cfg.Position, draws.CenterMiddle)
		jc.animations[i] = a
	}

	jc.tween = draws.Tween{}
	jc.tween.AppendTween(1, 0.15, mode.ToTick(25), draws.EaseLinear)
	jc.tween.AppendTween(1.15, -0.15, mode.ToTick(25), draws.EaseLinear)
	jc.tween.AppendTween(1, 0, mode.ToTick(200), draws.EaseLinear)
	jc.tween.AppendTween(1, -0.25, mode.ToTick(25), draws.EaseLinear)
	return
}

func (jc *JudgmentComponent) Update(js []mode.Judgment) {
	jc.worst = mode.Judgment{}
	for _, j := range js {
		if jc.worst.Window < j.Window { // j is worse
			jc.worst = j
		}
	}
	if !jc.worst.IsBlank() {
		jc.animations[jc.index()].Reset()
		jc.tween.Reset()
	}
}

func (jc JudgmentComponent) index() int {
	for i, j := range jc.Judgments {
		if j.Is(jc.worst) {
			return i
		}
	}
	return len(jc.Judgments) // blank judgment
}

func (jc *JudgmentComponent) Increment(j mode.Judgment) {
	for i, j2 := range jc.Judgments {
		if j.Is(j2) {
			jc.Counts[i]++
			break
		}
	}
}

func (jc JudgmentComponent) kool() mode.Judgment { return jc.Judgments[Kool] }
func (jc JudgmentComponent) cool() mode.Judgment { return jc.Judgments[Cool] }
func (jc JudgmentComponent) good() mode.Judgment { return jc.Judgments[Good] }
func (jc JudgmentComponent) miss() mode.Judgment { return jc.Judgments[Miss] }

// worstJudgment is guaranteed not to be blank,
// hence no panicked by index out of range.
func (jc JudgmentComponent) Draw(screen draws.Image) {
	if jc.tween.IsFinished() {
		return
	}
	a := jc.animations[jc.index()]
	a.MultiplyScale(jc.tween.Current())
	a.Draw(screen, draws.Op{})
}
