package mania

import (
	"log"
	"path/filepath"

	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/game"
)

// StartPoint, Width, Height, Name 총 4가지 알면 spritesheet 에서 이미지 빼올 수 있음
var Skin struct {
	Note [4]*ebiten.Image // one, two, middle, pinky
	// LNHead [4]*ebiten.Image
	LNBody [4]*ebiten.Image // animation
	// LNTail [4]*ebiten.Image   // never nil; LNHead 복사를 하더라도

	Judge      [5]*ebiten.Image
	Lighting   []*ebiten.Image
	LightingLN []*ebiten.Image

	StageLeft        *ebiten.Image
	StageRight       *ebiten.Image
	StageBottom      *ebiten.Image
	StageLight       *ebiten.Image // mask
	StageKeys        [4]*ebiten.Image
	StageKeysPressed [4]*ebiten.Image
	// MaskingBorder
	HPBar      *ebiten.Image
	HPBarColor *ebiten.Image
}

func LoadSkin(cwd string) {
	dir := filepath.Join(cwd, "skin")
	var path string
	var err error

	path = filepath.Join(dir, "mania-note1.png")
	Skin.Note[0], err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}
	path = filepath.Join(dir, "mania-note2.png")
	Skin.Note[1], err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}
	path = filepath.Join(dir, "mania-noteS.png")
	Skin.Note[2], err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}
	path = filepath.Join(dir, "mania-noteSC.png")
	Skin.Note[3], err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}

	// LN Body
	path = filepath.Join(dir, "mania-note1L-0.png")
	Skin.LNBody[0], err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}
	path = filepath.Join(dir, "mania-note2L-0.png")
	Skin.LNBody[1], err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}
	path = filepath.Join(dir, "mania-noteSL-0.png")
	Skin.LNBody[2], err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}
	path = filepath.Join(dir, "mania-noteSCL-0.png")
	Skin.LNBody[3], err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}
}
