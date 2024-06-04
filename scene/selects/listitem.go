package selects

import (
	"fmt"
	"image/color"

	"github.com/hndada/gosu/game"
)

// Music name itself may be duplicated.
// Artist + Title (Music name) may be unique.
func FolderText(c *game.ChartHeader) string {
	return fmt.Sprintf("%s - %s", c.MusicName, c.Artist)
}

// Todo: add level database. attach level info to the text.
// Then, sort by the level.

// Memo: make([]T, len) and make([]T, 0, len) is prone to be erroneous.
func ItemText(c *game.ChartHeader) string {
	// return fmt.Sprintf("[Lv. %.0f] %s [%s]", c.Level, c.MusicName, c.ChartName) // [Lv. %4.2f]
	return fmt.Sprintf("%s [%s]", c.MusicName, c.ChartName)
}

// Background brightness at Song select: 60% (153 / 255), confirmed.
// Score box color: Gray128 with 50% transparent
// Hovered Score box color: Gray96 with 50% transparent
func (s Scene) Draw(screen draws.Image) {
	s.drawBackground(screen)
	s.drawChartTree(screen)
}

func (s Scene) drawChartTree(dst draws.Image) {
	half := s.ChartTreeNodeCount()/2 + 1
	var n *Node

	// upper part
	n = s.chartTreeNode.Prev()
	for i := 1; i <= half; i++ {
		if n == nil {
			break
		}
		s.drawChartTreeNode(dst, n, -i)
		n = n.Prev()
	}

	// middle part
	n = s.chartTreeNode
	s.drawChartTreeNode(dst, n, 0)

	// lower part
	n = s.chartTreeNode.Next()
	for i := 1; i <= half; i++ {
		if n == nil {
			break
		}
		s.drawChartTreeNode(dst, n, i)
		n = n.Next()
	}
}

func (s Scene) drawChartTreeNode(dst draws.Image, n *Node, offset int) {
	// var (
	// 	black = color.NRGBA{R: 16, G: 16, B: 16, A: 128}
	// 	gray  = color.NRGBA{R: 128, G: 128, B: 128, A: 128}
	// )
	var clr color.NRGBA
	switch n.Type {
	case FolderNode:
		clr = color.NRGBA{R: 64, G: 255, B: 255, A: 128} // skyblue
	case ChartNode:
		clr = color.NRGBA{R: 255, G: 128, B: 255, A: 128} // pink
	}

	box := s.BoxMaskSprite
	text := s.newNodeTextSprite(n)
	op := draws.Op{}
	op.ColorM.ScaleWithColor(clr)
	switch offset {
	case 0:
		// Todo: glow effect
	default:
		dx := s.ChartTreeNodeShrink
		if offset > 0 {
			dx += 10 * float64(offset)
		} else {
			dx -= 10 * float64(offset)
		}
		dy := float64(offset) * s.ChartTreeNodeHeight
		box.Move(dx, dy)
		text.Move(dx, dy)
	}
	text.Draw(dst, op)
	box.Draw(dst, op)
}

// Todo: handle dx, dy automatically with face size.
func (s Scene) newNodeTextSprite(n *Node) draws.Sprite {
	const (
		dx = 9
		dy = 18
	)
	src := draws.NewText(n.Data, scene.Face(24))
	text := draws.NewSprite(src)
	text.Locate(s.ScreenSize.X-s.ChartTreeNodeWidth+dx, s.ScreenSize.Y/2+dy, draws.LeftMiddle)
	return text
}
