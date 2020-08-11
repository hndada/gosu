package game

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hndada/gosu/mode/mania"
	"image"
	"image/color"
	_ "image/jpeg"
	"io/ioutil"
)

type SceneMania struct { // aka Clavier
	Notes        []NoteImageInfo
	C            mania.Chart
	TickDuration float64
	scrollSpeed  float64
	Time         float64
	Score        float64
	HP           float64

	Cover *ebiten.Image
	// Audio
}

type NoteImageInfo struct {
	x, y, w, h float64
	clr        color.RGBA
}

// todo: 맵 로딩하는동안 업데이트 말고 로딩에만 집중하게
func (s *SceneMania) Update(g *Game) error {
	const endTime = 5.0 * 1000 // float64(s.C.Notes[len(s.C.Notes)-1].Time)
	s.Time += s.TickDuration
	for i := range s.Notes {
		// if s.Notes[i].y > 900 { continue }
		s.Notes[i].y += s.TickDuration * s.scrollSpeed
	}
	if s.Time > endTime {
		g.NextScene = &SceneResult{}
		g.TransCountdown = 99 // todo: set maxCount
	}
	return nil
}

// todo:범위 넘어간 애들은 Rect 안그리기 -> 오히려 fps 불안정
func (s *SceneMania) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	const ratio = float64(1600) / 1920
	op.GeoM.Scale(ratio, ratio)
	op.ColorM.ChangeHSV(0, 1, 0.30)

	screen.DrawImage(s.Cover, op)

	ebitenutil.DrawRect(screen, 565, 0, 70*7, 900, color.RGBA{0, 0, 0, 180})
	ebitenutil.DrawRect(screen, 565, 730, 70*7, 10, color.RGBA{252, 106, 111, 255})
	for i := range s.Notes {
		// if s.Notes[i].y>900 { continue }
		// if s.Notes[i].y+s.Notes[i].h<0 { continue }
		ebitenutil.DrawRect(screen, s.Notes[i].x, s.Notes[i].y, s.Notes[i].w, s.Notes[i].h, s.Notes[i].clr)
	}
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f\nTime: %.1fs", ebiten.CurrentFPS(), s.Time/1000))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\nTime: %.1fs", s.Time/1000))
}

func noteColor(n mania.Note, keys int) color.RGBA {
	switch n.Key {
	case 0, 2, 4, 6:
		return color.RGBA{239, 243, 247, 0xff} // white
	case 1, 5:
		return color.RGBA{66, 211, 247, 0xff} // blue
	case 3:
		return color.RGBA{255, 203, 82, 0xff} // yellow
	}
	panic("not reach")
}

// 노트, 그냥 네모 그리고 색깔 채워넣기
// &SceneMania{}로 하고 chart 로딩을 할까
func NewSceneMania(op Options, c *mania.Chart) (s *SceneMania) {
	ebiten.SetWindowTitle(fmt.Sprintf("gosu - %s [%s]", c.Title, c.ChartName)) // todo: can I change window title ?
	s = &SceneMania{}
	const w = 70
	const noteHeight = 25
	s.C = *c
	s.Notes = make([]NoteImageInfo, len(c.Notes))
	s.scrollSpeed = op.ScrollSpeed
	s.TickDuration = 1000 / float64(op.MaxTPS) // 스피드값 1 기준 초당 1000픽셀 내려오게 해야함
	for i, n := range c.Notes {
		var y, h float64
		x := float64(n.Key*w + 565)
		switch n.Type {
		case mania.TypeNote:
			y = -float64(n.Time)*op.ScrollSpeed + 1000
			h = noteHeight * op.ScrollSpeed
		case mania.LNHead:
			y = -float64(n.Time2)*op.ScrollSpeed + 1000
			h = float64(n.Time2-n.Time+noteHeight) * op.ScrollSpeed
		}
		s.Notes[i] = NoteImageInfo{x, y, w, h, noteColor(n, c.Keys)}
	}
	b, err := ioutil.ReadFile("C:\\Users\\hndada\\Documents\\GitHub\\hndada\\gosu\\mode\\mania\\test\\" + c.ImageFilename)
	if err != nil {
		panic(err)
	}
	cover, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	s.Cover, _ = ebiten.NewImageFromImage(cover, ebiten.FilterDefault)
	return
}
