package choose

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math"
	"net/http"

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

const (
	TPS         = scene.TPS
	ScreenSizeX = scene.ScreenSizeX
	ScreenSizeY = scene.ScreenSizeY
)

var modes = []int{piano.Mode, piano.Mode, drum.Mode}

// Background brightness at Song select: 60% (153 / 255), confirmed.
// Score box color: Gray128 with 50% transparent
// Hovered Score box color: Gray96 with 50% transparent

// Todo: fetch Score with Replay
// Todo: preview music. Rewind after preview has finished.
// Group1, Group2, Sort, Filter int
type Scene struct {
	volumeMusic   *float64
	volumeSound   *float64
	brightness    *float64
	offset        *int64
	speedFactors  []*float64
	exposureTimes []func(float64) float64

	// choose audios.Sound
	// musicCh      chan []byte
	// Music      audios.MusicPlayer
	bgCh       chan draws.Image
	Background mode.BackgroundDrawer

	mode int
	Mode ctrl.KeyHandler
	// subMode int
	// SubMode ctrl.KeyHandler
	query string
	Query input.TypeWriter
	page  int

	Focus     int
	ChartSets ChartSetList
	Charts    ChartList
	// lastChartSets [][]*ChartSet
	levelLimit bool
	LevelLimit ctrl.KeyHandler

	loading bool
	Loading LoadingDrawer
	// inited  bool

	Preview PreviewPlayer
}

const (
	FocusSearch = iota
	FocusChartSet
	FocusChart
)

func NewScene() *Scene {
	s := &Scene{}
	s.volumeMusic = &mode.S.VolumeMusic
	s.volumeSound = &mode.S.VolumeSound
	s.brightness = &mode.S.BackgroundBrightness
	s.offset = &mode.S.Offset
	s.speedFactors = []*float64{
		&piano.S.SpeedScale, &drum.S.SpeedScale}
	s.exposureTimes = []func(float64) float64{
		piano.ExposureTime, drum.ExposureTime}
	// s.Mode = ctrl.KeyHandler{
	// 	Handler: ctrl.IntHandler{
	// 		Value: &s.mode,
	// 		Min:   0,
	// 		Max:   3 - 1, // There are 3 modes.
	// 		Loop:  true,
	// 	},
	// 	Modifiers: []input.Key{},
	// 	Keys:      [2]input.Key{-1, input.KeyF1},
	// 	Sounds:    [2]audios.Sounder{scene.UserSkin.Swipe, scene.UserSkin.Swipe},
	// 	Volume:    &mode.S.VolumeSound,
	// }
	// s.subMode = 4
	// s.SubMode = ctrl.KeyHandler{
	// 	Handler: ctrl.IntHandler{
	// 		Value: &s.subMode,
	// 		Min:   4,
	// 		Max:   9,
	// 		Loop:  true,
	// 	},
	// 	Modifiers: []input.Key{},
	// 	Keys:      [2]input.Key{input.KeyF2, input.KeyF3},
	// 	Sounds:    [2]audios.Sounder{scene.UserSkin.Swipe, scene.UserSkin.Swipe},
	// 	Volume:    &mode.S.VolumeSound,
	// }

	// s.lastChartSets = make([][]*ChartSet, 3)
	// for i := range s.lastChartSets {
	// 	css, err := search("", i, 0)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	s.lastChartSets[i] = css
	// }
	s.levelLimit = true
	s.Loading = NewLoadingDrawer()
	s.handleEnter()
	// ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)
	// debug.SetGCPercent(100)
	ebiten.SetWindowTitle("gosu")
	return s
}
func isEnter() bool {
	return inpututil.IsKeyJustPressed(input.KeyEnter) ||
		inpututil.IsKeyJustPressed(input.KeyNumpadEnter)
}
func isBack() bool {
	return inpututil.IsKeyJustPressed(input.KeyEscape)
}

