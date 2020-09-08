package graphics

import (
	"github.com/hndada/gosu/settings"
)

type GameSprites struct {
	Name string
	// BoxLeft         Sprite // unscaled
	// BoxMiddle       Sprite // unscaled
	// BoxRight        Sprite // unscaled
	// ChartPanelFrame Sprite // unscaled

	Score           [10]Sprite         // unscaled
	ManiaCombo      [10]Sprite         // unscaled
	ManiaHitResults [5]Sprite          // unscaled
	ManiaStages     map[int]ManiaStage // 키별로 option 다름

	skin skin
}

func (s *GameSprites) SkinName() string { return s.skin.name }

// todo: combo, score 숫자 표시 어떻게 하지
func (s *GameSprites) Render(settings *settings.Settings) {
	// screenSize := settings.screenSize
	// settings.ComboPosition
	// settings.HitResultPosition
	s.skin.LoadSkin(`C:\Users\hndada\Documents\GitHub\hndada\gosu\test\Skin`)
	if s.ManiaStages == nil {
		s.ManiaStages = make(map[int]ManiaStage)
	}
	// for key := range maniaNoteKinds {
	for _, key := range []int{4, 7} {
		stage := ManiaStage{Keys: key}
		stage.Render(settings, s.skin)
		s.ManiaStages[key] = stage
	}
}
