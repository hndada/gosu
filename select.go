package gosu

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/mode/mania"
	"github.com/hndada/rg-parser/osugame/osu"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"io/ioutil"
	"path/filepath"
)

// todo: 로딩일 때 기다리는 로직
// Loading 이라는 별도의 Lock을 둔 이상, 특별히 채널은 필요없는거 아닌가?
// 비트맵 로딩 15초 후 timeout
type SceneSelect struct {
	g      *Game
	charts []chartPanel
	cursor int
	mods   mania.Mods
	// 그룹 (디렉토리 트리)
	// 현재 정렬 기준
	// Buttons     []ebitenui.Button
	// ChartPanels []ChartPanel
}

func (g *Game) NewSceneSelect() *SceneSelect {
	s := &SceneSelect{}
	s.g = g
	s.mods = mania.Mods{
		TimeRate: 1,
		Mirror:   false,
	}
	ebiten.SetWindowTitle("gosu")
	_ = s.checkCharts()
	return s
}

// 위쪽/왼쪽: 커서 -1
// 아래쪽/오른쪽: 커서 +1
// +시프트: 그룹 이동
func (s *SceneSelect) Init() {}
func (s *SceneSelect) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		s.g.SceneChanger.changeScene(s.g.NewSceneMania(s.charts[s.cursor].chart, s.mods))
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		s.cursor++
		if s.cursor >= len(s.charts) {
			s.cursor = 0
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		s.cursor--
		if s.cursor < 0 {
			s.cursor = len(s.charts) - 1
		}
	}
	screenSize := s.g.settings.ScreenSize()
	for i := range s.charts {
		mid := (screenSize.Y - 40) / 2 // 현재 선택된 차트 focus 틀 위치 고정
		s.charts[i].x = screenSize.X - 400
		s.charts[i].y = mid + 40*(i-s.cursor)
	}
	s.charts[s.cursor].x -= 30
	// for _, p := range s.ChartPanels {
	// 	p.Update()
	// }
	return nil
}

func (s *SceneSelect) Draw(screen *ebiten.Image) {
	for i, c := range s.charts {
		s.charts[i].op.GeoM.Reset()
		s.charts[i].op.GeoM.Translate(float64(c.x), float64(c.y))
		screen.DrawImage(c.box, s.charts[i].op)
	}
	// for _, p := range s.ChartPanels {
	// 	p.Draw(screen)
	// }
	// for _, b := range s.Buttons {
	// 	b.Draw(screen)
	// }
}

func (s *SceneSelect) checkCharts() error {
	// 폴더 변화 감지하면 LoadCharts()
	// 수정 날짜는 정직하다고 가정?
	return s.LoadCharts()
}

// 로드된 차트 데이터는 gob로 저장
func (s *SceneSelect) LoadCharts() error {
	s.charts = make([]chartPanel, 0, 100)
	dirs, err := ioutil.ReadDir(filepath.Join(s.g.path, "Music"))
	if err != nil {
		return err
	}
	for _, d := range dirs {
		if !d.IsDir() {
			continue
		}
		dpath := filepath.Join(s.g.path, "Music", d.Name())
		files, err := ioutil.ReadDir(dpath)
		if err != nil {
			return err
		}
		for _, f := range files {
			switch filepath.Ext(f.Name()) {
			case ".osu":
				fpath :=filepath.Join(dpath, f.Name())
				o, err := osu.Parse(fpath)
				if err != nil {
					panic(err) // todo: log and continue
				}
				switch o.Mode {
				case 3: // osu.ModeMania
					c, err := mania.NewChartFromOsu(o, fpath)
					if err != nil {
						panic(err) // todo: log and continue
					}
					s.charts = append(s.charts, newChartPanel(c))
				}
			}
		}
	}
	return nil
}

type chartPanel struct {
	box   *ebiten.Image
	x, y  int // todo: sprite-ize
	op    *ebiten.DrawImageOptions
	chart *mania.Chart
}

func newChartPanel(c *mania.Chart) chartPanel {
	var cp chartPanel
	img := image.NewRGBA(image.Rect(0, 0, 450, 40))
	col := color.RGBA{200, 100, 0, 255}
	x, y := 20, 30
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(fmt.Sprintf("%s [%s]", c.MusicName, c.ChartName))
	cp.box, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	cp.op = &ebiten.DrawImageOptions{}
	cp.chart = c
	return cp
}
