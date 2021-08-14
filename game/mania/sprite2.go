package mania

import "image"

var Skin2 struct {
	Note   [4]image.Image // one, two, middle, pinky
	LNHead [4]image.Image
	LNBody [4]image.Image
	LNTail [4]image.Image // never nil; LNHead 복사를 하더라도

	Judge      [5]image.Image
	Lighting   []image.Image
	LightingLN []image.Image

	StageLeft        image.Image
	StageRight       image.Image
	StageBottom      image.Image
	StageLight       image.Image // mask
	StageKeys        [4]image.Image
	StageKeysPressed [4]image.Image
	// MaskingBorder
	HPBar      image.Image
	HPBarColor image.Image
}
