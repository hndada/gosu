package choose

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"image/color"
	"io"
	"io/fs"
	"net/http"
	"runtime/debug"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
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

	choose audios.Sound
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
		&piano.S.SpeedScale, &piano.S.SpeedScale, &drum.S.SpeedScale}
	s.exposureTimes = []func(float64) float64{
		piano.ExposureTime, piano.ExposureTime, drum.ExposureTime}
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
	s.handleEnter()
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)
	debug.SetGCPercent(100)
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
	scene.VolumeMusic.Update()
	scene.VolumeSound.Update()
	scene.Brightness.Update()
	scene.Offset.Update()
	scene.SpeedScales[s.mode].Update()
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
		s.choose.Play(*s.volumeSound)
		err := s.handleEnter()
		if err != nil {
			fmt.Println(err)
		}
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
		if s.ChartSets.Update() {
			cset := s.ChartSets.Current()
			go func() {
				i, err := ebitenutil.NewImageFromURL(cset.URLCover("cover", Large))
				if err != nil {
					return
				}
				s.bgCh <- draws.Image{Image: i}
			}()
		}
	case FocusChart:
		s.Charts.Update()
	}
	return nil
}
func (s *Scene) handleEnter() any {
	switch s.Focus {
	case FocusSearch:
		s.LoadChartSetList()
	case FocusChartSet:
		fmt.Println("Load chart list")
		s.LoadChartList()
	case FocusChart:
		fmt.Println("Play chart")
		c := s.Charts.Current()
		if c == nil {
			return errors.New("no chart loaded")
		}
		fs, name, err := c.Choose()
		if err != nil {
			return err
		}
		s.choose.Play(*s.volumeSound)
		return Return{
			FS:     fs,
			Name:   name,
			Mode:   s.mode,
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
	{
		const (
			x = 1200
			y = 100
		)
		t := s.query
		if t == "" {
			t = "Type for search..."
		}
		text.Draw(screen.Image, t, scene.Face16, x, y, color.Black)
	}
	speed := *s.speedFactors[s.mode]
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
		"Speed (PageUp/Down): %.0f (Exposure time: %.0fms)\n",
		ebiten.ActualFPS(),
		ebiten.ActualTPS(),
		[]string{"Piano4", "Piano7", "Drum"}[s.mode],
		// fmt.Sprintf("%d Key", s.subMode),

		*s.volumeMusic*100,
		*s.volumeSound*100,
		*s.brightness*100,
		*s.offset,

		speed*100, s.exposureTimes[s.mode](speed)))
}

type Return struct {
	FS     fs.FS
	Name   string
	Mode   int
	Mods   interface{}
	Replay *osr.Format
}
