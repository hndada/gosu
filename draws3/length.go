package draws

import "fmt"

type Unit int

const (
	Pixel Unit = iota
	Percent
	// RootPercent

	// Em // A unit of length in typography.
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

// Default unit is Pixel.
// func NewLength(value float64) Length { return Length{nil, value, Pixel} }
func NewLength(base *Length, value float64, unit Unit) Length {
	return Length{base, value, unit}
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

func NewLength2(base *Length2, x, y float64, unit Unit) Length2 {
	return Length2{
		X: NewLength(&base.X, x, unit),
		Y: NewLength(&base.Y, y, unit),
	}
}

func (l Length2) String() string {
	return fmt.Sprintf("{%s, %s}", l.X, l.Y)
}

func (l Length2) Pixel() XY {
	return XY{l.X.Pixel(), l.Y.Pixel()}
}

func (l *Length2) Add(px XY) {
	l.X.Add(px.X)
	l.Y.Add(px.Y)
}

func (l *Length2) Mul(scale XY) {
	l.X.Value *= scale.X
	l.Y.Value *= scale.Y
}
