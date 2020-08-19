package gosu

import (
	"github.com/hajimehoshi/ebiten"
	input2 "github.com/hndada/gosu/input"
	"github.com/hndada/gosu/mode/mania"
	"github.com/hndada/rg-parser/osugame/osu"
	"path/filepath"
)

// 폴더명: id or hash
// 1. id
// 2. 등록 안되어있으면 `--` 이어서 md5 앞 6자리, 16진수
// 2-1. 겹치는게 있다면 똑같이 6자리 하고 비교 / 나중 거는 자리수 추가 / 둘 다 자리수 추가 (얘는 어려울 듯)
// 웹에 없는 건 업데이트가 안되니 tracking이 안됨
// mp3 플레이어, Scene에 저장; 연동
type Game struct {
	State
	Options
}

type State struct {
	Scene          Scene
	NextScene      Scene
	TransSceneFrom *ebiten.Image
	TransSceneTo   *ebiten.Image
	TransCountdown int
	Loading        bool

	input input2.Input
}

type Options struct {
	ScrollSpeed  float64
	KeysLayout   map[int][]ebiten.Key
	MaxTPS       int
	ScreenWidth  int
	ScreenHeight int
	HitPosition  float64 // object which is now set at 'options'

	DimValue  int
	VolumeSFX int
	VolumeBGM int
	Skin      *Skin
}

// scene이 game을 control하는 주체
type Scene interface {
	Update(g *Game) error
	Draw(screen *ebiten.Image)
}

func (g *Game) MaxTransCountDown() int { return int(0.8 * float64(g.MaxTPS)) } // 모든 time 관련 단위는 ms
func (g *Game) Update(screen *ebiten.Image) error {
	if g.Loading {
		return nil
	}
	if g.TransCountdown == 0 {
		return g.Scene.Update(g)
	}
	g.TransCountdown--
	if g.TransCountdown > 0 {
		return nil
	}
	// count down has just been from non-zero to zero
	g.Scene = g.NextScene
	g.NextScene = nil
	return nil
}

// scene의 Draw는 input으로 들어온 screen을 그리는 함수
// Game.Draw() 자체에서는 screen에 직접 그려봤자 반영 안됨
func (g *Game) Draw(screen *ebiten.Image) {
	if g.Loading {
		return
	}
	if g.TransCountdown == 0 {
		g.Scene.Draw(screen)
		return
	}
	var value float64
	{
		value = float64(g.TransCountdown) / float64(g.MaxTransCountDown())
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
	return g.ScreenWidth, g.ScreenHeight
}

func NewGame() (g *Game) {
	g = &Game{}
	g.State = State{
		Scene: &SceneTitle{},
	}
	g.Options = Options{ // todo: load settings
		MaxTPS:      240,
		ScrollSpeed: 1.33,
		HitPosition: 730,
	}
	g.ScreenWidth = 1600
	g.ScreenHeight = 900
	g.TransSceneFrom, _ = ebiten.NewImage(g.ScreenWidth, g.ScreenHeight, ebiten.FilterDefault)
	g.TransSceneTo, _ = ebiten.NewImage(g.ScreenWidth, g.ScreenHeight, ebiten.FilterDefault)
	return
}

// option에 관련 세팅이 들어갈 수 있을 것 같아 game의 method로
func (g *Game) PlayIntro() {

}
func (g *Game) PlayExit() {

}

func (g *Game) ChangeScene(s Scene) {
	g.NextScene = s
	g.TransCountdown = g.MaxTransCountDown()
}

// 수정 날짜는 정직하다고 가정
// 플레이 이외의 scene에서, 폴더 변화 감지하면 차트 리로드
// 로드된 차트 데이터는 게임 폴더의 db파일에 별도 저장
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
