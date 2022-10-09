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

func (a Align) Min(size Point) Point {

}
