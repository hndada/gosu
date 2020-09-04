package gosu

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/graphics"
	"github.com/hndada/gosu/settings"
)

// 채보 재생: 속도, 로딩 및 싱크
// 스코어, hp 계산, 리플레이 저장
// 기타 Sprite 그리기
// 레벨 계산 대충 마무리

// 운동, 코딩, 잡일, 독서, 글씨 연습
// 다른 곳으로 날아가도, 외출 때 코딩을 하겠음, 아니면 과외나 학원 알바를 하든가

// 레벨 튜닝
// 시스템 디자인: pp(그대로 갈듯), 심플 웹, 랭크 시스템, 채보 discussion and contribution
// ui
// gosu만의 특별한 기능
// (다른 파일 포맷 파싱)

const Millisecond = 1000

type Game struct {
	path string
	settings.Settings
	graphics.GameSprites
	Scene        Scene
	SceneChanger *SceneChanger
	// Input        input.Input
}

// todo: 소리 재생
// Scene이 Game을 control하는 주체
type Scene interface {
	Init()
	Update() error
	Draw(screen *ebiten.Image) // Draws scene to screen
}

func NewGame() *Game {
	g := &Game{}
	// var err error
	// if g.path, err = os.Executable(); err != nil {
	// 	panic(err)
	// }
	g.path = `C:\Users\hndada\Documents\GitHub\hndada\gosu\test\`
	g.Settings.Load()
	g.GameSprites.Render(&g.Settings)
	g.Scene = g.NewSceneSelect()
	g.SceneChanger = g.NewSceneChanger()
	ebiten.SetMaxTPS(g.Settings.MaxTPS())
	ebiten.SetWindowTitle("gosu")
	ebiten.SetWindowSize(g.Settings.ScreenSize().X, g.Settings.ScreenSize().Y) // fixed in prototype
	ebiten.SetRunnableOnUnfocused(true)
	return g
}
func (g *Game) Update(screen *ebiten.Image) error {
	if !g.SceneChanger.done() {
		return g.SceneChanger.Update()
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

func (g *Game) changeScene(s Scene) {
	g.SceneChanger.changeScene(s)
}

// todo: float64, fixed로 고치기 생각

// reset    save cancel
// 설정 켜면 임시 세팅이 생성, 임시 세팅으로 실시간 보여주기
// save 누르면 실제 세팅으로 값복사
// game에서 세팅 바꾸면 Sprite 자동 갱신

// PlayIntro, PlayExit
