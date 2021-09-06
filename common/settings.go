package common

import "image"

var Settings struct {
	ScreenSize image.Point
	MaxTPS     int

	JudgmentMeterScale float64 // 1ms 당 2px
	ScoreHeight        float64 // todo: ScoreHeight -> ScoreImageScale
	ComboHeight        float64 // todo: mode마다 combo 위치가 다르므로 game package에 있으면 안됨
	ComboPosition      float64
	ComboGap           float64
	BackgroundDimness  float64

	MusicVolume  float64
	SFXVolume    float64
	MasterVolume float64
}

func init() {
	Settings.ScreenSize = image.Pt(1600, 900)
	Settings.JudgmentMeterScale = 2
	Settings.ScoreHeight = 7
	Settings.ComboHeight = 10
	Settings.ComboPosition = 40
	Settings.ComboGap = 0.8
	Settings.MaxTPS = 60
	Settings.BackgroundDimness = 0.3
	Settings.MusicVolume = 0.25
	Settings.SFXVolume = 0.25
	Settings.MasterVolume = 1
}

func DisplayScale() float64 {
	return float64(Settings.ScreenSize.Y) / 100
}

// Scale returns scaled value based on screen size
func Scale(v float64) int {
	scale := float64(Settings.ScreenSize.Y) / 100
	return int(v * scale)
}
