package gosu

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/mode/mania"
	"github.com/hndada/rg-parser/osugame/osu"
	"io/ioutil"
	"path/filepath"
)

// 차트 box:
// music name: chart name
// todo: 로딩일 때 기다리는 로직
// Loading 이라는 별도의 Lock을 둔 이상, 특별히 채널은 필요없는거 아닌가?
// 비트맵 로딩 15초 후 timeout
type SceneSelect struct {
	g       *Game
	mcharts []mania.Chart
	cursor  int
	mods    mania.Mods
	// 그룹 (디렉토리 트리)
	// 현재 정렬 기준
	// Buttons     []ebitenui.Button
	// ChartPanels []ChartPanel
}

func (g *Game) NewSceneSelect() *SceneSelect {
	s := &SceneSelect{}
	ebiten.SetWindowTitle("gosu")
	_ = s.checkCharts()
	return s
}

// 위쪽/왼쪽: 커서 -1
// 아래쪽/오른쪽: 커서 +1
// +시프트: 그룹 이동
func (s *SceneSelect) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		s.g.SceneChanger.changeScene(s.g.NewSceneMania(&s.mcharts[s.cursor], s.mods))
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		s.cursor++
		if s.cursor <= len(s.mcharts) {
			s.cursor = 0
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		s.cursor--
		if s.cursor < 0 {
			s.cursor = len(s.mcharts) - 1
		}
	}
	// for _, p := range s.ChartPanels {
	// 	p.Update()
	// }
	return nil
}

// 현재 선택된 차트 focus 틀 위치 고정
func (s *SceneSelect) Draw(screen *ebiten.Image) {
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
	s.mcharts = make([]mania.Chart, 0, 100)
	dirs, err := ioutil.ReadDir(filepath.Join(s.g.path, "Music"))
	if err != nil {
		return err
	}
	for _, d := range dirs {
		if !d.IsDir() {
			continue
		}
		dpath, err := filepath.Abs(d.Name())
		if err != nil {
			return err
		}
		files, err := ioutil.ReadDir(dpath)
		for _, f := range files {
			switch filepath.Ext(f.Name()) {
			case ".osu":
				fpath, err := filepath.Abs(f.Name())
				if err != nil {
					continue
				}
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
					s.mcharts = append(s.mcharts, *c)
				}
			}
		}
	}
	return nil
}
