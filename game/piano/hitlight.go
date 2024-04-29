package piano

import (
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
)

type HitLightsComponent struct {
	keysAnim []draws.Animation
}

func NewHitLightsComponent(res *Resources, opts *Options, keyCount int) (cmp HitLightsComponent) {
	cmp.keysAnim = make([]draws.Animation, keyCount)
	xs := opts.keyPositionXsMap[keyCount]
	for k := range cmp.keysAnim {
		a := draws.NewAnimation(res.HitLightsFrames, 150)
		a.Scale(opts.HitLightImageScale)
		a.Locate(xs[k], opts.KeyPositionY, draws.CenterBottom)
		a.ColorScale.Scale(1, 1, 1, opts.HitLightOpacity)
		a.SetMaxLoop(1)
		cmp.keysAnim[k] = a
	}
	return
}

// Tail also makes hit lighting on.
func (cmp *HitLightsComponent) Update(kjk []game.JudgmentKind) {
	for k, jk := range kjk {
		if jk <= good {
			cmp.keysAnim[k].Reset()
		}
	}
}

// HitLightsComponent.Draw draws hit lights when Normal is Hit or Tail is Released.
func (cmp HitLightsComponent) Draw(dst draws.Image) {
	for _, a := range cmp.keysAnim {
		if a.IsFinished() {
			continue
		}
		a.Draw(dst)
	}
}
