package gosu

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/game/mania"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// anonymous struct: grouped globals
// reflect: fields should be exported
var argsSelectToMania struct {
	Chart      *mania.Chart
	Mods       mania.Mods
	ScreenSize image.Point
}

type sceneSelect struct {
	game.Scene // includes ScreenSize
	path       string
	mods       mania.Mods
	charts     []chartPanel
	cursor     int
	holdCount  int

	ready bool
	done  bool
}

func newSceneSelect(path string, size image.Point) *sceneSelect {
	s := new(sceneSelect)
	ebiten.SetWindowTitle("gosu")
	s.path = path
	s.mods = mania.Mods{
		TimeRate: 1,
		Mirror:   false,
	}
	_ = s.checkCharts()
	s.ready = true
	s.ScreenSize = size
	return s
}

func (s *sceneSelect) Ready() bool { return s.ready }
func (s *sceneSelect) Done(args *game.TransSceneArgs) bool {
	if s.done && args.Next == "" {
		args.Next = "mania.Scene"
		args.Args = argsSelectToMania // s.args
	}
	return s.done
}
func (s *sceneSelect) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		argsSelectToMania.Chart = s.charts[s.cursor].chart
		argsSelectToMania.Mods = s.mods
		argsSelectToMania.ScreenSize = s.ScreenSize
		s.done = true
		s.holdCount = 0
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		if s.holdCount >= 2 { // todo: MaxTPS가 변하여도 체감 시간은 그대로이게 설정
			s.cursor++
			if s.cursor >= len(s.charts) {
				s.cursor = 0
			}
			s.holdCount = 0
		} else {
			s.holdCount++
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyUp) {
		if s.holdCount >= 2 {
			s.cursor--
			if s.cursor < 0 {
				s.cursor = len(s.charts) - 1
			}
			s.holdCount = 0
		} else {
			s.holdCount++
		}
	} else {
		s.holdCount = 0
	}

	for i := range s.charts {
		mid := (s.ScreenSize.Y - 40) / 2 // 현재 선택된 차트 focus 틀 위치 고정
		s.charts[i].x = s.ScreenSize.X - 400
		s.charts[i].y = mid + 40*(i-s.cursor)
	}
	s.charts[s.cursor].x -= 30
	return nil
}

func (s *sceneSelect) Draw(screen *ebiten.Image) {
	for i, c := range s.charts {
		s.charts[i].op.GeoM.Reset()
		s.charts[i].op.GeoM.Translate(float64(c.x), float64(c.y))
		screen.DrawImage(c.box, s.charts[i].op)
	}
}

func (s *sceneSelect) checkCharts() error {
	return s.LoadCharts()
}

// 로드된 차트 데이터는 gob로 저장
func (s *sceneSelect) LoadCharts() error {
	s.charts = make([]chartPanel, 0, 100)
	dirs, err := ioutil.ReadDir(filepath.Join(s.path, "Music"))
	if err != nil {
		return err
	}
	for _, d := range dirs {
		if !d.IsDir() {
			continue
		}
		dpath := filepath.Join(s.path, "Music", d.Name())
		files, err := ioutil.ReadDir(dpath)
		if err != nil {
			return err
		}
		for _, f := range files {
			fpath := filepath.Join(dpath, f.Name())
			switch strings.ToLower(filepath.Ext(fpath)) {
			case ".osu":
				switch game.OsuMode(fpath) {
				case game.ModeMania:
					c, err := mania.NewChart(fpath)
					if err != nil {
						panic(err) // todo: log and continue
					}
					s.charts = append(s.charts, newChartPanel(c))
				}
			}
		}
	}
	sort.Slice(s.charts, func(i, j int) bool {
		if s.charts[i].chart.Keys == s.charts[j].chart.Keys {
			return s.charts[i].chart.Level < s.charts[j].chart.Level
		} else {
			return s.charts[i].chart.Keys < s.charts[j].chart.Keys
		}
	})
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
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(fmt.Sprintf("(%dKey Lv %.1f) %s [%s]", c.Keys, c.Level, c.MusicName, c.ChartName))
	cp.box, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	cp.op = &ebiten.DrawImageOptions{}
	cp.chart = c
	return cp
}

// var ModePrefix = map[int]string{0: "o", 1: "t", 2: "c", 3: "m"}
// 폴더 두번 스캔: mapsets, maps
func LoadSongList(root, ext string) ([]string, error) {
	var songs []string
	if info, err := os.Stat(root); err != nil || !info.IsDir() {
		return songs, errors.New("invalid root dir")
	}
	sets, err := ioutil.ReadDir(root)
	if err != nil {
		return songs, errors.New("invalid root dir")
	}

	var absSet, absMap string
	for _, set := range sets {
		absSet = filepath.Join(root, set.Name())
		if info, err := os.Stat(absSet); err != nil || !info.IsDir() {
			continue
		}
		maps, err := ioutil.ReadDir(absSet)
		if err != nil {
			continue
		}
		for _, mapFile := range maps {
			absMap = filepath.Join(absSet, mapFile.Name())
			if !mapFile.IsDir() && filepath.Ext(mapFile.Name()) == ext {
				songs = append(songs, absMap)
			}
		}
	}
	return songs, nil
}
