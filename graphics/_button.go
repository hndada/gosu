package graphic

import (
	"image/color"
)

const (
	bsDefault = iota
	bsPressed
	bsToggled
	bsPressedToggled
)

type Button struct {
	// widgetBase

	width  int
	height int
	state  int
	text   string

	// surfaces []surface
}

func NewButton() (*Button, error) {
	textColor := color.RGBA{R: 0x64, G: 0x64, B: 0x64, A: 0xff}
}
