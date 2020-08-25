package gosu

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/ebitenui"
	"github.com/hndada/gosu/mode/mania"
)

// todo: 차트 패널
// todo: Songs 폴더 읽는 로직 만들기 - rule 포함
// 모든 box 생성?

// todo: 로딩일 때 기다리는 로직
// Loading 이라는 별도의 Lock을 둔 이상, 특별히 채널은 필요없는거 아닌가?
// 비트맵 로딩 15초 후 timeout

type SceneSelect struct {
	g *Game
	Buttons     []ebitenui.Button
	ChartPanels []ChartPanel
	cursor      int
	// 그룹 (디렉토리 트리)
	// 현재 정렬 기준
}

func NewSceneSelect() *SceneSelect {
	ebiten.SetWindowTitle("gosu")
	return nil
}

// 위쪽/왼쪽: 커서 -1
// 아래쪽/오른쪽: 커서 +1
// +시프트: 그룹 이동
func (s *SceneSelect) Update() error {
	if ebiten.IsKeyPressed(ebiten.Key1) {
		c := &mania.Chart{}
		s.g.NextScene = NewSceneMania(g, c) // todo: go func()?
		s.g.TransCountdown = s.g.MaxTransCountDown()
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		s.cursor++
		// if s.cursor <= len() {
		// 	s.cursor = 0
		// }
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		s.cursor--
		// if s.cursor < 0 {
		// 	s.cursor = len() - 1
		// }
	}
	for _, p := range s.ChartPanels {
		p.Update()
	}
	return nil
}

// 현재 선택된 차트 focus 틀 위치 고정
func (s *SceneSelect) Draw(screen *ebiten.Image) {
	for _, p := range s.ChartPanels {
		p.Draw(screen)
	}
	for _, b := range s.Buttons {
		b.Draw(screen)
	}
}
