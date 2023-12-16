package mode

import "github.com/hndada/gosu/draws"

type Combo struct {
	Combo int
	// for drawing
	cfg       *ComboConfig
	sprites   [10]draws.Sprite
	timer     draws.Timer2
	lastCombo int
}

type ComboConfig struct {
	Position float64 // x
	Scale    float64
	DigitGap float64
	Bounce   float64 // 0.85
	Persist  bool
}

func NewCombo(sprites [10]draws.Sprite, cfg *ComboConfig) (c Combo) {
	const duration = 2000 // ms
	c.cfg = cfg
	c.sprites = sprites
	c.timer = draws.NewFiniteTimer(ToTick(duration))
	return
}

func (c *Combo) Update() {
	c.timer.Ticker()
	if c.lastCombo != c.Combo {
		c.lastCombo = c.Combo
		c.timer.Reset()
	}
}

// Each number has different width. Number 0's width is used as standard.
// ComboDrawer's Draw draws each number at constant x regardless of their widths.
func (c Combo) Draw(screen draws.Image) {
	if !c.cfg.Persist && c.timer.IsDone() {
		return
	}

	if c.Combo == 0 {
		return
	}

	vs := make([]int, 0)
	for v := c.Combo; v > 0; v /= 10 {
		vs = append(vs, v%10) // Little endian.
	}

	// Size of the whole image is 0.5w + (n-1)(w+gap) + 0.5w.
	// Since sprites are already at anchor, no need to care of two 0.5w.
	digitWidth := c.sprites[0].Width()
	w := digitWidth + c.cfg.DigitGap
	tx := float64(len(vs)-1) * w / 2
	const (
		boundary1 = 0.05
		boundary2 = 0.1
	)
	for _, v := range vs {
		sprite := c.sprites[v]
		sprite.Move(tx, 0)
		age := c.timer.Age()
		if age < boundary1 {
			scale := 0.1 * c.timer.Progress(0, boundary1)
			sprite.Move(0, c.cfg.Bounce*sprite.Height()*scale)
		}
		if age >= boundary1 && age < boundary2 {
			scale := 0.1 - 0.1*c.timer.Progress(boundary1, boundary2)
			sprite.Move(0, c.cfg.Bounce*sprite.Height()*scale)
		}
		sprite.Draw(screen, draws.Op{})
		tx -= w
	}
}
