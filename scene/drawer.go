package scene

import "github.com/hndada/gosu/draws"

// Order of fields of drawer: updating fields, others fields, sprites
type BackgroundDrawer struct {
	Brightness *float64
	Sprite     draws.Sprite
}

func (d BackgroundDrawer) Draw(dst draws.Image) {
	op := draws.Op{}
	op.ColorM.ChangeHSV(0, 1, *d.Brightness)
	d.Sprite.Draw(dst, op)
}

type IntDrawer struct {
	draws.Timer
	DigitWidth float64
	DigitGap   float64
	Value      int
	Bounce     float64
	Numbers    [10]draws.Sprite
	// Signs      [3]draws.Sprite
}

// Each number has different width. Number 0's width is used as standard.
func (d *IntDrawer) Update(v int) {
	d.Ticker()
	if d.Value != v {
		d.Value = v
		d.Timer.Reset()
	}
}

// IntDrawer's Draw draws each number at constant x regardless of their widths.
func (d IntDrawer) Draw(dst draws.Image) {
	if d.Done() {
		return
	}
	if d.Value == 0 {
		return
	}
	vs := make([]int, 0)
	for v := d.Value; v > 0; v /= 10 {
		vs = append(vs, v%10) // Little endian.
	}

	// Size of the whole image is 0.5w + (n-1)(w+gap) + 0.5w.
	// Since sprites are already at origin, no need to care of two 0.5w.
	w := d.DigitWidth + d.DigitGap
	tx := float64(len(vs)-1) * w / 2
	const (
		bound0 = 0.05
		bound1 = 0.1
	)
	for _, v := range vs {
		sprite := d.Numbers[v]
		sprite.Move(tx, 0)
		age := d.Age()
		if age < bound0 {
			scale := 0.1 * d.Progress(0, bound0)
			sprite.Move(0, d.Bounce*sprite.H()*scale)
		}
		if age >= bound0 && age < bound1 {
			scale := 0.1 - 0.1*d.Progress(bound0, bound1)
			sprite.Move(0, d.Bounce*sprite.H()*scale)
		}
		sprite.Draw(dst, draws.Op{})
		tx -= w
	}
}
