package draws

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// DrawImageOptions is not commutative.
// Do Translate at final stage: Do Rotate or Scale first.
// Effecter should belong to Drawer, not to Sprite. There might be an animation.
type BaseDrawer struct {
	// Sprites []Sprite
	// Index   int

	// Effecter
	Countdown    int
	MaxCountdown int
	Permanent    bool // Whether draws number endlessly.
	// Value        float64
	Effecters []Effecter
	Translater
}

func (d *BaseDrawer) Update(i int) {
	if d.Countdown > 0 {
		d.Countdown--
	}
	if d.Index != i {
		d.Countdown = d.MaxCountdown
		d.Index = i
	}
}
func (d BaseDrawer) Draw(screen *ebiten.Image) {
	s := d.Sprites[d.Index]
	s.Draw(screen, d.Effecter.Op(s))
}
func (d BaseDrawer) Age() float64 {
	return 1 - (float64(d.Countdown) / float64(d.MaxCountdown))
}

// For drawing Combo and Score.
// Suppose Each Sprite's X and Y indicate Origin's point.
type NumberDrawer struct {
	Sprites       [10]Sprite
	SignSprites   [3]Sprite // Dot, Comma, Percent.
	DigitWidth    float64   // Use number 0's width.
	DigitGap      float64
	Integer       int
	Fraction      int
	FractionDigit int // Negative or zero value infers no drawing fraction part.
	Effecter
}

func (d *NumberDrawer) Update(i, f int) {
	if d.Countdown > 0 {
		d.Countdown--
	}
	if d.Integer != i || d.Fraction != f {
		d.Countdown = d.MaxCountdown
		d.Integer = i
		d.Fraction = f
	}
}
func (d NumberDrawer) IsZero() bool { return d.Integer == 0 && d.Fraction == 0 }

// NumberDrawer's Draw draws each number at the center of constant-width bound.
func (d NumberDrawer) Draw(screen *ebiten.Image) {
	if d.IsZero() || (!d.Permanent && d.Countdown == 0) {
		return
	}
	vs := make([]int, 0)
	for v := d.Integer; v > 0; v /= 10 {
		vs = append(vs, v%10) // Little endian.
	}
	// Size of the whole image is 0.5w + (n-1)(w-gap) + 0.5w.
	w := d.DigitWidth - d.DigitGap
	tx := float64((len(vs)-1))*w + d.DigitWidth/2
	for _, v := range vs {
		s := d.Sprites[v]
		op := d.Op(s)
		op.GeoM.Translate(tx, 0)
		s.Draw(screen, op)
		tx -= w
	}
}

const (
	SignDot = iota
	SignComma
	SignPercent
)
