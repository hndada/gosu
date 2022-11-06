func (d RollDrawer) Draw(screen *ebiten.Image) {
	const (
		head = iota
		tail
		body
		dot
	)
	max := len(d.Rolls) - 1
	for i := range d.Rolls {
		headNote := d.Rolls[max-i]
		if headNote.Position(d.Time) > maxPosition {
			continue
		}
		tailNote := *headNote
		tailNote.Time += headNote.Duration
		if tailNote.Position(d.Time) < minPosition {
			continue
		}
		op := ebiten.DrawImageOptions{}
		op.ColorM.ScaleWithColor(ColorYellow)
		for kind, sprite := range d.Sprites[headNote.Size][:3] {
			if kind == body {
				length := tailNote.Position(d.Time) - headNote.Position(d.Time)
				sprite.SetSize(length, sprite.H())
			}
			if kind == tail {
				sprite.Move(tailNote.Position(d.Time), 0)
			} else {
				sprite.Move(headNote.Position(d.Time), 0)
			}
			sprite.Draw(screen, op)
		}
		// bodySprite := d.Sprites[head.Size][body]
		// length := tail.Position(d.Time) - head.Position(d.Time)
		// bodySprite.SetSize(length, bodySprite.H())
		// bodySprite.Move(head.Position(d.Time), 0)
		// bodySprite.Draw(screen, op)

		// headSprite := d.HeadSprites[head.Size]
		// headSprite.Move(head.Position(d.Time), 0)
		// headSprite.Draw(screen, op)

		// tailSprite := d.TailSprites[tail.Size]
		// tailSprite.Move(tail.Position(d.Time), 0)
		// tailSprite.Draw(screen, op)
	}
}