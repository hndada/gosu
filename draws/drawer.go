package draws

// Countdown is for drawing a sprite for a while.
// DrawImageOptions is not commutative. Do translate at final stage.
type BaseDrawer struct {
	Countdown    int
	MaxCountdown int // Draw permanently when value is zero.
	// Effecter     Effecter
	// Translater   Effecter
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

// // For drawing Combo and Score.
// // Suppose Each Sprite's X and Y indicate Origin's point.
// type NumberDrawer struct {
// 	BaseDrawer
// 	Sprites     [10]Sprite
// 	SignSprites [3]Sprite // Dot, Comma, Percent.
// 	Origin
// 	DigitWidth    float64 // Use number 0's width.
// 	DigitGap      float64
// 	FractionDigit int // Negative or zero value infers no drawing fraction part.
// 	ZeroFill      int

// 	integer  int
// 	fraction int
// }

// func (d *NumberDrawer) Update(i, f int) {
// 	var reloaded bool
// 	if d.integer != i || d.fraction != f {
// 		d.integer = i
// 		d.fraction = f
// 		reloaded = true
// 	}
// 	d.BaseDrawer.Update(reloaded)
// }

// // NumberDrawer's Draw draws each number at the center of constant-width bound.
// func (d NumberDrawer) Draw(screen *ebiten.Image) {
// 	if d.MaxCountdown != 0 && d.Countdown == 0 {
// 		return
// 	}
// 	vs := make([]int, 0)
// 	for v := d.integer; v > 0; v /= 10 {
// 		vs = append(vs, v%10) // Little endian.
// 	}
// 	for i := len(vs); i < d.ZeroFill; i++ {
// 		vs = append(vs, 0)
// 	}

// 	w := d.DigitWidth - d.DigitGap
// 	var tx float64
// 	switch d.Origin {
// 	case OriginRightTop:
// 		tx = 0
// 	case OriginCenter:
// 		// Size of the whole image is 0.5w + (n-1)(w-gap) + 0.5w.
// 		tx = (float64((len(vs)-1))*w + d.DigitWidth) / 2
// 	}
// 	for _, v := range vs {
// 		sprite := d.Sprites[v]
// 		sprite.Move(-tx, 0)
// 		sprite.Draw(screen, nil)
// 		tx += w
// 	}
// }
