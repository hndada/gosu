package mode

type CommonSprites struct {
	// Name string
	// BoxLeft         Sprite // unscaled
	// BoxMiddle       Sprite // unscaled
	// BoxRight        Sprite // unscaled
	// ChartPanelFrame Sprite // unscaled
	Score [10]Sprite // unscaled
}

// func (s *Sprites) SkinName() string { return s.skin.name }

// todo: combo, score 숫자 표시 어떻게 하지
// func (g *Game) RenderSprites() {}
func (s *CommonSprites) Render(settings *CommonSettings) {
	// screenSize := settings.screenSize
	// settings.ComboPosition
	// settings.HitResultPosition
	// sk := LoadSkin(`C:\Users\hndada\Documents\GitHub\hndada\gosu\test\Skin`)
}
