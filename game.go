package gosu

import (
	"github.com/hndada/gosu/graphics"
	"github.com/hndada/gosu/settings"
	"path/filepath"

	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/mode/mania"
	"github.com/hndada/rg-parser/osugame/osu"
)

// todo: 훈련소 및 군대 가서 할 것들 정리 필요
// todo: reference book: 참고용 ui; 곡선택창 같은거
// 다른 곳으로 날아가도, 외출 때 코딩을 하겠음

// keyboard - 게임
// gob, toml - 설정 저장
// beep (sound) - 음원 및 효과음 재생
// font - 패널 그리기
// todo: float64, fixed로 고치기 생각

// 체크박스 같은거 다시 그리기 -> 우선 메모장으로 직접 설정하게
// PlayIntro, PlayExit
const Millisecond = 1000

// reset    save cancel
// 설정 켜면 임시 세팅이 생성, 임시 세팅으로 실시간 보여주기
// save 누르면 실제 세팅으로 값복사
// game에서 세팅 바꾸면 Sprite 자동 갱신
type Game struct {
	settings.Settings
	graphics.GameSprites
	Scene        Scene
	SceneChanger *SceneChanger
	// Input        input.Input
}

// todo: 소리 재생
// Scene이 Game을 control하는 주체
type Scene interface {
	Update() error
	Draw(screen *ebiten.Image) // Draws scene to screen
}

func NewGame() *Game {
	g := &Game{}
	g.Settings.Load()
	g.Scene = g.NewSceneTitle()
	g.SceneChanger = NewSceneChanger()
	return g
}
func (g *Game) Update(screen *ebiten.Image) error {
	if !g.SceneChanger.done() {
		return g.SceneChanger.Update(g)
	}
	return g.Scene.Update()
}

// 이미지의 method Draw는 input으로 들어온 screen을 그리는 함수
func (g *Game) Draw(screen *ebiten.Image) {
	if !g.SceneChanger.done() {
		g.SceneChanger.Draw(screen)
	}
	g.Scene.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.ScreenSize().X, g.ScreenSize().Y
}

// 수정 날짜는 정직하다고 가정
// 플레이 이외의 scene에서, 폴더 변화 감지하면 차트 리로드
// 로드된 차트 데이터는 gob에 별도 저장
func (g *Game) LoadCharts() error {
	for _, d := range dirs {
		for _, f := range files {
			switch filepath.Ext(f.Name()) {
			case ".osu":
				o, err := osu.Parse(path)
				if err != nil {
					panic(err) // todo: log and continue
				}
				switch o.Mode {
				case 3: // todo: osu.ModeMania
					c, err := mania.NewChartFromOsu()
					if err != nil {
						panic(err) // todo: log and continue
					}
					s.Charts = append(s.Charts, c)
				}
			}
		}
	}
	return nil
}
