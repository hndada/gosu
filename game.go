package gosu

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/format/osr"
)

type Game struct {
	Scene
	Modes []Mode
	ModeType
	VolumeHandler ctrl.F64Handler
}
type Scene interface {
	Update() any
	Draw(screen *ebiten.Image)
}

func NewGame(modes []Mode) *Game {
	g := new(Game)
	// Todo: load settings here
	g.Modes = modes
	g.LoadChartInfosSet()    // 1. Load chart info and score data
	g.TidyChartInfosSet()    // 2. Check removed chart
	for i := range g.Modes { // 3. Check added chart
		// Each mode scans Music root independently.
		g.Modes[i].ChartInfos = LoadNewChartInfos(MusicRoot, &g.Modes[i])
	}
	g.SaveChartInfosSet() // 4. Save chart infos to local file
	LoadSounds("skin/sound")
	LoadGeneralSkin()
	for _, mode := range g.Modes {
		mode.LoadSkin()
	}
	g.VolumeHandler = NewVolumeHandler(&Volume)

	ebiten.SetWindowTitle("gosu")
	ebiten.SetWindowSize(WindowSizeX, WindowSizeY)
	ebiten.SetMaxTPS(MaxTPS)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	return g
}

func (g *Game) Update() error {
	g.VolumeHandler.Update()
	if g.Scene == nil {
		g.Scene = NewSceneSelect(g.Modes, &g.ModeType)
	}
	args := g.Scene.Update()
	switch args := args.(type) {
	case error:
		return args
	case PlayToResultArgs: // Todo: SceneResult
		g.Scene = NewSceneSelect(g.Modes, &g.ModeType)
	case SelectToPlayArgs:
		var err error
		g.Scene, err = g.Modes[args.ModeType].NewScenePlay(args.Path, args.Mods, args.Replay)
		if err != nil {
			return err
		}
	case nil:
		return nil
	}
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	g.Scene.Draw(screen)
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenSizeX, screenSizeY
}

type SelectToPlayArgs struct {
	ModeType
	Path   string
	Mods   Mods
	Replay *osr.Format
}

type PlayToResultArgs struct {
	Result
}
