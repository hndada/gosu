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
	index      int
	lastIndex  int
	animations [4]draws.Animation
	tween      draws.Tween
}

type JudgmentConfig struct {
	FieldPosition *float64
	Position      float64
	Scale         float64
}

func LoadJudgmentImages(fsys fs.FS) (framesList [4]draws.Frames) {
	for i, name := range []string{"kool", "cool", "good", "miss"} {
		fname := fmt.Sprintf("piano/judgment/%s.png", name)
		framesList[i] = draws.NewFramesFromFilename(fsys, fname)
	}
	return
}

func NewJudgmentComponent(framesList [4]draws.Frames, cfg JudgmentConfig) (jc JudgmentComponent) {
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

// worstJudgment is guaranteed not to be blank,
// hence no panicked by index out of range.
func (jc JudgmentComponent) Draw(screen draws.Image) {
	if jc.tween.IsFinished() {
		return
	}
	a := jc.animations[jc.index]
	a.MultiplyScale(jc.tween.Current())
	a.Draw(screen, draws.Op{})
}
