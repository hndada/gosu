package gosu

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/format/osr"
)

var (
	MusicVolumeHandler  ctrl.F64Handler
	EffectVolumeHandler ctrl.F64Handler
	VsyncSwitchHandler  ctrl.BoolHandler
)

type Game struct {
	Scene
	ModeProps []ModeProp
	Mode      int
}
type Scene interface {
	Update() any
	Draw(screen *ebiten.Image)
}

func init() {
	if runtime.GOOS == "windows" {
		fmt.Println("OpenGL mode has enabled.")
		os.Setenv("EBITEN_GRAPHICS_LIBRARY", "opengl")
	}
}
func NewGame(props []ModeProp) *Game {
	g := new(Game)
	// Todo: load settings here
	g.ModeProps = props
	g.LoadChartInfosSet()        // 1. Load chart info and score data
	g.TidyChartInfosSet()        // 2. Check removed chart
	for i := range g.ModeProps { // 3. Check added chart
		// Each mode scans Music root independently.
		g.ModeProps[i].ChartInfos = LoadNewChartInfos(MusicRoot, &g.ModeProps[i])
	}
	g.SaveChartInfosSet() // 4. Save chart infos to local file
	LoadSounds("skin/sound")
	LoadGeneralSkin()
	for _, mode := range g.ModeProps {
		mode.LoadSkin()
	}
	g.Mode = ModePiano4
	MusicVolumeHandler = NewVolumeHandler(
		&MusicVolume, []ebiten.Key{ebiten.Key2, ebiten.Key1})
	EffectVolumeHandler = NewVolumeHandler(
		&EffectVolume, []ebiten.Key{ebiten.Key4, ebiten.Key3})
	VsyncSwitchHandler = NewVsyncSwitchHandler(&VsyncSwitch)
	ebiten.SetWindowTitle("gosu")
	ebiten.SetWindowSize(WindowSizeX, WindowSizeY)
	ebiten.SetTPS(TPS)
	// ebiten.SetCursorMode(ebiten.CursorModeHidden)
	return g
}

func (g *Game) Update() (err error) {
	// g.MusicVolumeHandler.Update()
	// g.EffectVolumeHandler.Update()
	if g.Scene == nil {
		g.Scene = NewSceneSelect(g.ModeProps, &g.Mode)
	}
	args := g.Scene.Update()
	switch args := args.(type) {
	case error:
		return args
	case PlayToResultArgs: // Todo: SceneResult
		ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)
		VsyncSwitch = true
		debug.SetGCPercent(0)

		g.Scene = NewSceneSelect(g.ModeProps, &g.Mode)
	case SelectToPlayArgs:
		ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
		VsyncSwitch = false
		debug.SetGCPercent(100)

		g.Scene, err = g.ModeProps[args.Mode].NewScenePlay(
			args.Path, args.Replay, args.SpeedHandler)
		if err != nil {
			return
		}
	case nil:
		return
	}
	return
}
func (g *Game) Draw(screen *ebiten.Image) {
	g.Scene.Draw(screen)
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenSizeX, screenSizeY
}

type SelectToPlayArgs struct {
	Mode int
	Path string
	// Mods   Mods
	Replay       *osr.Format
	SpeedHandler ctrl.F64Handler
}

type PlayToResultArgs struct {
	Result
}
