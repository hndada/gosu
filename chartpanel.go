package gosu

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/ebitenui"
	"github.com/hndada/gosu/mode/mania"
	"image"
	"image/color"
)

type ChartPanel struct {
	// 차트 정보
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
		SongName:    c.SongName,
		SongUnicode: c.SongUnicode,
		ChartName:   c.ChartName,
	}
	cp.MinPt = minPt
	cp.Image = cp.Render()
	cp.SetOnPressed(func(b *ebitenui.Button) {
		// SceneMania로 넘어가기
	})
	return cp
}
