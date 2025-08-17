package piano

import (
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
)

type HitLights struct {
	keysAnim []draws.Animation
}

func NewHitLights(res *Resources, opts *Options, keyCount int) HitLights {
	hitLights := HitLights{}
	hitLights.keysAnim = make([]draws.Animation, keyCount)
	xs := opts.keyPositionXsMap[keyCount]
	for k := range hitLights.keysAnim {
		a := draws.NewAnimation(res.HitLightsFrames, 150)
		a.Scale(opts.HitLightImageScale)
		a.Locate(xs[k], opts.KeyPositionY-opts.HintHeight/2, draws.CenterMiddle)
		a.ColorScale.Scale(1, 1, 1, opts.HitLightOpacity)
		a.MaxLoop = 1
		hitLights.keysAnim[k] = a
	}
	return hitLights
}

// Tail also makes hit lighting on.
func (hitLights *HitLights) Update(kjk []game.JudgmentKind) {
	for k, jk := range kjk {
		if jk <= good {
			hitLights.keysAnim[k].Reset()
		}
	}
}

// HitLights.Draw draws hit lights when Normal is Hit or Tail is Released.
func (hitLights HitLights) Draw(dst draws.Image) {
	for _, a := range hitLights.keysAnim {
		if a.IsFinished() {
			continue
		}
		a.Draw(dst)
	}
}
