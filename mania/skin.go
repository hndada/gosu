package mania

import (
	"fmt"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/engine/ui"
)

// StartPoint, Width, Height, Name 총 4가지 알면 spritesheet 에서 이미지 빼올 수 있음
var Skin struct {
	Note   [4]*ebiten.Image // one, two, middle, pinky
	LNHead [4]*ebiten.Image // optional
	LNBody [4]*ebiten.Image // animation
	LNTail [4]*ebiten.Image // optional

	Judge      [5]*ebiten.Image
	Lighting   []*ebiten.Image
	LightingLN []*ebiten.Image

	StageLeft  *ebiten.Image
	StageRight *ebiten.Image
	StageHint  *ebiten.Image // todo: HitPosition 대신 필요할 듯
	// StageBottom *ebiten.Image
	StageLight *ebiten.Image // mask

	// MaskingBorder
	StageKeys        [4]*ebiten.Image
	StageKeysPressed [4]*ebiten.Image

	// Mania mode should have its own HPBar image: rotated version
	HPBar      *ebiten.Image
	HPBarColor *ebiten.Image // todo: Animation
}

func LoadSkin(cwd string) {
	dir := filepath.Join(cwd, "skin")
	var path string
	var err error

	path = filepath.Join(dir, "mania-note1.png")
	Skin.Note[0], err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}
	path = filepath.Join(dir, "mania-note2.png")
	Skin.Note[1], err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}
	path = filepath.Join(dir, "mania-noteS.png")
	Skin.Note[2], err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}
	path = filepath.Join(dir, "mania-noteSC.png")
	Skin.Note[3], err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}

	// LNHead: Note image is used when fails at loading
	path = filepath.Join(dir, "mania-note1H.png")
	Skin.LNHead[0], err = ui.LoadImageHD(path)
	if err != nil {
		Skin.LNHead[0] = Skin.Note[0]
	}
	path = filepath.Join(dir, "mania-note2H.png")
	Skin.LNHead[1], err = ui.LoadImageHD(path)
	if err != nil {
		Skin.LNHead[1] = Skin.Note[1]
	}
	path = filepath.Join(dir, "mania-noteSH.png")
	Skin.LNHead[2], err = ui.LoadImageHD(path)
	if err != nil {
		Skin.LNHead[2] = Skin.Note[2]
	}
	path = filepath.Join(dir, "mania-noteSCH.png")
	Skin.LNHead[3], err = ui.LoadImageHD(path)
	if err != nil {
		Skin.LNHead[3] = Skin.Note[3]
	}
	// LN Body // todo: animated sprites on LN
	path = filepath.Join(dir, "mania-note1L.png")
	Skin.LNBody[0], err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}
	path = filepath.Join(dir, "mania-note2L.png")
	Skin.LNBody[1], err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}
	path = filepath.Join(dir, "mania-noteSL.png")
	Skin.LNBody[2], err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}
	path = filepath.Join(dir, "mania-noteSCL.png")
	Skin.LNBody[3], err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}
	// LNHead image is used when fails at loading
	// Note image is used when fails even at loading LNHead image
	{
		path = filepath.Join(dir, "mania-note1T.png")
		i1, err := ui.LoadImageImage(path)
		if err != nil {
			path = filepath.Join(dir, "mania-note1H.png")
			i1, err = ui.LoadImageImage(path)
			if err != nil {
				path = filepath.Join(dir, "mania-note1.png")
				i1, err = ui.LoadImageImage(path)
				if err != nil {
					panic(err)
				}
			}
		}
		i2 := imaging.FlipV(i1)
		i3 := ebiten.NewImageFromImage(i2)
		Skin.LNTail[0] = i3
	}
	{
		path = filepath.Join(dir, "mania-note2T.png")
		i1, err := ui.LoadImageImage(path)
		if err != nil {
			path = filepath.Join(dir, "mania-note2H.png")
			i1, err = ui.LoadImageImage(path)
			if err != nil {
				path = filepath.Join(dir, "mania-note2.png")
				i1, err = ui.LoadImageImage(path)
				if err != nil {
					panic(err)
				}
			}
		}
		i2 := imaging.FlipV(i1)
		i3 := ebiten.NewImageFromImage(i2)
		Skin.LNTail[1] = i3
	}
	{
		path = filepath.Join(dir, "mania-noteST.png")
		i1, err := ui.LoadImageImage(path)
		if err != nil {
			path = filepath.Join(dir, "mania-noteSH.png")
			i1, err = ui.LoadImageImage(path)
			if err != nil {
				path = filepath.Join(dir, "mania-noteS.png")
				i1, err = ui.LoadImageImage(path)
				if err != nil {
					panic(err)
				}
			}
		}
		i2 := imaging.FlipV(i1)
		i3 := ebiten.NewImageFromImage(i2)
		Skin.LNTail[2] = i3
	}
	{
		path = filepath.Join(dir, "mania-noteSCT.png")
		i1, err := ui.LoadImageImage(path)
		if err != nil {
			path = filepath.Join(dir, "mania-noteSCH.png")
			i1, err = ui.LoadImageImage(path)
			if err != nil {
				path = filepath.Join(dir, "mania-noteSC.png")
				i1, err = ui.LoadImageImage(path)
				if err != nil {
					panic(err)
				}
			}
		}
		i2 := imaging.FlipV(i1)
		i3 := ebiten.NewImageFromImage(i2)
		Skin.LNTail[3] = i3
	}
	// judge
	path = filepath.Join(dir, "mania-hit300g.png")
	Skin.Judge[0], err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}
	path = filepath.Join(dir, "mania-hit300.png")
	Skin.Judge[1], err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}
	path = filepath.Join(dir, "mania-hit200.png")
	Skin.Judge[2], err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}
	path = filepath.Join(dir, "mania-hit50.png")
	Skin.Judge[3], err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}
	path = filepath.Join(dir, "mania-hit0.png")
	Skin.Judge[4], err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}

	// stage
	// key-hit은 기본 pressed, key-glow는 점수 나는 pressed인가?
	// todo: StageLeft가 없다면 StageRight 쓰게 하기
	path = filepath.Join(dir, "mania-stage-left.png")
	Skin.StageLeft, err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}
	path = filepath.Join(dir, "mania-stage-right.png")
	Skin.StageRight, err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}
	// path = filepath.Join(dir, "mania-stage-bottom.png")
	// Skin.StageBottom, err = ui.LoadImageHD(path)
	// if err != nil {
	// 		panic(err)
	// }
	path = filepath.Join(dir, "mania-stage-light.png")
	Skin.StageLight, err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}
	path = filepath.Join(dir, "mania-stage-hint.png")
	Skin.StageHint, err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}

	path = filepath.Join(dir, "mania-key1.png")
	Skin.StageKeys[0], err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}
	path = filepath.Join(dir, "mania-key2.png")
	Skin.StageKeys[1], err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}
	path = filepath.Join(dir, "mania-keyS.png")
	Skin.StageKeys[2], err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}
	path = filepath.Join(dir, "mania-keyS.png") // temp: use keyS
	Skin.StageKeys[3], err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}

	path = filepath.Join(dir, "mania-key1D.png")
	Skin.StageKeysPressed[0], err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}
	path = filepath.Join(dir, "mania-key2D.png")
	Skin.StageKeysPressed[1], err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}
	path = filepath.Join(dir, "mania-keySD.png")
	Skin.StageKeysPressed[2], err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}
	path = filepath.Join(dir, "mania-keySD.png") // temp: use keyS
	Skin.StageKeysPressed[3], err = ui.LoadImageHD(path)
	if err != nil {
		panic(err)
	}

	var name string
	var img *ebiten.Image
	var count int
	Skin.Lighting = make([]*ebiten.Image, 0, 10)
	for {
		name = fmt.Sprintf("lightingN-%d.png", count)
		path = filepath.Join(dir, name)
		img, err = ui.LoadImageHD(path)
		if err != nil {
			break
		} else {
			Skin.Lighting = append(Skin.Lighting, img)
		}
		count++
	}
	if len(Skin.Lighting) == 0 {
		path = filepath.Join(dir, "lightingN.png")
		img, err = ui.LoadImageHD(path)
		if err != nil {
			panic(err)
		}
		Skin.Lighting = append(Skin.Lighting, img)
	}
	count = 0
	Skin.LightingLN = make([]*ebiten.Image, 0, 10)
	for {
		name = fmt.Sprintf("lightingL-%d.png", count)
		path = filepath.Join(dir, name)
		img, err = ui.LoadImageHD(path)
		if err != nil {
			break
		} else {
			Skin.LightingLN = append(Skin.LightingLN, img)
		}
		count++
	}
	if len(Skin.LightingLN) == 0 {
		path = filepath.Join(dir, "lightingL.png")
		img, err = ui.LoadImageHD(path)
		if err != nil {
			panic(err)
		}
		Skin.LightingLN = append(Skin.LightingLN, img)
	}
	{
		path = filepath.Join(dir, "scorebar-bg.png")
		i1, err := ui.LoadImageImage(path)
		if err != nil {
			panic(err)
		}
		i2 := imaging.Rotate90(i1)
		i3 := ebiten.NewImageFromImage(i2)
		Skin.HPBar = i3
	}
	{
		path = filepath.Join(dir, "scorebar-colour.png")
		i1, err := ui.LoadImageImage(path)
		if err != nil {
			panic(err)
		}
		i2 := imaging.Rotate90(i1)
		i3 := ebiten.NewImageFromImage(i2)
		Skin.HPBarColor = i3
	}
}
