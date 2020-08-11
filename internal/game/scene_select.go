package game

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hndada/gosu/mode/mania"
)

// 오디오 플레이어?
type SceneSelect struct {
	// 차트 리스트
	// 커서
	// 그룹 (디렉토리 트리)
	// 현재 정렬 기준
}

// 모든 box 생성?
// 현재 선택된 차트 focus (커서) 위치 고정

// 위쪽/왼쪽: 커서 -1
// 아래쪽/오른쪽: 커서 +1
// +시프트: 그룹 이동
func (s *SceneSelect) Update(g *Game) error {
	if ebiten.IsKeyPressed(ebiten.Key1) {
		c := mania.NewChart(`C:\Users\hndada\Documents\GitHub\hndada\gosu\mode\mania\test\test_ln.osu`)
		g.NextScene = NewSceneMania(g.Options, c)
		g.TransCountdown = 99
	}
	// 키 입력 받으면 play scene으로 이동
	return nil
}

func (s *SceneSelect) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "SceneSelect: Press Key 1")
}
