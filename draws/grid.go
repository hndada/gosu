package draws

import "github.com/hajimehoshi/ebiten/v2"

// Todo: implement nested Grid?
// type Grid struct {
// 	Outer  Box
// 	Inners [][]Box
// 	Gap    Point
// 	Widths  []float64
// 	Heights []float64
// }

// Grid implements inteface Subject.
type Grid [][]Box

func NewGrid(inners [][]Box, ws, hs []float64, gap Point) *Grid {
	var g Grid
	var offset Point
	for i, h := range hs {
		offset.X = 0
		for j, w := range ws {
			g[i][j].Outer.SetSize(Point{w, h})
			g[i][j].Point = offset
			g[i][j].Origin2 = AtMin
			offset.X += w + gap.X
		}
		offset.Y += h + gap.Y
	}
	return &g
}
func (g Grid) Size() (p Point) {
	row := g[len(g)-1]
	box := row[len(row)-1]
	return box.OuterMax()
}

// Todo: implement
func (g *Grid) SetSize(size Point) {}
func (g Grid) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions, p Point) {
	for _, row := range g {
		for _, box := range row {
			box.Draw(screen, op, p)
		}
	}
}

// // Todo: should Grid and Box's Draw be given external translate value?
// func (g Grid) Draw1(screen *ebiten.Image, op ebiten.DrawImageOptions) {
// 	g.Outer.Draw(screen, op)
// 	var offset Point
// 	for _, row := range g.Inners {
// 		offset.X = 0
// 		offset.Y += g.Gap.Y
// 		for _, box := range row {
// 			offset.X += g.Gap.X
// 			box.Point.Add(offset)
// 			box.Draw(screen, op, p)
// 			offset.X += box.Outer.Size().X
// 		}
// 		offset.Y += row[0].Outer.Size().Y
// 	}
// }
