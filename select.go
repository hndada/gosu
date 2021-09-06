package gosu

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu/common"
	"github.com/hndada/gosu/engine/scene"
	"github.com/hndada/gosu/engine/ui"
	"github.com/hndada/gosu/mania"
)

var sceneSelect *SceneSelect

type SceneSelect struct {
	ready        bool // todo: 채널 이용하여 wait 시도?
	close        bool
	panelHandler ui.PanelHandler
	mods         mania.Mods
	defaultBG    ui.FixedSprite
	boxSkin      ui.BoxSkin
}

func newSceneSelect(cwd string) *SceneSelect {
	s := new(SceneSelect)
	s.panelHandler = ui.NewPanelHandler(common.Settings.ScreenSize)
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
		mania.Settings.GeneralSpeed += 0.005
		if mania.Settings.GeneralSpeed > 0.4 {
			mania.Settings.GeneralSpeed = 0.4
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyP) {
		mania.Settings.GeneralSpeed -= 0.005
		if mania.Settings.GeneralSpeed < 0.01 {
			mania.Settings.GeneralSpeed = 0.01
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
	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		`Speed: %.1f`, mania.Settings.GeneralSpeed*100))
}

// todo: args != nil 필요한가?
func (s *SceneSelect) Close(args *scene.Args) bool {
	if s.close && args.Next == "" {
		args.Next = "mania.Scene"
		args.Args = argsSelectToMania
	}
	return s.close
}

// 폴더에서 직접 삭제하면 새로 고침해야함
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
