package mode

import (
	"github.com/hndada/gosu/draws"
)

type Combo struct {
	Combo     int
	lastCombo int // to reset tween
	sprites   [10]draws.Sprite
	w         float64 // Combo's width is fixed.
	tween     draws.Tween
}

type ComboConfig struct {
	ScreenSize    *draws.Vector2
	FieldPosition *float64

	Position float64 // x
	Scale    float64
	DigitGap float64
	Bounce   float64 // 0.85
	Persist  bool
}

// Let's make NewCombo everytime when Combo is changed.
func NewCombo(imgs [10]draws.Image, cfg ComboConfig) (c Combo) {
	x := cfg.ScreenSize.X * *cfg.FieldPosition
	y := cfg.ScreenSize.Y * cfg.Position
	for i := 0; i < 10; i++ {
		sprite := draws.NewSprite(imgs[i])
		sprite.MultiplyScale(cfg.Scale)
		sprite.Locate(x, y, draws.CenterMiddle)
		c.sprites[i] = sprite
	}
	// Size of the whole image is 0.5w + (n-1)(w+gap) + 0.5w.
	// Since sprites are already at anchor, no need to care of two 0.5w.
	c.w = c.sprites[0].Width() + cfg.DigitGap

	tw := draws.Tween{}
	tw.AppendTween(0, cfg.Bounce, ToTick(200), draws.EaseLinear)
	tw.AppendTween(cfg.Bounce, 0, ToTick(100), draws.EaseLinear)
	tw.AppendTween(0, 0, ToTick(1500), draws.EaseLinear)
	if !cfg.Persist {
		tw.SetLoop(1)
	}
	return
}

func (c *Combo) Update() {
	c.tween.Tick()
	if c.lastCombo != c.Combo {
		c.lastCombo = c.Combo
		c.tween.Reset()
	}
}

// Each number has different width. Number 0's width is used as standard.
// ComboDrawer's Draw draws each number at constant x regardless of their widths.
func (c Combo) Draw(screen draws.Image) {
	if c.tween.IsFinished() {
		return
	}
	if c.Combo == 0 {
		return
	}

	vs := make([]int, 0)
	for v := c.Combo; v > 0; v /= 10 {
		vs = append(vs, v%10) // Little endian.
	}

	tx := float64(len(vs)-1) * c.w / 2
	for _, v := range vs {
		sprite := c.sprites[v]
		ty := c.tween.Current() * sprite.Height()
		sprite.Move(tx, ty)
		sprite.Draw(screen, draws.Op{})
		tx -= c.w
	}
}
