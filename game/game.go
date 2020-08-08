package game

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

// Game Programming Draw
type Game struct {
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
