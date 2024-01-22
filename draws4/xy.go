package draws

type XY struct{ X, Y float64 }

// func Vec2(x, y float64) XY          { return XY{x, y} }
func NewXYFromInts(x, y int) XY         { return XY{float64(x), float64(y)} }
func NewXYFromScalar(v float64) XY      { return XY{v, v} }
func (a XY) Add(b XY) XY                { return XY{a.X + b.X, a.Y + b.Y} }
func (a XY) Sub(b XY) XY                { return XY{a.X - b.X, a.Y - b.Y} }
func (a XY) Mul(b XY) XY                { return XY{a.X * b.X, a.Y * b.Y} }
func (a XY) Div(b XY) XY                { return XY{a.X / b.X, a.Y / b.Y} }
func (a XY) Scale(scale float64) XY     { return XY{a.X * scale, a.Y * scale} }
func (a XY) Values() (float64, float64) { return a.X, a.Y }
func (a XY) IntValues() (int, int)      { return int(a.X), int(a.Y) }
