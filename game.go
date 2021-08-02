package gosu

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/game/mania"

	// _ "github.com/silbinarywolf/preferdiscretegpu"
	"path/filepath"
)

// 스코어, hp 고치기
// internal 삭제
// pp
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
	// g.path = `C:\Users\hndada\Documents\GitHub\hndada\gosu\test\`
	g.path = `F:\projects\gosu\cmd\gosu`

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
