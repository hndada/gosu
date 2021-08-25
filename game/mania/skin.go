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

	StageLeft   *ebiten.Image
	StageRight  *ebiten.Image
	StageHint   *ebiten.Image // todo: HitPosition 대신 필요할 듯
	StageBottom *ebiten.Image
	StageLight  *ebiten.Image // mask

	// MaskingBorder
	StageKeys        [4]*ebiten.Image
	StageKeysPressed [4]*ebiten.Image
}

func LoadSkin(cwd string) {
	// loadFont(cwd) // temp: load font

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
	// todo: mania-note1"H" and flip for tail
	// todo: animated sprites on LN
	path = filepath.Join(dir, "mania-note1L.png")
	Skin.LNBody[0], err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}
	path = filepath.Join(dir, "mania-note2L.png")
	Skin.LNBody[1], err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}
	path = filepath.Join(dir, "mania-noteSL.png")
	Skin.LNBody[2], err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}
	path = filepath.Join(dir, "mania-noteSCL.png")
	Skin.LNBody[3], err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}

	// judge
	path = filepath.Join(dir, "mania-hit300g.png")
	Skin.Judge[0], err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}
	path = filepath.Join(dir, "mania-hit300.png")
	Skin.Judge[1], err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}
	path = filepath.Join(dir, "mania-hit200.png")
	Skin.Judge[2], err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}
	path = filepath.Join(dir, "mania-hit50.png")
	Skin.Judge[3], err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}
	path = filepath.Join(dir, "mania-hit0.png")
	Skin.Judge[4], err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}

	// stage
	// key-hit은 기본 pressed, key-glow는 점수 나는 pressed인가?
	// todo: StageLeft가 없다면 StageRight 쓰게 하기
	path = filepath.Join(dir, "mania-stage-left.png")
	Skin.StageLeft, err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}
	path = filepath.Join(dir, "mania-stage-right.png")
	Skin.StageRight, err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}
	// path = filepath.Join(dir, "mania-stage-bottom.png")
	// Skin.StageBottom, err = game.LoadImage(path)
	// if err != nil {
	// 		log.Fatal(err)
	// }
	path = filepath.Join(dir, "mania-stage-light.png")
	Skin.StageLight, err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}
	path = filepath.Join(dir, "mania-stage-hint.png")
	Skin.StageHint, err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}

	path = filepath.Join(dir, "mania-key1.png")
	Skin.StageKeys[0], err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}
	path = filepath.Join(dir, "mania-key2.png")
	Skin.StageKeys[1], err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}
	path = filepath.Join(dir, "mania-keyS.png")
	Skin.StageKeys[2], err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}
	path = filepath.Join(dir, "mania-keyS.png") // temp: use keyS
	Skin.StageKeys[3], err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}

	path = filepath.Join(dir, "mania-key1D.png")
	Skin.StageKeysPressed[0], err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}
	path = filepath.Join(dir, "mania-key2D.png")
	Skin.StageKeysPressed[1], err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}
	path = filepath.Join(dir, "mania-keySD.png")
	Skin.StageKeysPressed[2], err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}
	path = filepath.Join(dir, "mania-keySD.png") // temp: use keyS
	Skin.StageKeysPressed[3], err = game.LoadImage(path)
	if err != nil {
		log.Fatal(err)
	}
}
