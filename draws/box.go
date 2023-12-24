package draws

type Box struct {
	Size     Vector2
	Position Vector2
	Anchor   Anchor
	// Filter   ebiten.Filter
	// Color    colorm.ColorM
}

type Sprite2 struct {
	Box
	Source Source
}
type Animation2 struct {
	Box
	Sources []Source
}
