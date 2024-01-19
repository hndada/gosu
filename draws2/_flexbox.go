package draws

// Box is a set of multiple drawable structs.
// Node vs Box: Node feels like for logical ones, Box feels like for visual ones.
// Other candidates: Folder, Nest
// l.Source.Text = "Hello, World!" // It works!

// No Z index. Reorder children manually on each component.
type Flexbox struct {
	Children []Drawable
}

// Extend vs Expand
// Extend: Make something larger by adding to it.
// Expand: Make something larger by stretching it
type ExtendOptions struct {
	Spacing          Length
	Direction        int
	CollapseFirstBox bool
}

// Extend works as Flexible box.
// X, Y, Aligns, Parent will be newly set.
// Y: Height + Spacing
func (fb *Flexbox) Extend(ds []Drawable, opts ExtendOptions) {

}
