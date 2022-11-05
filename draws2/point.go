package draws

type Point struct{ X, Y float64 }

// Input integers.
func IntPt(x, y int) Point   { return Point{float64(x), float64(y)} }
func Pt(x, y float64) Point  { return Point{x, y} }
func Scalar(v float64) Point { return Point{v, v} }
func (p Point) Add(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}
func (p Point) Sub(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}
func (p Point) Mul(q Point) Point {
	return Point{p.X * q.X, p.Y * q.Y}
}
func (p Point) Div(q Point) Point {
	return Point{p.X / q.X, p.Y / q.Y}
}
func (p Point) XY() (float64, float64) { return p.X, p.Y }

// Output integers.
func (p Point) XYInt() (int, int) { return int(p.X), int(p.Y) }
