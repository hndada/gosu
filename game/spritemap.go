package game

// TilesTemplate
type SpriteMapTemplate struct {
	// Name string
	// BoxLeft         Sprite // unscaled
	// BoxMiddle       Sprite // unscaled
	// BoxRight        Sprite // unscaled
	// ChartPanelFrame Sprite // unscaled
	Score  [10]Sprite // unscaled
	loaded bool
}

// func (s *Sprites) SkinName() string { return s.skin.name }
var SpriteMap SpriteMapTemplate

// todo: combo, score 숫자 표시 어떻게 하지
func LoadSpriteMap(skinPath string) {
	// screenSize := settings.screenSize
	// settings.ComboPosition
	// settings.HitResultPosition
	// sk := LoadSkin(`C:\Users\hndada\Documents\GitHub\hndada\gosu\test\Skin`)
}

func (sm *SpriteMapTemplate) Loaded() bool { return sm.loaded }
