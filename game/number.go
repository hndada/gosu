package game

import "github.com/hajimehoshi/ebiten"

const (
	NumberCombo = iota
	NumberScore
)

func LoadNumbers(mode int) [10]Sprite {
	var scale float64
	var numbers [10]Sprite
	var height float64
	var position float64
	var srcs [13]*ebiten.Image
	switch mode {
	case NumberCombo:
		height = Settings.ComboHeight
		position = Settings.ComboPosition
		srcs = Skin.Number1
	case NumberScore:
		height = Settings.ScoreHeight
		position = height / 2
		srcs = Skin.Number2
	}
	for i := 0; i < 10; i++ {
		numbers[i].SetEbitenImage(srcs[i])
		numbers[i].H = int(height * DisplayScale())
		if i == 0 {
			scale = float64(numbers[i].H) / float64(srcs[i].Bounds().Size().Y)
		}
		numbers[i].W = int(float64(srcs[i].Bounds().Size().X) * scale)
		numbers[i].Y = int(position*DisplayScale()) - numbers[i].H/2
	}
	return numbers
}
