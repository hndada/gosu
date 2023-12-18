package piano

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/draws"
	mode "github.com/hndada/gosu/mode2"
)

type JudgmentConfig struct {
	positionX float64 // from FieldPositionX
	PositionY float64
	Scale     float64
}

func NewJudgmentConfig(screen mode.ScreenConfig, stage StageConfig) JudgmentConfig {
	return JudgmentConfig{
		positionX: stage.FieldPositionX,
		PositionY: 0.66 * screen.Size.Y,
		Scale:     0.33,
	}
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

type JudgmentComponent struct {
	Judgments  [4]mode.Judgment
	Counts     [4]int
	worst      mode.Judgment // worst judgment
	animations [4]draws.Animation
	tween      draws.Tween
}

func NewJudgmentComponent(cfg JudgmentConfig, framesList [4]draws.Frames) (jc JudgmentComponent) {
	jc.Judgments = [4]mode.Judgment(DefaultJudgments())
	for i, frames := range framesList {
		a := draws.NewAnimation(frames, mode.ToTick(40))
		a.MultiplyScale(cfg.Scale)
		a.Locate(cfg.positionX, cfg.PositionY, draws.CenterMiddle)
		jc.animations[i] = a
	}

	jc.tween = draws.Tween{}
	jc.tween.AppendTween(1, 0.15, mode.ToTick(25), draws.EaseLinear)
	jc.tween.AppendTween(1.15, -0.15, mode.ToTick(25), draws.EaseLinear)
	jc.tween.AppendTween(1, 0, mode.ToTick(200), draws.EaseLinear)
	jc.tween.AppendTween(1, -0.25, mode.ToTick(25), draws.EaseLinear)
	return
}

func LoadJudgmentImages(fsys fs.FS) (framesList [4]draws.Frames) {
	for i, name := range []string{"kool", "cool", "good", "miss"} {
		fname := fmt.Sprintf("piano/judgment/%s.png", name)
		framesList[i] = draws.NewFramesFromFilename(fsys, fname)
	}
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
func (jc JudgmentComponent) Draw(dst draws.Image) {
	if jc.tween.IsFinished() {
		return
	}
	a := jc.animations[jc.index()]
	a.MultiplyScale(jc.tween.Current())
	a.Draw(dst, draws.Op{})
}
