package gosu

import (
	"github.com/faiface/beep"
	"github.com/hajimehoshi/ebiten"
)

const PlaySceneBufferTime float64 = 1500

// todo: bufferTime 1500ms 정도로 넣기, 그런데 현재 fps 5 나옴
// todo: mp3, 버퍼에 올리자. 오래 걸려도 되니까
// todo: examples/camera 참고하여 플레이 프레임 안정화
// sync with mp3, position
type BasePlayScene struct {
	Streamer     beep.StreamSeekCloser
	StreamFormat beep.Format

	Tick  func() float64
	Time  float64
	Score int64
	HP    float64
	Combo int32
}

func (s *BasePlayScene) RenderScore() *ebiten.Image {
	// ScoreOverlap
}

func (s *BasePlayScene) RenderHP() *ebiten.Image {

}
func (s *BasePlayScene) RenderCombo() *ebiten.Image {
	// ComboOverlap
}

type PlayScene interface {
}