// KeyDrawer draws KeyDownSprite at least for 30ms, KeyUpSprite otherwise.
// KeyDrawer uses MinCountdown instead of MaxCountdown.
type KeyDrawer struct {
	MinCountdown   int
	Countdowns     []int
	KeyUpSprites   []draws.Sprite
	KeyDownSprites []draws.Sprite
	lastPressed    []bool
	pressed        []bool
}

func (d *KeyDrawer) Update(lastPressed, pressed []bool) {
	d.lastPressed = lastPressed
	d.pressed = pressed
	for k, countdown := range d.Countdowns {
		if countdown > 0 {
			d.Countdowns[k]--
		}
	}
	for k, now := range d.pressed {
		last := d.lastPressed[k]
		if !last && now {
			d.Countdowns[k] = d.MinCountdown
		}
	}
}
func (d KeyDrawer) Draw(screen *ebiten.Image) {
	for k, p := range d.pressed {
		if p || d.Countdowns[k] > 0 {
			d.KeyDownSprites[k].Draw(screen, ebiten.DrawImageOptions{})
		} else {
			d.KeyUpSprites[k].Draw(screen, ebiten.DrawImageOptions{})
		}
	}
}
