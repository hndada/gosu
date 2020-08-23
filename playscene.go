package gosu

import (
	"github.com/faiface/beep"
	"github.com/hajimehoshi/ebiten"
)

const PlaySceneBufferTime float64 = 1500

// todo: bufferTime 1500ms 정도로 넣기, 그런데 현재 fps 5 나옴
// todo: mp3, 버퍼에 올리자. 오래 걸려도 되니까
// todo: examples/camera 참고하여 플레이 프레임 안정화
// todo: float64, fixed로 고치기 생각
// sync with mp3, position


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
