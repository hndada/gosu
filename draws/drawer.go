package draws

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

// DrawImageOptions is not commutative.
// Rotate or Scale first, do translate at final stage.
// Effecter should belong to Drawer, not to Sprite. There might be an animation.
type BaseDrawer struct {
	// Sprites []Sprite
	// Index        int
	// Value        float64
	Countdown    int
	MaxCountdown int // Draw permanently when value is zero.
	// Permanent    bool // Whether draws endlessly.
	Effecter
	Translater
}

func (d *BaseDrawer) Update(reloaded bool) {
	if d.Countdown > 0 {
		d.Countdown--
	}
	if reloaded {
		d.Countdown = d.MaxCountdown
	}
}
func (d BaseDrawer) Age() float64 {
	return 1 - (float64(d.Countdown) / float64(d.MaxCountdown))
}

// func (d Drawer) Draw(screen *ebiten.Image) {
// 	s := d.Sprites[d.Index]
// 	s.Draw(screen, d.Effecter.Op(s))
// }

// For drawing Combo and Score.
// Suppose Each Sprite's X and Y indicate Origin's point.
type NumberDrawer struct {
	BaseDrawer
	Sprites       [10]Sprite
	SignSprites   [3]Sprite // Dot, Comma, Percent.
	DigitWidth    float64   // Use number 0's width.
	DigitGap      float64
	Integer       int
	Fraction      int
	FractionDigit int // Negative or zero value infers no drawing fraction part.
	ZeroFill      int
	Origin
	// DrawZero      bool
	// Effecter
}

func (d *NumberDrawer) Update(i, f int) {
	// if d.Countdown > 0 {
	// 	d.Countdown--
	// }
	var reload bool
	if d.Integer != i || d.Fraction != f {
		d.Integer = i
		d.Fraction = f
		reload = true
		// d.Countdown = d.MaxCountdown
	}
	d.BaseDrawer.Update(reload)
}

// func (d NumberDrawer) IsZero() bool { return d.Integer == 0 && d.Fraction == 0 }

// NumberDrawer's Draw draws each number at the center of constant-width bound.
func (d NumberDrawer) Draw(screen *ebiten.Image) {
	if d.MaxCountdown != 0 && d.Countdown == 0 {
		return
	}
	vs := make([]int, 0)
	for v := d.Integer; v > 0; v /= 10 {
		vs = append(vs, v%10) // Little endian.
	}
	for i := len(vs); i < d.ZeroFill; i++ {
		vs = append(vs, 0)
	}
	// Size of the whole image is 0.5w + (n-1)(w-gap) + 0.5w.
	w := d.DigitWidth - d.DigitGap
	var tx float64
	fmt.Println(d.Origin)
	switch d.Origin {
	case OriginRightTop:

		tx = 0
	case OriginCenter:
		fmt.Println("uh")
		tx = (float64((len(vs)-1))*w + d.DigitWidth) / 2
	}
	for _, v := range vs {
		sprite := d.Sprites[v]
		sprite.Move(-tx, 0)
		sprite.Draw(screen, nil)
		// op := d.Op(s)
		// op.GeoM.Translate(tx, 0)
		// s.Draw(screen, op)
		tx += w
	}
}

const (
	SignDot = iota
	SignComma
	SignPercent
)
