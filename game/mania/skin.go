package mania

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/game"
	"path/filepath"
)

type skinTemplate struct {
	note             [4]*ebiten.Image
	lnHead           [4]*ebiten.Image
	lnBody           [4][]*ebiten.Image
	lnTail           [4]*ebiten.Image
	keyButton        [4]*ebiten.Image
	keyButtonPressed [4]*ebiten.Image

	hitResults   [5]*ebiten.Image
	noteLighting []*ebiten.Image
	lnLighting   []*ebiten.Image
	stageRight   *ebiten.Image // 폭맞춤x, screenHeigth
	stageBottom  *ebiten.Image // fieldWidth, 폭맞춤 y
	stageHint    *ebiten.Image // fieldWidth, 설정값 ('노트와 동일한 높이로' 옵션 추가)
}

var skin skinTemplate
var defaultSkin skinTemplate

// todo: for loop으로 1/4로 줄일 수 있음
func loadSkin(skinPath string) {
	var filename, path string
	var ok bool
	for i, fname := range []string{"mania-note1.png", "mania-note2.png", "mania-noteS.png", "mania-noteSC.png"} {
		path = filepath.Join(skinPath, fname)
		if skin.note[i], ok = game.LoadImage(path); !ok {
			skin.note[i] = defaultSkin.note[i]
		}
	}
	skin.lnHead = skin.note
	// filename =
	// 	path = filepath.Join(skinPath, filename)
	// if skin.lnHead[0], ok = graphics.LoadImage(path); !ok {
	// 	skin.lnHead[0] = defaultSkin.lnHead[0]
	// }
	// filename =
	// 	path = filepath.Join(skinPath, filename)
	// if skin.lnHead[1], ok = graphics.LoadImage(path); !ok {
	// 	skin.lnHead[1] = defaultSkin.lnHead[1]
	// }
	// filename =
	// 	path = filepath.Join(skinPath, filename)
	// if skin.lnHead[2], ok = graphics.LoadImage(path); !ok {
	// 	skin.lnHead[2] = defaultSkin.lnHead[2]
	// }
	// filename =
	// 	path = filepath.Join(skinPath, filename)
	// if skin.lnHead[3], ok = graphics.LoadImage(path); !ok {
	// 	skin.lnHead[3] = defaultSkin.lnHead[3]
	// }

	for i := range skin.lnBody {
		skin.lnBody[i] = make([]*ebiten.Image, 1)
	}
	filename = "mania-note1L-0.png"
	path = filepath.Join(skinPath, filename)
	if skin.lnBody[0][0], ok = game.LoadImage(path); !ok {
		skin.lnBody[0][0] = defaultSkin.lnBody[0][0]
	}
	filename = "mania-note2L-0.png"
	path = filepath.Join(skinPath, filename)
	if skin.lnBody[1][0], ok = game.LoadImage(path); !ok {
		skin.lnBody[1][0] = defaultSkin.lnBody[1][0]
	}
	filename = "mania-noteSL-0.png"
	path = filepath.Join(skinPath, filename)
	if skin.lnBody[2][0], ok = game.LoadImage(path); !ok {
		skin.lnBody[2][0] = defaultSkin.lnBody[2][0]
	}
	filename = "mania-noteSCL-0.png"
	path = filepath.Join(skinPath, filename)
	if skin.lnBody[3][0], ok = game.LoadImage(path); !ok {
		skin.lnBody[3][0] = defaultSkin.lnBody[3][0]

	}
	// skin 에서는 setting 상태와 관계 없이 있는 대로 이미지를 불러온다
	skin.lnTail = skin.lnHead
	// filename =
	// 	path = filepath.Join(skinPath, filename)
	// if skin.lnTail[0], ok = graphics.LoadImage(path); !ok {
	// 	skin.lnTail[0] = defaultSkin.lnTail[0]
	// }
	// filename =
	// 	path = filepath.Join(skinPath, filename)
	// if skin.lnTail[1], ok = graphics.LoadImage(path); !ok {
	// 	skin.lnTail[1] = defaultSkin.lnTail[1]
	// }
	// filename =
	// 	path = filepath.Join(skinPath, filename)
	// if skin.lnTail[2], ok = graphics.LoadImage(path); !ok {
	// 	skin.lnTail[2] = defaultSkin.lnTail[2]
	// }
	// filename =
	// 	path = filepath.Join(skinPath, filename)
	// if skin.lnTail[3], ok = graphics.LoadImage(path); !ok {
	// 	skin.lnTail[3] = defaultSkin.lnTail[3]
	// }

	filename = "mania-key1.png"
	path = filepath.Join(skinPath, filename)
	if skin.keyButton[0], ok = game.LoadImage(path); !ok {
		skin.keyButton[0] = defaultSkin.keyButton[0]
	}
	filename = "mania-key2.png"
	path = filepath.Join(skinPath, filename)
	if skin.keyButton[1], ok = game.LoadImage(path); !ok {
		skin.keyButton[1] = defaultSkin.keyButton[1]
	}
	filename = "mania-keyS.png"
	path = filepath.Join(skinPath, filename)
	if skin.keyButton[2], ok = game.LoadImage(path); !ok {
		skin.keyButton[2] = defaultSkin.keyButton[2]
	}
	filename = "mania-keyS.png" // todo: 4th image
	path = filepath.Join(skinPath, filename)
	if skin.keyButton[3], ok = game.LoadImage(path); !ok {
		skin.keyButton[3] = defaultSkin.keyButton[3]
	}
	filename = "mania-key1D.png"
	path = filepath.Join(skinPath, filename)
	if skin.keyButtonPressed[0], ok = game.LoadImage(path); !ok {
		skin.keyButtonPressed[0] = defaultSkin.keyButtonPressed[0]
	}
	filename = "mania-key2D.png"
	path = filepath.Join(skinPath, filename)
	if skin.keyButtonPressed[1], ok = game.LoadImage(path); !ok {
		skin.keyButtonPressed[1] = defaultSkin.keyButtonPressed[1]
	}
	filename = "mania-keySD.png"
	path = filepath.Join(skinPath, filename)
	if skin.keyButtonPressed[2], ok = game.LoadImage(path); !ok {
		skin.keyButtonPressed[2] = defaultSkin.keyButtonPressed[2]
	}
	filename = "mania-keySD.png" // todo: 4th image
	path = filepath.Join(skinPath, filename)
	if skin.keyButtonPressed[3], ok = game.LoadImage(path); !ok {
		skin.keyButtonPressed[3] = defaultSkin.keyButtonPressed[3]
	}
	filename = "mania-hit300g.png"
	path = filepath.Join(skinPath, filename)
	if skin.hitResults[0], ok = game.LoadImage(path); !ok {
		skin.hitResults[0] = defaultSkin.hitResults[0]
	}
	filename = "mania-hit300.png"
	path = filepath.Join(skinPath, filename)
	if skin.hitResults[1], ok = game.LoadImage(path); !ok {
		skin.hitResults[1] = defaultSkin.hitResults[1]
	}
	filename = "mania-hit200.png"
	path = filepath.Join(skinPath, filename)
	if skin.hitResults[2], ok = game.LoadImage(path); !ok {
		skin.hitResults[2] = defaultSkin.hitResults[2]
	}
	filename = "mania-hit50.png"
	path = filepath.Join(skinPath, filename)
	if skin.hitResults[3], ok = game.LoadImage(path); !ok {
		skin.hitResults[3] = defaultSkin.hitResults[3]
	}
	filename = "mania-hit0.png"
	path = filepath.Join(skinPath, filename)
	if skin.hitResults[4], ok = game.LoadImage(path); !ok {
		skin.hitResults[4] = defaultSkin.hitResults[4]
	}
	skin.noteLighting = make([]*ebiten.Image, 1)
	filename = "lightingN-3.png"
	path = filepath.Join(skinPath, filename)
	if skin.noteLighting[0], ok = game.LoadImage(path); !ok {
		skin.noteLighting[0] = defaultSkin.noteLighting[0]
	}
	skin.lnLighting = make([]*ebiten.Image, 1)
	filename = "lightingL-0.png"
	path = filepath.Join(skinPath, filename)
	if skin.lnLighting[0], ok = game.LoadImage(path); !ok {
		skin.lnLighting[0] = defaultSkin.lnLighting[0]
	}
	filename = "mania-stage-right.png"
	path = filepath.Join(skinPath, filename)
	if skin.stageRight, ok = game.LoadImage(path); !ok {
		skin.stageRight = defaultSkin.stageRight
	}
	filename = "mania-stage-bottom.png"
	path = filepath.Join(skinPath, filename)
	if skin.stageBottom, ok = game.LoadImage(path); !ok {
		skin.stageBottom = defaultSkin.stageBottom
	}
	filename = "mania-stage-hint.png"
	path = filepath.Join(skinPath, filename)
	if skin.stageHint, ok = game.LoadImage(path); !ok {
		skin.stageHint = defaultSkin.stageHint
	}
}
