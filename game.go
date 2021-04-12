package gosu

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/game/mania"

	// _ "github.com/silbinarywolf/preferdiscretegpu"
	"path/filepath"
)

// select에서 곡 로드 방식 바꾸기, 그 외 간단히 고칠 것들 todo 나 메모에서
// 리플레이 저장, 이후 스코어/hp 고치기
// 기타 Sprite 그리기
// docs, markdown 정리 및 internal 삭제

// 운동, 코딩, 잡일, 독서, 글씨 연습
// 다른 곳으로 날아가도, 외출 때 코딩을 하겠음, 아니면 과외나 학원 알바를 하든가

// 레벨 튜닝
// 시스템 디자인: pp(그대로 갈듯), 심플 웹, 랭크 시스템, 채보 discussion and contribution
// ui
// gosu만의 특별한 기능: 다른 파일 포맷 파싱 등
// game programming beginner's guide with go?
var MaxTransCountDown int

type Game struct {
	path           string
	Scene          Scene
	NextScene      Scene
	TransSceneFrom *ebiten.Image
	TransSceneTo   *ebiten.Image
	TransCountdown int
}

// 실제 파일을 자주 불러오는듯
// skin: *ebiten.Image, raw 이미지들
// sprites: *ebiten.Image, 크기랑 position 맞춰진 이미지들

// todo: Load, Save settings
// type settings struct {
// 	common game.SettingsTemplate
// 	mania  mania.SettingsTemplate
// }
type Scene interface { // Scene이 Game을 control하는 주체
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

	game.LoadSettings()
	mania.ResetSettings()
	mania.LoadSpriteMap(filepath.Join(g.path, "Skin"))
	g.Scene = g.NewSceneSelect()

	// g.SceneChanger = g.NewSceneChanger()
	p := game.ScreenSize()
	g.TransSceneFrom, _ = ebiten.NewImage(p.X, p.Y, ebiten.FilterDefault)
	g.TransSceneTo, _ = ebiten.NewImage(p.X, p.Y, ebiten.FilterDefault)
	MaxTransCountDown = game.MaxTPS() * 4 / 5

	ebiten.SetWindowTitle("gosu")
	ebiten.SetRunnableOnUnfocused(true)
	return g
}

func (g *Game) Update(screen *ebiten.Image) error {
	if g.TransCountdown <= 0 { // == 0
		return g.Scene.Update()
	}
	g.TransCountdown--
	if g.TransCountdown > 0 {
		return nil
	}
	// count down has just been from non-zero to zero
	g.Scene = g.NextScene
	g.NextScene = nil
	g.Scene.Init()
	return nil
	// if !g.SceneChanger.done() {
	// 	return g.SceneChanger.Update()
	// }
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.TransCountdown == 0 {
		g.Scene.Draw(screen)
		return
	}
	var value float64
	{
		value = float64(g.TransCountdown) / float64(MaxTransCountDown)
		g.TransSceneFrom.Clear()
		g.Scene.Draw(g.TransSceneFrom)
		op := ebiten.DrawImageOptions{}
		op.ColorM.ChangeHSV(0, 1, value)
		screen.DrawImage(g.TransSceneFrom, &op)
	}
	{
		value = 1 - value
		g.TransSceneTo.Clear()
		g.NextScene.Draw(g.TransSceneTo)
		op := ebiten.DrawImageOptions{}
		op.ColorM.ChangeHSV(0, 1, value)
		screen.DrawImage(g.TransSceneTo, &op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return game.ScreenSize().X, game.ScreenSize().Y
}

func (g *Game) changeScene(s Scene) {
	g.NextScene = s
	g.TransCountdown = MaxTransCountDown
}

// reset    save cancel
// 설정 켜면 임시 세팅이 생성, 임시 세팅으로 실시간 보여주기
// save 누르면 실제 세팅으로 값복사
// game에서 세팅 바꾸면 Sprite 자동 갱신

// PlayIntro, PlayExit
// BasePlayScene (base struct), PlayScene (interface)
// PlayScene, 각 mode 패키지에다가 구현해야 할까?

// float64, fixed로 고칠건 따로 안 보임
