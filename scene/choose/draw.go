package choose

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/scene"
)

func (s Scene) Draw(screen draws.Image) {
	// I don't want to put BaseXxx
	// s.drawBackground() // from baseScene
	s.drawList(screen)
	s.drawSearchBox(screen)
	// Todo: s.drawPanel()?
}

// May add extra effect to box arrangement. e.g., x -= y / 5
func (s Scene) drawList(screen draws.Image) {
	const (
		tx = 20
		ty = 20
	)
	for i := range l.Texts[:l.cursor] {
		if i > RowCount/2 {
			break
		}
		row := boxMask
		op := draws.Op{}
		op.GeoM.Translate(ScreenSizeX-1*RowWidth, ScreenSizeY/2-RowHeight/2-float64(i+1)*RowHeight)
		// row.Move(0, -float64(i+1)*RowHeight)
		row.Draw(dst, op)

		t := draws.NewText(l.Texts[l.cursor-i-1], draws.LoadDefaultFace(20))
		op.GeoM.Translate(tx, ty)
		t.Draw(dst, op)
	}
	for i := range l.Texts[l.cursor:] {
		if i > RowCount/2 {
			break
		}
		row := boxMask
		op := draws.Op{}
		op.GeoM.Translate(ScreenSizeX-1*RowWidth, ScreenSizeY/2-RowHeight/2+float64(i)*RowHeight)
		// row.Move(0, -float64(i+1)*RowHeight)
		row.Draw(dst, op)
		if i == 0 {
			op.GeoM.Translate(-RowShrink+10, 0)
		}
		t := draws.NewText(l.Texts[l.cursor+i], draws.LoadDefaultFace(20))
		op.GeoM.Translate(tx, ty)
		t.Draw(dst, op)
	}
}

func (s Scene) drawSearchBox(screen draws.Image) {
	t := *&s.queryTypeWriter.Text()
	if t == "" {
		t = "Type for search..."
	}
	const a = "searching..."
	count := 0
	if count == 0 {
		fmt.Sprintf("found no charts", count)
	} else {
		fmt.Sprintf("%s charts found", count)
	}
	text.Draw(screen, t, scene.Face16, int(d.X), int(d.Y)+25, color.White)
}

// imported from old code.
// It was for displaying text when loading.
func loadingText(t draws.Timer) string {
	s := "Loading"
	c := int(3*t.Age() + 1)
	s += strings.Repeat(".", c)
	return s
}
