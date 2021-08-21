package game

import "image"

var Settings struct {
	ScreenSize         image.Point
	JudgmentMeterScale float64
}

func Scale() float64 {
	return float64(Settings.ScreenSize.Y) / 100
}
