type BoxOptions struct {
	BackgroundColor Color
}

func (b Box[T]) Draw(dst Image) {
	if !b.Collapsed || !b.Exposed(dst) {
		return
	}

	// Draw children in order.
	idxs := make([]int, len(b.Children))
	for i := range idxs {
		idxs[i] = i
	}
	sort.Slice(idxs, func(i, j int) bool {
		b1 := b.Children[idxs[i]].(Box[T])
		b2 := b.Children[idxs[j]].(Box[T])
		return b1.ZIndex < b2.ZIndex
	})

	after := false
	for _, idx := range idxs {
		if !after && b.Children[idx].(Box[T]).ZIndex >= 0 {
			after = true
			b.draw(dst)
		}
		b.Children[idx].Draw(dst)
	}
}