func (s *Scene) Update() any {
	if s.loading {
		s.Loading.Update()
		return nil
	}
	scene.VolumeMusic.Update()
	scene.VolumeSound.Update()
	scene.Brightness.Update()
	scene.Offset.Update()
	scene.SpeedScales[modes[s.mode]].Update()
	s.Preview.Update()
	select {
	case i := <-s.bgCh:
		sprite := draws.NewSpriteFromSource(i)
		sprite.SetScaleToW(ScreenSizeX)
		sprite.Locate(ScreenSizeX/2, ScreenSizeY/2, draws.CenterMiddle)
		s.Background.Sprite = sprite
	default:
	}
	if s.query != s.Query.Text {
		s.query = s.Query.Text
		s.Focus = FocusSearch
	}
	if inpututil.IsKeyJustPressed(input.KeyF1) {
		s.mode++
		s.mode %= 3
		scene.UserSkin.Enter.Play(*s.volumeSound)
		s.Focus = FocusSearch
		err := s.handleEnter()
		if err != nil {
			fmt.Println(err)
		}
	}
	if inpututil.IsKeyJustPressed(input.KeyF4) {
		s.levelLimit = !s.levelLimit
		scene.UserSkin.Swipe.Play(*s.volumeSound)
	}
	if isEnter() {
		return s.handleEnter()
	}
	if isBack() {
		switch s.Focus {
		case FocusSearch, FocusChartSet:
			s.Query.Reset()
			s.Focus = FocusSearch
		case FocusChart:
			s.Focus = FocusChartSet
		}
	}
	if s.Mode.Update() { // || s.SubMode.Update()
		go s.LoadChartSetList()
	}
	switch s.Focus {
	case FocusSearch:
		s.Query.Update()
	case FocusChartSet:
		fired, state := s.ChartSets.Update()
		if !fired {
			break
		}
		switch state {
		case prev:
			if s.page == 0 {
				break
			}
			s.page--
			go func() {
				s.LoadChartSetList()
				s.ChartSets.cursor = RowCount - 1
			}()
		case next:
			css := s.ChartSets
			s.page++
			go func() {
				s.LoadChartSetList()
				if len(s.ChartSets.ChartSets) == 0 {
					s.ChartSets = css
					s.page--
				}
			}()
		case stay:
		}
		cset := s.ChartSets.Current()
		go func() {
			resp, err := http.Get(cset.URLPreview())
			if err != nil || resp.StatusCode == 404 {
				fmt.Println(err)
				return
			}
			defer resp.Body.Close()
			if s.Preview.IsValid() {
				s.Preview.Close()
			}
			s.Preview, err = NewPreviewPlayer(resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}
		}()
		go func() {
			i, err := ebitenutil.NewImageFromURL(cset.URLCover("cover", Large))
			if err != nil {
				return
			}
			s.bgCh <- draws.Image{Image: i}
		}()

	case FocusChart:
		s.Charts.Update()
	}
	return nil
}
func (s *Scene) handleEnter() any {
	switch s.Focus {
	case FocusSearch:
		fmt.Println("Load chart sets")
		s.query = s.Query.Text
		s.LoadChartSetList()
	case FocusChartSet:
		fmt.Println("Load chart")
		s.LoadChartList()
	case FocusChart:
		fmt.Println("Play chart")
		scene.UserSkin.Enter.Play(*s.volumeSound)
		c := s.Charts.Current()
		if c == nil {
			return errors.New("no chart loaded")
		}
		fs, name, err := c.Choose()
		if err != nil {
			return err
		}
		return Return{
			FS:     fs,
			Name:   name,
			Mode:   modes[s.mode],
			Mods:   nil,
			Replay: nil,
		}
	}
	return nil
}

// err will be assigned to return value 'err'.
func (c Chart) Choose() (fsys fs.FS, name string, err error) {
	// const noVideo = 1
	// u := fmt.Sprintf("%s%d?n=%d", APIDownload, c.ParentSetId, noVideo)
	u := c.URLDownload()
	// fmt.Printf("download URL: %s\n", u)
	resp, err := http.Get(u)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fsys, err = zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return
	}
	return fsys, c.OsuFile, err
}
func (s Scene) Draw(screen draws.Image) {
	if s.loading {
		s.Loading.Draw(screen)
	}
	s.Background.Draw(screen)
	switch s.Focus {
	case FocusSearch, FocusChartSet:
		s.ChartSets.Draw(screen)
	case FocusChart:
		s.Charts.Draw(screen)
	}
	s.DebugPrint(screen)
}
func (s Scene) DebugPrint(screen draws.Image) {
	speed := *s.speedFactors[modes[s.mode]]
	ebitenutil.DebugPrint(screen.Image, fmt.Sprintf("FPS: %.2f\n"+
		"TPS: %.2f\n"+
		"Mode (F1): %s\n"+
		// "Sub mode (F2/F3): %s\n"+
		"\n"+
		"Music volume (Ctrl+ Left/Right): %.0f%%\n"+
		"Sound volume (Alt+ Left/Right): %.0f%%\n"+
		"Brightness (Ctrl+ O/P): %.0f%%\n"+
		"Offset (Shift+ Left/Right): %dms\n"+
		"\n"+
		"Speed (PageUp/Down): %.0f (Exposure time: %.0fms)\n"+
		"\n"+
		"Query: %s\n"+
		"No: %d (Page: %d)\n"+
		"Level limit to 10 (F4): %v\n",
		ebiten.ActualFPS(),
		ebiten.ActualTPS(),
		[]string{"Piano4", "Piano7", "Drum"}[s.mode],
		// fmt.Sprintf("%d Key", s.subMode),

		*s.volumeMusic*100,
		*s.volumeSound*100,
		*s.brightness*100,
		*s.offset,

		speed*100, s.exposureTimes[modes[s.mode]](speed),
		s.query,
		s.page*RowCount+s.ChartSets.cursor, s.page,
		s.levelLimit,
	))
}

type Return struct {
	FS     fs.FS
	Name   string
	Mode   int
	Mods   interface{}
	Replay *osr.Format
}

// It goes roughly triangular number.
func Level(sr float64) int { return int(math.Pow(sr, 1.7)) }
