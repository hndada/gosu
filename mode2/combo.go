package mode

import (
	"github.com/hndada/gosu/draws"
)

type ComboConfig struct {
	// Position x and y may be derived from other configuration values.
	PositionX float64 // e.g. from FieldPositionX
	PositionY float64 // e.g. from ComboPositionY

	Scale    float64
	DigitGap float64
	Bounce   float64 // 0.85
	Persist  bool
}

func NewComboConfig(screen ScreenConfig) ComboConfig {
	return ComboConfig{
		PositionX: 0.50 * screen.Size.X,
		PositionY: 0.50 * screen.Size.Y,
		Scale:     0.75,
		DigitGap:  -1,
		Bounce:    0.85,
		Persist:   false,
	}
}

type ComboComponent struct {
	Combo     int
	lastCombo int // to reset tween
	sprites   [10]draws.Sprite
	w         float64 // Combo's width is fixed.
	tween     draws.Tween
}

// Let's make NewCombo everytime when Combo is changed.
func NewComboComponent(cfg ComboConfig, imgs [10]draws.Image) (cc ComboComponent) {
	for i := 0; i < 10; i++ {
		sprite := draws.NewSprite(imgs[i])
		sprite.MultiplyScale(cfg.Scale)
		sprite.Locate(cfg.PositionX, cfg.PositionY, draws.CenterMiddle)
		cc.sprites[i] = sprite
	}
	// Size of the whole image is 0.5w + (n-1)(w+gap) + 0.5w.
	// Since sprites are already at anchor, no need to care of two 0.5w.
	cc.w = cc.sprites[0].Width() + cfg.DigitGap

	tw := draws.Tween{}
	tw.AppendTween(0, cfg.Bounce, ToTick(200), draws.EaseLinear)
	tw.AppendTween(cfg.Bounce, -cfg.Bounce, ToTick(100), draws.EaseLinear)
	tw.AppendTween(0, 0, ToTick(1500), draws.EaseLinear)
	if !cfg.Persist {
		tw.SetLoop(1)
	}
	return
}

func (cc *ComboComponent) Update() {
	cc.tween.Tick()
	if cc.lastCombo != cc.Combo {
		cc.lastCombo = cc.Combo
		cc.tween.Reset()
	}
}

// Each number has different width. Number 0's width is used as standard.
// ComboDrawer's Draw draws each number at constant x regardless of their widths.
func (cc ComboComponent) Draw(dst draws.Image) {
	if cc.tween.IsFinished() {
		return
	}
	if cc.Combo == 0 {
		return
	}

	vs := make([]int, 0)
	for v := cc.Combo; v > 0; v /= 10 {
		vs = append(vs, v%10) // Little endian.
	}

	tx := float64(len(vs)-1) * cc.w / 2
	for _, v := range vs {
		sprite := cc.sprites[v]
		ty := cc.tween.Current() * sprite.Height()
		sprite.Move(tx, ty)
		sprite.Draw(dst, draws.Op{})
		tx -= cc.w
	}
}
