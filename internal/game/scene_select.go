package game

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hndada/gosu/graphics"
	"github.com/hndada/gosu/mode/mania"
	"image/color"

	"fmt"
)

// 오디오 플레이어? // 필수는 아님
type SceneSelect struct {
	// 차트 리스트
	cursor int
	// 그룹 (디렉토리 트리)
	// 현재 정렬 기준
}

// 모든 box 생성?
// 현재 선택된 차트 focus (커서) 위치 고정

// 위쪽/왼쪽: 커서 -1
// 아래쪽/오른쪽: 커서 +1
// +시프트: 그룹 이동
func (s *SceneSelect) Update(g *Game) error {
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
	if ebiten.IsKeyPressed(ebiten.Key1) {
		c := &mania.Chart{}
		g.NextScene = NewSceneMania(g, c) // todo: go func()?
		g.TransCountdown = g.MaxTransCountDown()
	}
	return nil
}

func (s *SceneSelect) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "SceneSelect: Press Key 1")
}

type ChartPanel struct { // todo: mode-inspecific -> keys 분리 필요
	Keys        int
	Level       float64
	SongName    string // todo: Title -> SongName
	SongUnicode string
	ChartName   string
}

func NewChartPanel(c *mania.Chart) *ChartPanel { // todo: mode-inspecific
	return &ChartPanel{
		Keys:        c.Keys,
		Level:       0, // todo: temp, c.Level
		SongName:    c.Title,
		SongUnicode: c.TitleUnicode,
		ChartName:   c.ChartName,
	}
}

const (
	textWidthKeys   = 25 // X or 1X
	textWidthLevel  = 75 // X.X or XX.X, XXX, 999 이상은 999로.
	textWidthName   = 400
	textHeightPanel = 60
)

var (
	chartPanelColor = color.RGBA{166, 233, 240, 128}
)

func DrawChartPanel(cp *ChartPanel) *ebiten.Image {
	keys := fmt.Sprintf("%2d", cp.Keys)
	a := graphics.DrawTextBox(
		graphics.DrawText(keys, graphics.FontVarelaNormal, color.Black),
		chartPanelColor)

	var lv string
	switch {
	case cp.Level < 100:
		lv = fmt.Sprintf("%2.1f", cp.Level)
	case cp.Level < 999:
		lv = fmt.Sprintf("%3d", int(cp.Level))
	default:
		lv = "999"
	}
	b := graphics.DrawTextBox(
		graphics.DrawText(lv, graphics.FontVarelaNormal, color.Black), // todo: various font
		chartPanelColor)

	name := fmt.Sprintf(`%s\n%s`, cp.SongName, cp.ChartName)
	c := graphics.DrawTextBox(
		graphics.DrawText(name, graphics.FontVarelaNormal, color.Black),
		chartPanelColor)
	return graphics.AttachH(a, b, c)
}
