package common

import "image"

var Settings struct {
	ScreenSizeX int
	ScreenSizeY int
	MaxTPS      int

	ScoreHeight       float64 // TODO: ScoreHeight -> ScoreImageScale
	ComboHeight       float64 // TODO: It should be in each game package: a game has different combo position
	ComboPosition     float64
	ComboGap          float64
	BackgroundDimness float64

	MusicVolume  float64
	SFXVolume    float64
	MasterVolume float64

	ScoreMode       int
	IsAuto          bool
	AutoUnstability float64 // 0~100; 0 makes the play 'Perfect'
}

func init() {
	Settings.ScreenSizeX = 800
	Settings.ScreenSizeY = 600
	Settings.MaxTPS = 60

	Settings.ScoreHeight = 7
	Settings.ComboHeight = 10
	Settings.ComboPosition = 40
	Settings.ComboGap = 0.8
	Settings.BackgroundDimness = 0.3

	// Settings.MusicVolume = 0.25
	// Settings.SFXVolume = 0.25
	Settings.MasterVolume = 0.5

	Settings.IsAuto = true
	Settings.AutoUnstability = 0 // 0~100; 0 makes the play 'Perfect'
}

func DisplayScale() float64 {
	return float64(Settings.ScreenSizeY) / 100
}

// Scale returns scaled value based on screen size
func Scale(v float64) int {
	scale := float64(Settings.ScreenSizeY) / 100
	return int(v * scale)
}

func ScreenSize() image.Point {
	return image.Pt(Settings.ScreenSizeX, Settings.ScreenSizeX)
}
