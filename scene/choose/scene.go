package choose

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"runtime/debug"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hndada/gosu/audios"
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

	choose     audios.Sound
	Music      audios.MusicPlayer
	Background mode.BackgroundDrawer

	mode    int
	Mode    ctrl.KeyHandler
	subMode int
	SubMode ctrl.KeyHandler
	query   string
	Query   TypeWriter
	page    int

	Focus     int
	ChartSets ChartSetList
	Charts    ChartList
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
	s.Mode = ctrl.KeyHandler{
		Handler: ctrl.IntHandler{
			Value: &s.mode,
			Min:   0,
			Max:   2 - 1, // There are two modes.
			Loop:  true,
		},
		Modifiers: []input.Key{},
		Keys:      [2]input.Key{-1, input.KeyF1},
		Sounds:    [2]audios.Sounder{scene.UserSkin.Swipe, scene.UserSkin.Swipe},
		Volume:    &mode.S.VolumeSound,
	}
	s.SubMode = ctrl.KeyHandler{
		Handler: ctrl.IntHandler{
			Value: &s.subMode,
			Min:   4,
			Max:   9,
			Loop:  true,
		},
		Modifiers: []input.Key{},
		Keys:      [2]input.Key{input.KeyF2, input.KeyF3},
		Sounds:    [2]audios.Sounder{scene.UserSkin.Swipe, scene.UserSkin.Swipe},
		Volume:    &mode.S.VolumeSound,
	}
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)
	debug.SetGCPercent(100)
	ebiten.SetWindowTitle("gosu")
	return &Scene{}
}
func isEnter() bool {
	return inpututil.IsKeyJustPressed(input.KeyEnter) ||
		inpututil.IsKeyJustPressed(input.KeyNumpadEnter)
}
func isBack() bool {
	return inpututil.IsKeyJustPressed(input.KeyEscape)
}

// Todo: resolve nested blocks
func (s *Scene) Update() any {
	scene.VolumeMusic.Update()
	scene.VolumeSound.Update()
	scene.Brightness.Update()
	scene.Offset.Update()
	scene.SpeedScales[s.mode].Update()
	if isEnter() {
		switch s.Focus {
		case FocusSearch:
			css, err := s.Search()
			if err != nil {
				fmt.Println(err)
			}
			s.ChartSets = NewChartSetList(css)
		case FocusChartSet:
			if len(s.ChartSets.ChartSets) > 0 {
				s.Charts = s.ChartSets.NewChartList()
				s.Focus = FocusChart
			}
		case FocusChart:
			c := s.Charts.Current()
			if c != nil {
				fs, name, err := c.Choose()
				if err != nil {
					fmt.Println(err)
				} else {
					s.choose.Play(*s.volumeSound)
					return Return{
						FS:     fs,
						Name:   name,
						Mode:   s.mode,
						Mods:   nil,
						Replay: nil,
					}
				}
			}
		}
	}
	if isBack() {
		switch s.Focus {
		case FocusSearch:
			s.Query.Reset()
		case FocusChartSet:
			s.Query.Reset()
			s.Focus = FocusSearch
		case FocusChart:
			s.Focus = FocusChartSet
		}
	}
	if s.Mode.Update() || s.SubMode.Update() {
		css, err := s.Search()
		if err != nil {
			fmt.Println(err)
		}
		s.ChartSets = NewChartSetList(css)
		// Background
	}
	switch s.Focus {
	case FocusChartSet:
		s.ChartSets.Update()
	case FocusChart:
		s.Charts.Update()
	}
	return nil
}

// err will be assigned to return value 'err'.
func (c Chart) Choose() (fsys fs.FS, name string, err error) {
	// const noVideo = 1
	// u := fmt.Sprintf("%s%d?n=%d", APIDownload, c.ParentSetId, noVideo)
	u := c.URLDownload()
	fmt.Printf("download URL: %s\n", u)
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
	s.Background.Draw(screen)
	switch s.Focus {
	case FocusChartSet:
		s.ChartSets.Draw(screen)
	case FocusChart:
		s.Charts.Draw(screen)
	}
	s.DebugPrint(screen)
}
func (s Scene) DebugPrint(screen draws.Image) {
	speed := *s.speedFactors[s.mode]
	ebitenutil.DebugPrint(screen.Image,
		fmt.Sprintf(
			"Mode (F1): %s\n"+
				"Sub mode (F2/F3): %s\n"+
				"\n"+
				"Music volume (Alt+ Left/Right): %.0f%%\n"+
				"Sound volume (Ctrl+ Left/Right): %.0f%%\n"+
				"Offset (Shift+ Left/Right): %dms\n"+
				"Brightness (Ctrl+ O/P): %.0f%%\n"+
				"\n"+
				"Speed (PageUp/Down): %.0f (Exposure time: %.0fms)\n"+
				[]string{"Piano", "Drum"}[s.mode],
			fmt.Sprintf("%d Key", s.subMode),

			*s.volumeMusic*100,
			*s.volumeSound*100,
			*s.offset,
			*s.brightness*100,

			speed*100, s.exposureTimes[s.mode](speed)))
}

type Return struct {
	FS     fs.FS
	Name   string
	Mode   int
	Mods   interface{}
	Replay *osr.Format
}
