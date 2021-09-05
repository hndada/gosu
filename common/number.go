package common

import "github.com/hajimehoshi/ebiten/v2"

const (
	NumberCombo = iota
	NumberScore
)

func LoadNumbers(mode int) [10]Sprite {
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
		numbers[i].SetImage(srcs[i])
		numbers[i].H = int(height * DisplayScale())
		// set scale on every image: each image may has different size
		scale := float64(numbers[i].H) / float64(srcs[i].Bounds().Size().Y)
		numbers[i].W = int(float64(srcs[i].Bounds().Size().X) * scale)
		numbers[i].Y = int(position*DisplayScale()) - numbers[i].H/2
		numbers[i].Saturation = 1
		numbers[i].Dimness = 1
	}
	return numbers
}
