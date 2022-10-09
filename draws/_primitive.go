package draws

// type Position struct {
// 	// XY
// 	X, Y float64
// 	Origin
// }

// func (p Position) LeftTopX(w float64) float64 {
// 	switch p.Origin.PositionX() {
// 	case OriginLeft:
// 		return p.X
// 	case OriginCenter:
// 		return p.X - w/2
// 	case OriginRight:
// 		return p.X - w
// 	}
// 	panic("no reach")
// }
// func (p Position) LeftTopY(h float64) float64 {
// 	switch p.Origin.PositionY() {
// 	case OriginTop:
// 		return p.Y
// 	case OriginMiddle:
// 		return p.Y - h/2
// 	case OriginBottom:
// 		return p.Y - h
// 	}
// 	panic("no reach")
// }

// func (p Position) Move(tx, ty float64) {
// 	p.X += tx
// 	p.Y += ty
// }

// type Size struct {
// 	W, H float64
// 	Src  struct{ W, H float64 }
// 	// Scale struct{ W, H float64 }
// }

// func (s *Size) SetScale(scale float64) {
// 	s.SetScaleWH(scale, scale)
// }
// func (s *Size) SetScaleWH(scaleW, scaleH float64) {
// 	s.W *= scaleW
// 	s.H *= scaleH
// }

// type XY struct{ X, Y float64 }

// func (xy *XY) Move(xy2 XY) {
// 	xy.X += xy2.X
// 	xy.Y += xy2.Y
// }
// func (xy XY) Move2(xy2 XY) XY {
// 	return XY{
// 		X: xy.X + xy2.X,
// 		Y: xy.Y + xy2.Y,
// 	}
// }
