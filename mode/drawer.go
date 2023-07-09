package mode

import (
	"math"

	"github.com/hndada/gosu/draws"
)

type ScoreDrawer struct {
	digitWidth float64 // Use number 0's width.
	DigitGap   float64
	ZeroFill   int
	Score      Delayed
	Sprites    []draws.Sprite
}

func NewScoreDrawer(sprites []draws.Sprite, gap float64) ScoreDrawer {
	return ScoreDrawer{
		digitWidth: sprites[0].W(),
		DigitGap:   gap,
		ZeroFill:   1,
		Score:      NewDelayed(),
		Sprites:    sprites[:10],
	}
}
func (d *ScoreDrawer) Update(score float64) {
	d.Score.Update(score)
}

func (d ScoreDrawer) Draw(dst draws.Image) {
	vs := make([]int, 0)
	score := int(math.Floor(d.Score.Delayed + 0.1))
	for v := score; v > 0; v /= 10 {
		vs = append(vs, v%10) // Little endian.
	}
	for i := len(vs); i < d.ZeroFill; i++ {
		vs = append(vs, 0)
	}
	w := d.digitWidth + d.DigitGap
	var tx float64
	for _, v := range vs {
		sprite := d.Sprites[v]
		sprite.Move(tx, 0)
		sprite.Move(-w/2+sprite.W()/2, 0) // Need to set at center since anchor is RightTop.
		sprite.Draw(dst, draws.Op{})
		tx -= w
	}
}

// Todo: add combo *int and skip passing combo value?
type ComboDrawer struct {
	draws.Timer
	DigitWidth float64 // Use number 0's width.
	DigitGap   float64
	Combo      int
	Bounce     float64
	Sprites    [10]draws.Sprite
}

// Each number has different width. Number 0's width is used as standard.
func (d *ComboDrawer) Update(combo int) {
	d.Ticker()
	if d.Combo != combo {
		d.Combo = combo
		d.Timer.Reset()
	}
}

// ComboDrawer's Draw draws each number at constant x regardless of their widths.
func (d ComboDrawer) Draw(dst draws.Image) {
	if d.IsDone() {
		return
	}
	if d.Combo == 0 {
		return
	}
	vs := make([]int, 0)
	for v := d.Combo; v > 0; v /= 10 {
		vs = append(vs, v%10) // Little endian.
	}

	// Size of the whole image is 0.5w + (n-1)(w+gap) + 0.5w.
	// Since sprites are already at anchor, no need to care of two 0.5w.
	w := d.DigitWidth + d.DigitGap
	tx := float64(len(vs)-1) * w / 2
	const (
		boundary1 = 0.05
		boundary2 = 0.1
	)
	for _, v := range vs {
		sprite := d.Sprites[v]
		sprite.Move(tx, 0)
		age := d.Age()
		if age < boundary1 {
			scale := 0.1 * d.Progress(0, boundary1)
			sprite.Move(0, d.Bounce*sprite.H()*scale)
		}
		if age >= boundary1 && age < boundary2 {
			scale := 0.1 - 0.1*d.Progress(boundary1, boundary2)
			sprite.Move(0, d.Bounce*sprite.H()*scale)
		}
		sprite.Draw(dst, draws.Op{})
		tx -= w
	}
}
