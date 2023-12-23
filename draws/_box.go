package draws

// type Box interface {
// 	Size() Vector2
// 	At() Vector2
// }

// type Box interface {
// 	W() float64
// 	H() float64
// 	X() float64
// 	Y() float64
// }

type WHXY struct{ W, H, X, Y float64 }
type Box = WHXY
