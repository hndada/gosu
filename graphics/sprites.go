package graphics

import (
	"github.com/hndada/gosu/settings"
)

type GameSprites struct {
	Name            string
	BoxLeft         Sprite    // unscaled
	BoxMiddle       ExpSprite // unscaled
	BoxRight        Sprite    // unscaled
	ChartPanelFrame Sprite    // unscaled

	Score           [10]Sprite   // unscaled
	ManiaCombo      [10]Sprite   // unscaled
	ManiaHitResults [5]Sprite    // unscaled
	ManiaStages     []ManiaStage // 키별로 option 다름

	skin skin
}

func (s *GameSprites) SkinName() string { return s.skin.name }

// todo: combo, score 숫자 표시 어떻게 하지
func (s *GameSprites) Render(settings *settings.Settings) {
	// screenSize := settings.screenSize
	// settings.ComboPosition
	// settings.HitResultPosition
	for key := range maniaNoteKinds {
		s.ManiaStages[key].Render(settings, s.skin)
	}
}
