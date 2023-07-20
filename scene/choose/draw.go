package choose

import (
	"image/color"

	"github.com/hndada/gosu/draws"
)

// Background brightness at Song select: 60% (153 / 255), confirmed.
// Score box color: Gray128 with 50% transparent
// Hovered Score box color: Gray96 with 50% transparent
func (s Scene) Draw(dst draws.Image) {
	s.drawBackground()
	s.drawChartTree(dst)
	// Todo: s.drawSearchBox(screen)
	// Todo: s.drawPanel()
}

var (
	black = color.NRGBA{R: 16, G: 16, B: 16, A: 128}
	gray  = color.NRGBA{R: 128, G: 128, B: 128, A: 128}
)

func (s Scene) drawChartTree(dst draws.Image) {
	half := s.ListItemCount()/2 + 1

	// upper part
	var c int
	for n := s.currentNode.Prev(); n != nil; n = n.Prev() {
		c++

		box := s.BoxMaskSprite
		dy := float64(c) * s.ListItemHeight
		box.Move(0, -dy)

		op := draws.Op{}
		switch n.Type {
		case FolderNode:
			op.ColorM.ScaleWithColor(black)
		case ChartNote:
			op.ColorM.ScaleWithColor(gray)
		}
		box.Draw(dst, op)

		if c >= half {
			break
		}
	}

	// lower part
	c = 0 // reset
	for n := s.currentNode.Next(); n != nil; n = n.Next() {
		c++

		box := s.BoxMaskSprite
		dy := float64(c) * s.ListItemHeight
		box.Move(0, dy)

		op := draws.Op{}
		switch n.Type {
		case FolderNode:
			op.ColorM.ScaleWithColor(black)
		case ChartNote:
			op.ColorM.ScaleWithColor(gray)
		}
		box.Draw(dst, op)

		if c >= half {
			break
		}
	}
}
