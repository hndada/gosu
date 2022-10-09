package draws

import "github.com/hajimehoshi/ebiten/v2"

// Todo: implement nested Grid?
type Grid struct {
	Outer  Box
	Inners [][]Box
	Gap    Point
	// Widths  []float64
	// Heights []float64
}

func (g *Grid) SetSizes(ws, hs []float64) {
	for i, w := range ws {
		for j, h := range hs {
			g.Inners[i][j].Outer.SetSize(Point{w, h})
		}
	}
	// Todo: set Inners's XY
}

// Todo: should Grid and Box's Draw be given external translate value?
func (g Grid) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions) {
	g.Outer.Draw(screen, op)
	var offset Point
	for _, row := range g.Inners {
		offset.X = 0
		offset.Y += g.Gap.Y
		for _, box := range row {
			offset.X += g.Gap.X
			box.Point.Add(offset)
			box.Draw(screen, op, p)
			offset.X += box.Outer.Size().X
		}
		offset.Y += row[0].Outer.Size().Y
	}
}
