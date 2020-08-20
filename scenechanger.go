package gosu

import "github.com/hajimehoshi/ebiten"

// the word changer is more descriptive than manager
type SceneChanger struct {
	Scene          Scene
	NextScene      Scene
	TransSceneFrom *ebiten.Image
	TransSceneTo   *ebiten.Image
	TransCountdown int
}

func (sc *SceneChanger) Update(g *Game) error {
	if sc.done() {
		return g.Scene.Update(g)
	}
	sc.TransCountdown--
	if sc.TransCountdown > 0 {
		return nil
	}
	// count down has just been from non-zero to zero
	g.Scene = sc.NextScene
	sc.NextScene = nil
	return nil
}

func (sc *SceneChanger) Draw(screen *ebiten.Image) {
	if sc.done() {
		return
	}
	var value float64
	{
		value = float64(sc.TransCountdown) / float64(sc.MaxTransCountDown())
		sc.TransSceneFrom.Clear()
		sc.Scene.Draw(sc.TransSceneFrom)
		op := ebiten.DrawImageOptions{}
		op.ColorM.ChangeHSV(0, 1, value)
		screen.DrawImage(sc.TransSceneFrom, &op)
	}
	{
		value = 1 - value
		sc.TransSceneTo.Clear()
		sc.NextScene.Draw(sc.TransSceneTo)
		op := ebiten.DrawImageOptions{}
		op.ColorM.ChangeHSV(0, 1, value)
		screen.DrawImage(sc.TransSceneTo, &op)
	}
}

// The function is called every time when settings has been updated
func NewSceneChanger() *SceneChanger {
	sc := &SceneChanger{}
	sc.TransSceneFrom, _ = ebiten.NewImage(g.ScreenWidth, g.ScreenHeight, ebiten.FilterDefault)
	sc.TransSceneTo, _ = ebiten.NewImage(g.ScreenWidth, g.ScreenHeight, ebiten.FilterDefault)
}

func (g *Game) ChangeScene(s Scene) {
	g.NextScene = s
	g.TransCountdown = g.MaxTransCountDown()
}

// 모든 time 관련 단위는 ms
func (sc *SceneChanger) MaxTransCountDown() int {
	return int(0.8 * float64(g.MaxTPS))
}

func (sc *SceneChanger) done() bool { return sc.TransCountdown == 0 }
