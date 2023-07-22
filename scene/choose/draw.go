package choose

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu/draws"
)

// Background brightness at Song select: 60% (153 / 255), confirmed.
// Score box color: Gray128 with 50% transparent
// Hovered Score box color: Gray96 with 50% transparent
func (s Scene) Draw(screen draws.Image) {
	s.drawBackground(screen)
	s.drawChartTree(screen)
	// Todo: s.drawSearchBox(screen)
	// Todo: s.drawPanel(screen)
	if s.DebugPrint {
		f := fmt.Fprintf
		b := strings.Builder{}
		f(&b, s.Config.DebugString())
		ebitenutil.DebugPrint(screen.Image, b.String())
	}
}

func (s Scene) drawChartTree(dst draws.Image) {
	half := s.ListItemCount()/2 + 1
	var n *Node

	// upper part
	n = s.chartTreeNode.Prev()
	for i := 0; i < half; i++ {
		if n == nil {
			break
		}
		s.drawChartTreeNode(dst, s.chartTreeNode, i)
		n = n.Prev()
	}

	// lower part
	n = s.chartTreeNode.Next()
	for i := 0; i < half; i++ {
		if n == nil {
			break
		}
		s.drawChartTreeNode(dst, s.chartTreeNode, i)
		n = n.Next()
	}

	// middle part
	s.drawChartTreeNode(dst, s.chartTreeNode, 0)
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
	switch offset {
	case 0:
		// Todo: glow effect
	default:
		dx := s.ListItemShrink
		dy := float64(offset) * s.ListItemHeight
		box.Move(dx, dy)
	}
	op := draws.Op{}
	op.ColorM.ScaleWithColor(clr)
	box.Draw(dst, op)
}
