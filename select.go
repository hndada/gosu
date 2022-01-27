package gosu

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu/common"
	"github.com/hndada/gosu/engine/scene"
	"github.com/hndada/gosu/engine/ui"
	"github.com/hndada/gosu/mania"
)

var sceneSelect *SceneSelect

type SceneSelect struct {
	ready        bool // TODO: use channel?
	close        bool
	panelHandler ui.PanelHandler
	mods         mania.Mods
	defaultBG    ui.FixedSprite
	boxSkin      ui.BoxSkin
}

// TODO: bg, music preview
func newSceneSelect(cwd string) *SceneSelect {
	s := new(SceneSelect)
	{
		dir := filepath.Join(cwd, "skin")
		name := "soft-slidertick.wav"
		sePath := filepath.Join(dir, name)
		s.panelHandler = ui.NewPanelHandler(common.Settings.ScreenSize, sePath)
	}
	s.mods = mania.NewMods()
	s.defaultBG = common.DefaultBG()
	s.boxSkin = ui.BoxSkin{
		Left:   common.Skin.BoxLeft,
		Middle: common.Skin.BoxMiddle,
		Right:  common.Skin.BoxRight,
	}
	s.reload()
	return s
}

func (s *SceneSelect) Ready() bool { return s.ready }

// anonymous struct: grouped globals
// reflect: fields should be exported
var argsSelectToMania struct {
	Chart *mania.Chart
	Mods  mania.Mods
}

func (s *SceneSelect) Update() error {
	idx := s.panelHandler.Update()
	if idx != -1 {
		argsSelectToMania.Chart = charts[idx]
		argsSelectToMania.Mods = s.mods
		s.close = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyO) {
		mania.Settings.GeneralSpeed -= 0.005
		if mania.Settings.GeneralSpeed < 0.01 {
			mania.Settings.GeneralSpeed = 0.01
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyP) {
		mania.Settings.GeneralSpeed += 0.005
		if mania.Settings.GeneralSpeed > 0.4 {
			mania.Settings.GeneralSpeed = 0.4
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyDigit1) {
		common.Settings.ScoreMode = common.ScoreModeNaive
		cwd, _ := os.Getwd()
		loadCharts(cwd)
	} else if ebiten.IsKeyPressed(ebiten.KeyDigit2) {
		common.Settings.ScoreMode = common.ScoreModeWeighted
		cwd, _ := os.Getwd()
		loadCharts(cwd)
	} else if ebiten.IsKeyPressed(ebiten.KeyDigit3) {
		common.Settings.ScoreMode = common.ScoreModeOsuLegacy
		cwd, _ := os.Getwd()
		loadCharts(cwd)
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		common.Settings.IsAuto = !common.Settings.IsAuto
	}
	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		common.Settings.AutoInstability -= 1
		if common.Settings.AutoInstability < 0 {
			common.Settings.AutoInstability = 0
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyX) {
		common.Settings.AutoInstability += 1
		if common.Settings.AutoInstability > 100 {
			common.Settings.AutoInstability = 100
		}
	}
	return nil
}

func (s *SceneSelect) Draw(screen *ebiten.Image) {
	// for i, cp := range s.chartPanels {
	// 	if i == s.cursor {
	// 		screen.DrawImage(cp.BG, cp.OpBG)
	// 		break
	// 	}
	// }
	s.defaultBG.Draw(screen)
	s.panelHandler.Draw(screen)

	var scoreMode string
	switch common.Settings.ScoreMode {
	case common.ScoreModeNaive:
		scoreMode = "Naive"
	case common.ScoreModeWeighted:
		scoreMode = "Weighted"
	case common.ScoreModeOsuLegacy:
		scoreMode = "osu! legacy"
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		`Speed(Press O/P): %.1f
Score mode(Press 1/2/3): %s
Auto mode(Press A): %t
Auto instability(Press Z/X): %.0f
`, mania.Settings.GeneralSpeed*100, scoreMode, common.Settings.IsAuto, common.Settings.AutoInstability))
}

// TODO: Dose it need args != nil ?
func (s *SceneSelect) Close(args *scene.Args) bool {
	if s.close && args.Next == "" {
		args.Next = "mania.Scene"
		args.Args = argsSelectToMania
	}
	return s.close
}

// Need to refresh manually when one has deleted directly on file explorer
func (s *SceneSelect) reload() {
	s.ready = false
	ebiten.SetWindowTitle("gosu")
	cs := updateCharts(cwd)
	for _, c := range cs {
		t := fmt.Sprintf("(%dKey Lv %.1f) %s [%s]", c.KeyCount, c.Level, c.MusicName, c.ChartName)
		p := ui.NewPanel(t, s.boxSkin)
		s.panelHandler.Append(p)
	}
	s.ready = true
}
