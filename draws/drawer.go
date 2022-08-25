package draws

import "github.com/hajimehoshi/ebiten/v2"

type Effect = func(float64, *ebiten.DrawImageOptions)

// DrawImageOptions is not commutative.
// Rotate -> Scale -> Translate.
// Combo and Score will have separate drawer called NumberDrawer.
type Drawer struct {
	Sprites      []Sprite
	Index        int
	Countdown    int
	MaxCountdown int

	Rotater    Effect
	Scaler     Effect
	Translater Effect
	Colorer    Effect
}

func (d *Drawer) Update(v any) {}

// First parameter of Effection function is age.
