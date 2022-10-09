package draws

type Point struct{ X, Y float64 }

func Pt(x, y int) Point { return Point{float64(x), float64(y)} }
func (p Point) Add(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}
func (p Point) Sub(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}
func (p Point) Mul(q Point) Point {
	return Point{p.X * q.X, p.Y * q.Y}
}
func (p Point) XY() (float64, float64) { return p.X, p.Y }
func (p Point) XYInt() (int, int)      { return int(p.X), int(p.Y) }
