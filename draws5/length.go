package draws

import "fmt"

type Unit int

const (
	Pixel Unit = iota

	// Percent is a unit of relative length of the viewport.
	Percent

	// Em is a unit of length in typography.
	// Em
	// RootEm

	// Extra stands for additional size to parent's size.
	Extra
)

// Length.Pixel requires the kind of whxy to calculate pixel.
// This makes the function way more complicated.
type Length struct {
	Base  *Length
	Value float64
	Unit  Unit
}

func NewLength(v float64) Length {
	return Length{Value: v, Unit: Pixel}
}

func (l *Length) SetBase(base *Length, unit Unit) {
	l.Base = base
	l.Unit = unit
}

func (l Length) String() string {
	switch l.Unit {
	case Pixel:
		return fmt.Sprintf("%.2fpx", l.Value)
	case Percent:
		return fmt.Sprintf("%.2f%%", l.Value)
	case Extra:
		return fmt.Sprintf("+%.2fpx", l.Value)
	}
	return fmt.Sprintf("%.2fpx", l.Value)
}

func (l Length) Pixel() float64 {
	switch l.Unit {
	case Pixel:
		return l.Value
	case Percent:
		ratio := l.Value / 100.0
		return ratio * l.Base.Pixel()
	case Extra:
		return l.Value + l.Base.Pixel()
	}
	return l.Value
}

func (l *Length) Add(px float64) {
	switch l.Unit {
	case Pixel:
		l.Value += px
	case Percent:
		ratio := px / l.Base.Pixel()
		l.Value += ratio * 100
	case Extra:
		l.Value += px
	default:
		l.Value += px
	}
}

type Length2 struct{ X, Y Length }

func NewLength2(x, y float64) Length2 {
	return Length2{
		Length{Value: x},
		Length{Value: y},
	}
}

func (l2 *Length2) SetValues(x, y float64) {
	l2.X.Value = x
	l2.Y.Value = y
}

func (l2 *Length2) SetBase(base *Length2, unit Unit) {
	l2.X.SetBase(&base.X, unit)
	l2.Y.SetBase(&base.Y, unit)
}

func (l Length2) String() string {
	return fmt.Sprintf("{%s, %s}", l.X, l.Y)
}

func (l Length2) Pixels() XY {
	return XY{l.X.Pixel(), l.Y.Pixel()}
}

func (l *Length2) Add(vs XY) {
	l.X.Value += vs.X
	l.Y.Value += vs.Y
}

func (l *Length2) Mul(scales XY) {
	l.X.Value *= scales.X
	l.Y.Value *= scales.Y
}
