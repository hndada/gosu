package gosu

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
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
	defaultBG    common.FixedSprite
}

func newSceneSelect(cwd string) *SceneSelect {
	s := new(SceneSelect)
	s.mods = mania.NewMods()
	s.defaultBG = common.DefaultBG()
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
}

// todo: args != nil 필요한가?
func (s *SceneSelect) Close(args *scene.Args) bool {
	if s.close && args.Next == "" {
		args.Next = "mania.Scene"
		args.Args = argsSelectToMania
	}
	return s.close
}

func (s *SceneSelect) reload() {
	s.ready = false
	ebiten.SetWindowTitle("gosu")
	cs := updateCharts(cwd)
	s.updatePanels(cs)
	s.ready = true
}

// 폴더에서 직접 삭제하면 새로 고침해야함
func (s *SceneSelect) updatePanels(cs []*mania.Chart) {
	for _, c := range cs {
		t := fmt.Sprintf("(%dKey Lv %.1f) %s [%s]", c.KeyCount, c.Level, c.MusicName, c.ChartName)
		p := ui.NewPanel(t)
		s.panelHandler.Append(p)
	}
}
