package draws

// Screen is the root Box for the screen.
var Screen = Box{
	Size: NewLength2(640, 480),
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
