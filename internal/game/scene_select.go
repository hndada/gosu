package game

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hndada/gosu/ebitenui"
	"github.com/hndada/gosu/mode/mania"
	"image"
	"image/color"

	"fmt"
)

// todo: 곡패널 만들기
// todo: Songs 폴더 읽는 로직 만들기 - rule 포함

// todo: SceneMania 네모 잘그리기
// todo: score, hp (race는 나중에)
// todo: 리절트창 고민
// todo: 입력, 스코어/HP 시뮬레이터 (기본->custom)

// 오디오 플레이어? // 필수는 아님
type SceneSelect struct {
	ChartPanels []ChartPanel
	cursor      int
	// 그룹 (디렉토리 트리)
	// 현재 정렬 기준
}

// 모든 box 생성?
// 현재 선택된 차트 focus (커서) 위치 고정

// 위쪽/왼쪽: 커서 -1
// 아래쪽/오른쪽: 커서 +1
// +시프트: 그룹 이동

// 폴더명: id or hash
// 1. id
// 2. 등록 안되어있으면 `--` 이어서 md5 앞 6자리, 16진수
// 2-1. 겹치는게 있다면 똑같이 6자리 하고 비교 / 나중 거는 자리수 추가 / 둘 다 자리수 추가 (얘는 어려울 듯)
// 웹에 없는 건 업데이트가 안되니 tracking이 안됨
func (s *SceneSelect) Update(g *Game) error {
	if ebiten.IsKeyPressed(ebiten.Key1) {
		c := &mania.Chart{}
		g.NextScene = NewSceneMania(g, c) // todo: go func()?
		g.TransCountdown = g.MaxTransCountDown()
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

func (s *SceneSelect) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "SceneSelect: Press Key 1")
	for _, p := range s.ChartPanels {
		p.Draw(screen)
	}
}

type ChartPanel struct {
	ChartPanelInfo
	ebitenui.Button
}

type ChartPanelInfo struct { // todo: mode-inspecific -> keys 분리 필요
	Keys        int
	Level       float64
	SongName    string // todo: Title -> SongName
	SongUnicode string
	ChartName   string
}

func (cp *ChartPanelInfo) Render() *ebiten.Image {
	const (
		twKeys  = 50 // X or 1X
		twLevel = 75 // X.X or XX.X, XXX, 999 이상은 999로.
		twName  = 400
		thPanel = 60
	)

	img := ebitenui.RenderBox(image.Pt(twKeys+twLevel+twName, thPanel),
		color.RGBA{117, 249, 102, 128}) // todo: skinnable
	op := &ebiten.DrawImageOptions{}

	keys := fmt.Sprintf("%2d", cp.Keys)
	img.DrawImage(ebitenui.RenderText(keys, ebitenui.MplusNormalFont, color.Black), op)
	op.GeoM.Translate(twKeys, 0)

	var lv string
	switch {
	case cp.Level < 100:
		lv = fmt.Sprintf("%2.1f", cp.Level)
	case cp.Level < 999:
		lv = fmt.Sprintf("%3d", int(cp.Level))
	default:
		lv = "999"
	}
	img.DrawImage(ebitenui.RenderText(lv, ebitenui.MplusNormalFont, color.Black), op)
	op.GeoM.Translate(twLevel, 0)

	name := fmt.Sprintf("%s\n%s", cp.SongName, cp.ChartName)
	fmt.Println(name)
	img.DrawImage(ebitenui.RenderText(name, ebitenui.MplusNormalFont, color.Black), op)
	return img
}

func NewChartPanel(c *mania.Chart, minPt image.Point) *ChartPanel {
	cp := &ChartPanel{}
	cp.ChartPanelInfo = ChartPanelInfo{ // todo: mode-inspecific
		Keys:        c.Keys,
		Level:       0, // todo: temp, c.Level
		SongName:    c.Title,
		SongUnicode: c.TitleUnicode,
		ChartName:   c.ChartName,
	}
	cp.MinPt = minPt
	cp.Image = cp.Render()
	cp.SetOnPressed(func(b *ebitenui.Button) {
		// SceneMania로 넘어가기
	})
	return cp
}
