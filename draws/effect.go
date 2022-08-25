package draws

import "github.com/hajimehoshi/ebiten/v2"

// Todo: need a test for variadic parameters 'vs'
type Effect func(op *ebiten.DrawImageOptions, vs ...float64)
type Effecter struct {
	Countdown    int
	MaxCountdown int
	Value        float64
	Permanent    bool // Whether draws number endlessly.

	Colorer Effect
	Rotater Effect
	Scaler  Effect
	Translater
}

// AffineEffect indicates Translate effect.
// type AffineEffect func(float64, struct{ float64 float64 }) struct{ float64 float64 }

// Todo: need a test for variadic parameters 'vs'
// type Translater = func(age, w, h float64) (tx, ty float64)
type Translater func(op *ebiten.DrawImageOptions, vs ...float64)

// Any general-purpose value can be passed to Effect's Value.
// Todo: generalize input parameter of Op
func (e Effecter) Op(sprite Sprite) *ebiten.DrawImageOptions {
	age := e.Age()
	op := &ebiten.DrawImageOptions{}
	if e.Colorer != nil {
		v := age
		e.Colorer(op, v)
	}
	if e.Rotater != nil {
		v := age
		e.Rotater(op, v)
	}
	if e.Scaler != nil {
		v := age
		e.Scaler(op, v)
	}
	if e.Translater != nil {
		v := age
		w, h := sprite.Size()
		e.Translater(op, v, w, h)
	}
	return op
}

func (e Effecter) Age() float64 {
	return 1 - (float64(e.Countdown) / float64(e.MaxCountdown))
}

var Bower = func(op *ebiten.DrawImageOptions, age, w, h float64) {
	var tx, ty float64
	switch {
	case age < 0.1:
		ty = 0.85 * age * h
	case age >= 0.1 && age < 0.2:
		ty = 0.85 * (0.2 - age) * h
	}
	op.GeoM.Translate(tx, ty)
}
