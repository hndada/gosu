package game

import "github.com/hajimehoshi/ebiten/audio"

type PlayScene struct {
	Tick        int64
	AudioPlayer *audio.Player
}

type TransSceneArgs struct {
	// 자기 자신은 .(type)으로 알 수 있음
	Next string // next scene name
	Args interface{}
}
