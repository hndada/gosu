package draws

var Screen = Box{
	Size: NewLength2(1280, 720),
}

// func ScreenSize() XY {
// 	return XY{
// 		Screen.Size.X.Value,
// 		Screen.Size.Y.Value,
// 	}
// }

func SetScreenSize(x, y float64) {
	Screen.Size.X.Value = x
	Screen.Size.Y.Value = y
}
