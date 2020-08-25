package gosu

import (
	"github.com/faiface/beep"
	"github.com/hajimehoshi/ebiten"
)

const PlaySceneBufferTime float64 = 1500

type BasePlayScene struct {
	Streamer     beep.StreamSeekCloser
	StreamFormat beep.Format

	Tick  int64
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

// 이 방법을 하려면 tps가 게임 중에 변하지 않아야 함
// CurrentTPS가 약간 딱 떨어지지 않는 게 마음에 걸리지만, 곧 보충되어 결과적으로 일정히 유지 된다고 상정하겠음
func (s *BasePlayScene) Time() int64 {
	return s.Tick * Millisecond / int64(ebiten.MaxTPS())
}

type PlayScene interface {
}
