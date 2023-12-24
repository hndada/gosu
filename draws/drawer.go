package draws

// Sprite, Animation, TextBox implement Drawer.
type Drawer interface {
	Draw(dst Image)
}
