package graphics

import (
	"github.com/hajimehoshi/ebiten"
	"image"
	"image/color"
)

const (
	checkboxBorder  = 2
	checkboxRadius  = 8
	checkboxPadding = 8
)

var (
	imgUncheckedBox *ebiten.Image
	imgCheckedBox   *ebiten.Image
)

func init() {
	bw := 2 * checkboxRadius
	bh := bw
	fw := bw - 2*checkboxBorder
	fh := fw
	w := bw + 2*checkboxPadding
	h := bh + checkboxPadding

	fop := &ebiten.DrawImageOptions{}
	fop.GeoM.Translate(checkboxBorder, checkboxBorder)
	bop := &ebiten.DrawImageOptions{}
	bop.GeoM.Translate(checkboxPadding, checkboxPadding/2)

	clrBorder := color.RGBA{102, 0, 95, 255}
	clrUncheck := color.Transparent
	clrCheck := color.RGBA{95, 240, 252, 255}
	{
		fill, _ := ebiten.NewImage(fw, fh, ebiten.FilterDefault)
		fill.Fill(clrUncheck)
		box, _ := ebiten.NewImage(bw, bh, ebiten.FilterDefault)
		box.Fill(clrBorder)
		box.DrawImage(fill, fop)
		imgUncheckedBox, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
		imgUncheckedBox.DrawImage(box, bop)
	}
	{
		fill, _ := ebiten.NewImage(fw, fh, ebiten.FilterDefault)
		fill.Fill(clrCheck)
		box, _ := ebiten.NewImage(bw, bh, ebiten.FilterDefault)
		box.Fill(clrBorder)
		box.DrawImage(fill, fop)
		imgCheckedBox, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
		imgCheckedBox.DrawImage(box, bop)
	}
}

// todo: 글자 크기가 안 맞는 체크박스가 나올 수 있음, 현재 수동으로 해야함
// -> 폰트 및 글자 크기. dpi고정. text받고 렌더하고 체크박스까지 그리자
type Checkbox struct {
	MinPt     image.Point
	Text      *ebiten.Image
	checked   bool
	mouseDown bool
	onChanged func(c *Checkbox)
}

// 화면 등이 움직이면 point값 바꿔주면 됨
func NewCheckbox(s string, p image.Point) *Checkbox {
	c := &Checkbox{}
	c.MinPt = p
	c.Text = DrawText(s, mplusNormalFont, color.Black)
	return c
}

// todo: button과 한번에 합칠 수 있을까?
func (c *Checkbox) Update() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		p := image.Pt(ebiten.CursorPosition())
		if p.In(c.Rect()) {
			c.mouseDown = true
		} else {
			c.mouseDown = false
		}
	} else {
		if c.mouseDown {
			c.checked = !c.checked
			if c.onChanged != nil {
				c.onChanged(c)
			}
		}
		c.mouseDown = false
	}
}

// 패딩 넣어 그려진 체크박스 그리고 텍스트 이미지 그리기
func (c *Checkbox) Draw(screen *ebiten.Image) {
	var box *ebiten.Image
	if c.checked {
		box = imgCheckedBox
	} else {
		box = imgUncheckedBox
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.MinPt.X), float64(c.MinPt.Y))
	screen.DrawImage(box, op)

	bx, _ := imgUncheckedBox.Size()
	op.GeoM.Translate(float64(bx), 0)
	screen.DrawImage(c.Text, op)
}

func (c *Checkbox) Checked() bool {
	return c.checked
}

func (c *Checkbox) SetOnChanged(f func(c *Checkbox)) {
	c.onChanged = f
}

func (c *Checkbox) Rect() image.Rectangle {
	bw, bh := imgUncheckedBox.Size()
	tw, th := c.Text.Size()
	w := bw + tw
	h := th
	if h < bh {
		h = bh
	}
	maxPt := c.MinPt.Add(image.Pt(w, h))
	return image.Rectangle{c.MinPt, maxPt}
}

// func Rect(p image.Point, w, h int) image.Rectangle {
// 	return image.Rect(p.X, p.Y, p.X+w, p.Y+h)
// }
