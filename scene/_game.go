package game

import (
	"runtime/debug"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osr"
)

// 1. Load chart info and score data
// 2. Check removed chart
// 3. Check added chart
// Each mode scans Music root independently.
// 4. Save chart infos to local file
func NewGame(props []ModeProp) *Game {
	g := &Game{}

	LoadChartInfosSet(props)
	TidyChartInfosSet(props)
	for i, prop := range modeProps {
		modeProps[i].ChartInfos = prop.LoadNewChartInfos(MusicRoot)
	}
	SaveChartInfosSet(props)
	LoadGeneralSkin()
	for _, mode := range modeProps {
		mode.LoadSkin()
	}
	LoadHandlers(props)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
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
}
