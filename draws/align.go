package draws

type Align struct{ X, Y int }

const (
	AlignLeft = iota
	AlignCenter
	AlignRight
)
const (
	AlignTop = iota
	AlignMiddle
	AlignBottom
)

func (a Align) Min(min, size Point) Point {
	switch a.X {
	case AlignLeft:

	case AlignCenter:
	case AlignBottom:
	}
}
