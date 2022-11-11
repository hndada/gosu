package draws

type Vector2 struct{ X, Y float64 }

// Input integers.
func IntVec2(x, y int) Vector2  { return Vector2{float64(x), float64(y)} }
func Vec2(x, y float64) Vector2 { return Vector2{x, y} }
func Scalar(v float64) Vector2  { return Vector2{v, v} }
func (v Vector2) Add(q Vector2) Vector2 {
	return Vector2{v.X + q.X, v.Y + q.Y}
}
func (v Vector2) Sub(q Vector2) Vector2 {
	return Vector2{v.X - q.X, v.Y - q.Y}
}
func (v Vector2) Mul(q Vector2) Vector2 {
	return Vector2{v.X * q.X, v.Y * q.Y}
}
func (v Vector2) Div(q Vector2) Vector2 {
	return Vector2{v.X / q.X, v.Y / q.Y}
}
func (v Vector2) XY() (float64, float64) { return v.X, v.Y }

// Output integers.
func (v Vector2) XYInt() (int, int) { return int(v.X), int(v.Y) }
