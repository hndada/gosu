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
	musicVolume   *float64
	volumeSound   *float64
	brightness    *float64
	offset        *int64
	speedFactors  []*float64
	exposureTimes []func(float64) float64

	// choose audios.Sound
	// musicCh      chan []byte
	// Music      audios.MusicPlayer
	// bgCh       chan draws.Image
	Background mode.BackgroundDrawer

	mode int
	// Mode ctrl.KeyHandler
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

	lastFocus   int
	keySettings []string //[]input.Key
	// KeySettings KeySettingsDrawer
	// returnCh chan Return
	r *Return
}

const (
	FocusSearch = iota
	FocusChartSet
	FocusChart
	FocusKeySettings
)

var DefaultBackground = draws.NewImage(ScreenSizeX, 444)

func NewScene() *Scene {
	s := &Scene{}
	s.musicVolume = &mode.S.MusicVolume
	s.volumeSound = &mode.S.SoundVolume
	s.brightness = &mode.S.BackgroundBrightness
	s.offset = &mode.S.MusicOffset
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
	// 	Volume:    &mode.S.SoundVolume,
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
	// 	Volume:    &mode.S.SoundVolume,
	// }

	// s.lastChartSets = make([][]*ChartSet, 3)
	// for i := range s.lastChartSets {
	// 	css, err := search("", i, 0)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	s.lastChartSets[i] = css
	// }
	{
		sprite := draws.NewSprite(DefaultBackground)
		sprite.SetScaleToW(ScreenSizeX)
		sprite.Locate(ScreenSizeX/2, ScreenSizeY/2, draws.CenterMiddle)
		s.Background.Sprite = sprite
	}
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
	if s.r != nil {
		r := s.r
		s.r = nil
		return *r
	}
	// select {
	// case r := <-s.returnCh:
	// 	return r
	// default:
	// }
	if s.loading {
		s.Loading.Update()
		return nil
	}
	scene.MusicVolume.Update()
	scene.SoundVolume.Update()
	scene.Brightness.Update()
	scene.MusicOffset.Update()
	scene.SpeedScales[modes[s.mode]].Update()
	s.Preview.Update()
	// select {
	// case i := <-s.bgCh:
	// 	sprite := draws.NewSprite(i)
	// 	sprite.SetScaleToW(ScreenSizeX)
	// 	sprite.Locate(ScreenSizeX/2, ScreenSizeY/2, draws.CenterMiddle)
	// 	s.Background.Sprite = sprite
	// default:
	// }
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
	if inpututil.IsKeyJustPressed(input.KeyF5) {
		if s.Focus != FocusKeySettings {
			s.lastFocus = s.Focus
		}
		s.keySettings = make([]string, 0)
		s.Focus = FocusKeySettings
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
		case FocusKeySettings:
			s.Focus = s.lastFocus
		}
	}
	// if s.Mode.Update() { // || s.SubMode.Update()
	// 	go s.LoadChartSetList()
	// }
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
			// var req *http.Request
			// var resp *http.Response
			// var err error
			u := cset.URLPreview()
			resp, err := http.Get(u)
			if err != nil {
				return
			}
			// if runtime.GOARCH == "wasm" {
			// 	client := &http.Client{}
			// 	req, err = http.NewRequest("GET", u, nil)
			// 	if err != nil {
			// 		return
			// 	}
			// 	req.Header.Add("js.fetch:mode", "no-cors")
			// 	resp, err = client.Do(req)
			// 	if err != nil {
			// 		return
			// 	}
			// } else {
			// 	resp, err = http.Get(u)
			// 	if err != nil {
			// 		return
			// 	}
			// }
			if err != nil || resp.StatusCode == 404 {
				fmt.Println(err)
				// if runtime.GOARCH != "wasm" {
				// 	return
				// }
			}
			defer resp.Body.Close()
			if !s.Preview.IsEmpty() {
				s.Preview.Close()
			}
			s.Preview, err = NewPreviewPlayer(resp.Body)
			if err != nil {
				fmt.Println("preview:", err)
				// if runtime.GOARCH != "wasm" {
				// 	return
				// }
			}
		}()
		go func() {
			// i, err := draws.NewImageFromURL("https://upload.wikimedia.org/wikipedia/commons/1/1f/As08-16-2593.jpg")
			i, err := draws.NewImageFromURL(cset.URLCover("cover", Large))
			if err != nil {
				return
			}
			s.Background.Sprite.Source = i
			// s.bgCh <- draws.Image{Image: i}
		}()
	case FocusChart:
		s.Charts.Update()
	case FocusKeySettings:
		for k := input.Key(0); k < input.KeyReserved0; k++ {
			if inpututil.IsKeyJustPressed(k) {
				name := input.KeyToName(k)
				if name[0] == 'F' && name != "F" {
					continue
				}
				s.keySettings = append(s.keySettings, name)
			}
		}
		switch s.mode {
		case 0:
			if len(s.keySettings) >= 4 {
				s.keySettings = s.keySettings[:4]
				s.keySettings = mode.NormalizeKeys(s.keySettings)
				piano.S.KeySettings[4] = s.keySettings
				s.Focus = s.lastFocus
			}
		case 1:
			if len(s.keySettings) >= 7 {
				s.keySettings = s.keySettings[:7]
				s.keySettings = mode.NormalizeKeys(s.keySettings)
				piano.S.KeySettings[7] = s.keySettings
				s.Focus = s.lastFocus
			}
		case 2:
			if len(s.keySettings) >= 4 {
				s.keySettings = s.keySettings[:4]
				s.keySettings = mode.NormalizeKeys(s.keySettings)
				drum.S.KeySettings[4] = s.keySettings
				s.Focus = s.lastFocus
			}
		}
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
		go func() {
			s.loading = true
			scene.UserSkin.Enter.Play(*s.volumeSound)
			c := s.Charts.Current()
			if c == nil {
				fmt.Println(errors.New("no chart loaded"))
				return
			}
			fs, name, err := c.Choose()
			if err != nil {
				fmt.Println(err)
			}
			s.Preview.Close()
			s.r = &Return{
				FS:     fs,
				Name:   name,
				Mode:   modes[s.mode],
				Mods:   nil,
				Replay: nil,
			}
			s.loading = false
		}()
	}
	return nil
}

