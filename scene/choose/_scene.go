package choose

import (
	"fmt"
	"io/fs"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/mode/drum"
	"github.com/hndada/gosu/mode/piano"
	"github.com/hndada/gosu/scene"
)

type Scene struct {
	// musicCh      chan []byte
	// bgCh       chan draws.Image
	// returnCh chan Return
	page        int
	loading     bool
	keySettings []string //[]input.Key
	chartSets   []ChartSet
	charts      []*Chart
	ChartSets   *List
	Charts      *List
}

func NewScene() *Scene {
	s := &Scene{}

	// s.lastChartSets = make([][]*ChartSet, 3)
	// for i := range s.lastChartSets {
	// 	css, err := search("", i, 0)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	s.lastChartSets[i] = css
	// }

	// chartSets := make([]string, 0)
	// dirs, err := os.ReadDir("./music")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// for _, dir := range dirs {
	// 	if dir.IsDir() {
	// 		chartSets = append(chartSets, dir.Name())
	// 	}
	// }
	s.chartSets = LoadChartSets()
	{
		rows := make([]string, 0, len(s.chartSets))
		for _, cs := range s.chartSets {
			row := fmt.Sprintf("%s - %s", cs.Artist, cs.Title)
			rows = append(rows, row)
		}
		s.ChartSets = NewList(rows)
	}
	s.Focus = FocusChartSet
	return s
}

func (s *Scene) Update() any {
	// s.Preview.Update()
	// select {
	// case i := <-s.bgCh:
	// 	sprite := draws.NewSprite(i)
	// 	sprite.SetScaleToW(ScreenSizeX)
	// 	sprite.Locate(ScreenSizeX/2, ScreenSizeY/2, draws.CenterMiddle)
	// 	s.Background.Sprite = sprite
	// default:
	// }
	// if s.query != s.Query.Text {
	// 	s.query = s.Query.Text
	// 	s.Focus = FocusSearch
	// }

	if isEnter() {
		switch s.Focus {
		case FocusSearch:
			fmt.Println("Load chart sets")
			s.query = s.Query.Text
			s.LoadChartSetList()
		case FocusChartSet:
			fmt.Println("Load chart")
			idx := s.ChartSets.Current()
			s.charts = s.chartSets[idx].ChildrenBeatmaps
			{
				rows := make([]string, 0, len(s.charts))
				for _, cs := range s.charts {
					row := fmt.Sprintf("%s [%s]", cs.Title, cs.DiffName)
					rows = append(rows, row)
				}
				s.Charts = NewList(rows)
			}
			// base := s.ChartSets.Current()
			// charts := make([]string, 0)
			// files, err := os.ReadDir("./music/" + s.ChartSets.Current())
			// if err != nil {
			// 	fmt.Println("uhh")
			// 	return err
			// }
			// for _, file := range files {
			// 	if file.IsDir() {
			// 		continue
			// 	}
			// 	// check whether file is a chart
			// 	if !strings.HasSuffix(file.Name(), ".osu") {
			// 		continue
			// 	}
			// 	charts = append(charts, base+file.Name())
			// }
			// s.Charts = NewList(charts)
			s.Focus = FocusChart
		case FocusChart:
			fmt.Println("Play chart")
			go func() {
				s.loading = true
				scene.UserSkin.Enter.Play(*s.volumeSound)
				// if c == nil {
				// 	fmt.Println(errors.New("no chart loaded"))
				// 	return
				// }
				// fs, name, err := c.Choose()
				// if err != nil {
				// 	fmt.Println(err)
				// }
				// s.Preview.Close()
				sub, err := fs.Sub(s.MusicFS, s.chartSets[s.ChartSets.Current()].Path)
				if err != nil {
					fmt.Println(err)
				}
				// csi := s.ChartSets.Current()
				// cs:= s.chartSets[csi]
				idx := s.Charts.Current()
				name := s.charts[idx].OsuFile
				s.r = &Return{
					// FS:     fs,
					// Name:   name,
					FS:     sub,
					Name:   name, // suppose it contains a whole path
					Mode:   modes[s.mode],
					Mods:   nil,
					Replay: nil,
				}
				s.loading = false
			}()
		}
		return nil
	}
	if isBack() {
		switch s.Focus {
		// case FocusSearch, FocusChartSet:
		// 	s.Query.Reset()
		// 	s.Focus = FocusSearch
		case FocusChart:
			s.Focus = FocusChartSet
		case FocusKeySettings:
			s.Focus = s.lastFocus
		}
	}
	// if s.Mode.Update() { // || s.SubMode.Update()
	// 	go s.LoadChartSetList()
	// }
	switch s.Focus {
	// case FocusSearch:
	// 	s.Query.Update()
	case FocusChartSet:
		if !s.ChartSets.Update() {
			break
		}
		// cset := s.ChartSets.Current()
		// go func() {
		// 	// i, err := draws.NewImageFromURL("https://upload.wikimedia.org/wikipedia/commons/1/1f/As08-16-2593.jpg")
		// 	i, err := draws.NewImageFromURL(cset.URLCover("cover", Large))
		// 	if err != nil {
		// 		return
		// 	}
		// 	s.Background.Sprite.Source = i
		// 	// s.bgCh <- draws.Image{Image: i}
		// }()
	case FocusChart:
		s.Charts.Update()
	case FocusKeySettings:
	}
	return nil
}
