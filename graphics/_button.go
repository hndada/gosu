package graphic


type Button struct {
}

func init() {
	introImage, _ = ebiten.NewImage(screenWidth, screenHeigth, ebiten.FilterDefault)
	const (
		buttonWidth  = 250
		buttonHeigth = 100
	)
	bigButtons := []string{"gosu!"}
	for i, b := range bigButtons {
		width := buttonWidth * 2
		height := buttonHeigth * 3
		x := i*width + 300
		y := i*height + 300
		ebitenutil.DrawRect(introImage, float64(x), float64(y), float64(width), float64(height), color.White)
		text.Draw(introImage, b, mFont, x+width/4, y+height/2, color.Black)
	}
	smallButtons := []string{"(Multi)", "(Edit)", "Options"}
	for i, b := range smallButtons {
		width := buttonWidth
		height := buttonHeigth
		x := 300 + buttonWidth*2
		y := i*height + 300
		ebitenutil.DrawRect(introImage, float64(x), float64(y), float64(width), float64(height), color.White)
		text.Draw(introImage, b, mFont, x+width/4, y+height/2, color.Black)
	}
}