// err will be assigned to return value 'err'.
func (c Chart) Choose() (fsys fs.FS, name string, err error) {
	// const noVideo = 1
	// u := fmt.Sprintf("%s%d?n=%d", APIDownload, c.ParentSetId, noVideo)
	u := c.URLDownload()
	resp, err := http.Get(u)
	if err != nil {
		return
	}
	// fmt.Printf("download URL: %s\n", u)
	// var req *http.Request
	// var resp *http.Response
	// if runtime.GOARCH == "wasm" {
	// 	client := &http.Client{}
	// 	req, err = http.NewRequest("GET", u, nil)
	// 	if err != nil {
	// 		return
	// 	}
	// 	req.Header.Add("js.fetch:mode", "no-cors")
	// 	resp, err = client.Do(req)
	// 	if err != nil {
	// 		return
	// 	}
	// } else {
	// 	resp, err = http.Get(u)
	// 	if err != nil {
	// 		return
	// 	}
	// }
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
	s.Background.Draw(screen)
	// if s.Focus == FocusKeySettings {
	// 	s.Focus = s.lastFocus
	// }
	switch s.Focus {
	case FocusSearch, FocusChartSet:
		s.ChartSets.Draw(screen)
	case FocusChart:
		s.Charts.Draw(screen)
	}
	if s.loading {
		s.Loading.Draw(screen)
	}
	s.DebugPrint(screen)
}
func (s Scene) DebugPrint(screen draws.Image) {
	speed := *s.speedFactors[modes[s.mode]]
	keySettings := [][]string{piano.S.KeySettings[4], piano.S.KeySettings[7], drum.S.KeySettings[4]}[s.mode]
	if s.Focus == FocusKeySettings {
		keySettings = s.keySettings
	}
	ebitenutil.DebugPrint(screen.Image, fmt.Sprintf("FPS: %.2f\n"+
		"TPS: %.2f\n"+
		"Mode (F1): %s\n"+
		// "Sub mode (F2/F3): %s\n"+
		"\n"+
		"Music volume (Ctrl+ Left/Right): %.0f%%\n"+
		"Sound volume (Alt+ Left/Right): %.0f%%\n"+
		"Brightness (Ctrl+ O/P): %.0f%%\n"+
		"MusicOffset (Shift+ Left/Right): %dms\n"+
		"\n"+
		"Speed (PageUp/Down): %.0f (Exposure time: %.0fms)\n"+
		"\n"+
		"Query: %s\n"+
		"No: %d (Page: %d)\n"+
		"Level limit to 10 (F4): %v\n"+
		"\n"+
		"Key settings: %v (F5) listening: %v)\n",
		ebiten.ActualFPS(),
		ebiten.ActualTPS(),
		[]string{"Piano4", "Piano7", "Drum"}[s.mode],
		// fmt.Sprintf("%d Key", s.subMode),

		*s.musicVolume*100,
		*s.volumeSound*100,
		*s.brightness*100,
		*s.offset,

		speed*100, s.exposureTimes[modes[s.mode]](speed),
		s.query,
		s.page*RowCount+s.ChartSets.cursor, s.page,
		s.levelLimit,

		keySettings, s.Focus == FocusKeySettings,
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
