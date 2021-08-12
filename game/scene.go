package game

import (
	"image"

	"github.com/hajimehoshi/ebiten/audio"
)

type Scene struct {
	Tick        int64
	ScreenSize  image.Point
	AudioPlayer *audio.Player
}

type TransSceneArgs struct {
	// 자기 자신은 .(type)으로 알 수 있음
	Next string // next scene name
	Args interface{}
}
