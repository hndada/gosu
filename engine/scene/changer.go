package scene

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/game"
)

// todo: tick-based? time-based?
var initCountDown = ebiten.MaxTPS() * 4 / 5

type Changer struct {
	scene      Scene
	nextScene  Scene
	transScene *ebiten.Image
	countdown  int
	phase      phase
}

type phase int

const (
	phaseIdle phase = iota
	phaseFadeOut
	phaseFadeIn
)

var blackScreen *ebiten.Image

func init() {
	blackScreen, _ = ebiten.NewImage(game.Settings.ScreenSize.X,
		game.Settings.ScreenSize.Y, ebiten.FilterDefault)
	blackScreen.Fill(color.Black)
}

func NewChanger() *Changer {
	c := &Changer{}
	c.transScene, _ = ebiten.NewImage(game.Settings.ScreenSize.X,
		game.Settings.ScreenSize.Y, ebiten.FilterDefault)

	return c
}

func (c *Changer) Done() bool { return c.phase == phaseIdle }

func (c *Changer) Change(s1, s2 Scene) {
	c.scene = s1
	c.nextScene = s2
	c.countdown = initCountDown
	c.phase = phaseFadeOut
}

// countdown 동안 scene이 fade되면서 변화
func (c *Changer) Update() error {
	switch c.phase {
	case phaseIdle:
	case phaseFadeOut:
		if c.countdown == 0 && c.nextScene.Ready() {
			c.phase = phaseFadeIn
			c.countdown = initCountDown
		}
		if c.countdown > 0 {
			c.countdown--
		}
	case phaseFadeIn:
		if c.countdown == 0 {
			c.phase = phaseIdle
		}
		if c.countdown > 0 {
			c.countdown--
		}
	}
	return nil
}

func (c *Changer) Draw(screen *ebiten.Image) {
	switch c.phase {
	case phaseIdle:
	case phaseFadeOut:
		screen.DrawImage(blackScreen, &ebiten.DrawImageOptions{})

		value := float64(c.countdown) / float64(initCountDown)
		c.transScene.Clear()
		c.scene.Draw(c.transScene)

		op := &ebiten.DrawImageOptions{}
		op.ColorM.ChangeHSV(0, 1, value)
		screen.DrawImage(c.transScene, op)
	case phaseFadeIn:
		screen.DrawImage(blackScreen, &ebiten.DrawImageOptions{})

		value := 1 - float64(c.countdown)/float64(initCountDown)
		c.transScene.Clear()
		c.nextScene.Draw(c.transScene)

		op := &ebiten.DrawImageOptions{}
		op.ColorM.ChangeHSV(0, 1, value)
		screen.DrawImage(c.transScene, op)
	}
}
