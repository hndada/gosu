package gosu

import (
	"image"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/game/mania"
)

// anonymous struct: grouped globals
// reflect: fields should be exported
var argsSelectToMania struct {
	Chart      *mania.Chart
	Mods       mania.Mods
	ScreenSize image.Point
}

var sceneSelect *SceneSelect

type SceneSelect struct {
	game.Scene // includes ScreenSize
	cwd        string
	mods       mania.Mods
	// charts      []*mania.Chart // temp
	chartPanels []ChartPanel
	cursor      int
	holdCount   int

	ready bool
	done  bool

	playSE    func()
	defaultBG game.Sprite

	bornTime time.Time
}

func newSceneSelect(cwd string, size image.Point) *SceneSelect {
	s := new(SceneSelect)
	ebiten.SetWindowTitle("gosu")
	s.cwd = cwd
	s.mods = mania.Mods{
		TimeRate: 1,
		Mirror:   false,
	}

	updateCharts(cwd)
	s.chartPanels = make([]ChartPanel, len(charts))
	for i, c := range charts {
		s.chartPanels[i] = s.NewChartPanel(c)
	}
	s.ready = true
	s.ScreenSize = size

	s.playSE = mania.SEPlayer(cwd)
	s.defaultBG = game.DefaultBG()
	s.bornTime = time.Now()
	return s
}
func (s *SceneSelect) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		argsSelectToMania.Chart = charts[s.cursor]
		if argsSelectToMania.Chart.KeyCount == 8 { // temp
			argsSelectToMania.Mods.ScratchMode = mania.LeftScratch
		}
		argsSelectToMania.Mods = s.mods
		argsSelectToMania.ScreenSize = s.ScreenSize
		s.done = true
		s.holdCount = 0
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		if s.holdCount >= 2 { // todo: MaxTPS가 변하여도 체감 시간은 그대로이게 설정
			s.playSE()
			s.cursor++
			if s.cursor >= len(charts) {
				s.cursor = 0
			}
			s.holdCount = 0
		} else {
			s.holdCount++
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyUp) {
		if s.holdCount >= 2 {
			s.playSE()
			s.cursor--
			if s.cursor < 0 {
				s.cursor = len(charts) - 1
			}
			s.holdCount = 0
		} else {
			s.holdCount++
		}
	} else {
		s.holdCount = 0
	}

	for i := range s.chartPanels {
		mid := s.ScreenSize.Y / 2 // 현재 선택된 차트 focus 틀 위치 고정
		x := game.Settings.ScreenSize.X - 400
		d := i - s.cursor
		if d < 0 {
			d = -d
		}
		x += int(math.Pow(1.55, float64(d)))
		y := mid + 40*(i-s.cursor)
		if d == 0 {
			x -= 40
		}
		s.chartPanels[i].SetXY(x, y)
	}
	return nil
}

var bgs []*ebiten.Image

func (s *SceneSelect) Draw(screen *ebiten.Image) {
	//for i, cp := range s.chartPanels {
	//		if i == s.cursor {
	//screen.DrawImage(cp.BG, cp.OpBG)
	//break
	//	}
	//}
	s.defaultBG.Draw(screen)
	for _, cp := range s.chartPanels {
		cp.Draw(screen)
	}
}

func (s *SceneSelect) Ready() bool { return s.ready }
func (s *SceneSelect) Done(args *game.TransSceneArgs) bool {
	if s.done && args.Next == "" {
		args.Next = "mania.Scene"
		args.Args = argsSelectToMania // s.args
	}
	return s.done
}
