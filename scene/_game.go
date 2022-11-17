package game

import (
	"runtime/debug"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osr"
)

var (
	modeProps     []ModeProp
	sceneSelect   *SceneSelect
	tailExtraTime *float64 // For cache.
)

// Todo: load settings
func NewGame(props []ModeProp) *Game {
	modeProps = props
	tailExtraTime = modeProps[ModePiano4].Settings["TailExtraTime"]
	g := &Game{}
	SetKeySettings(props)
	// 1. Load chart info and score data
	// 2. Check removed chart
	// 3. Check added chart
	// Each mode scans Music root independently.
	LoadChartInfosSet(props)
	TidyChartInfosSet(props)
	for i, prop := range modeProps {
		modeProps[i].ChartInfos = prop.LoadNewChartInfos(MusicRoot)
	}
	SaveChartInfosSet(props) // 4. Save chart infos to local file
	LoadGeneralSkin()
	for _, mode := range modeProps {
		mode.LoadSkin()
	}
	LoadHandlers(props)
	ebiten.SetWindowTitle("gosu")
	ebiten.SetWindowSize(WindowSizeX, WindowSizeY)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(TPS)
	modeHandler.Max = len(props)
	sceneSelect = NewSceneSelect()
	// ebiten.SetCursorMode(ebiten.CursorModeHidden)
	return g
}

func (g *Game) Update() (err error) {
	VolumeMusicKeyHandler.Update()
	VolumeSoundKeyHandler.Update()
	SpeedScaleKeyHandler.Update()
	OffsetKeyHandler.Update()
	TailExtraTimeKeyHandler.Update()
	if g.Scene == nil {
		g.Scene = sceneSelect
	}
	args := g.Scene.Update()
	switch args := args.(type) {
	case error:
		return args
	case PlayToResultArgs: // Todo: SceneResult
		// VolumeSound = 0.25 // Todo: resolve delayed effect sound playing
		ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)
		debug.SetGCPercent(100)
		g.Scene = sceneSelect
		ebiten.SetWindowTitle("gosu")
	case SelectToPlayArgs:
		// VolumeSound = 0 // Todo: resolve delayed effect sound playing
		ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
		debug.SetGCPercent(0)
		prop := modeProps[currentMode]
		g.Scene, err = prop.NewScenePlay(args.Path, args.Replay)
		if err != nil {
			return
		}
	}
	return
}
