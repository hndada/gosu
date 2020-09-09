package mania

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/mode"
	"path/filepath"
)

type skin struct {
	*mode.CommonSkin
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

var defaultSkin skin

// todo: for loop으로 1/4로 줄일 수 있음
func (s *skin) load(skinPath string) {
	var filename, path string
	var ok bool
	for i, fname := range []string{"mania-note1.png", "mania-note2.png", "mania-noteS.png", "mania-noteSC.png"} {
		path = filepath.Join(skinPath, fname)
		if s.note[i], ok = mode.LoadImage(path); !ok {
			s.note[i] = defaultSkin.note[i]
		}
	}
	s.lnHead = s.note
	// filename =
	// 	path = filepath.Join(skinPath, filename)
	// if s.lnHead[0], ok = graphics.LoadImage(path); !ok {
	// 	s.lnHead[0] = defaultSkin.lnHead[0]
	// }
	// filename =
	// 	path = filepath.Join(skinPath, filename)
	// if s.lnHead[1], ok = graphics.LoadImage(path); !ok {
	// 	s.lnHead[1] = defaultSkin.lnHead[1]
	// }
	// filename =
	// 	path = filepath.Join(skinPath, filename)
	// if s.lnHead[2], ok = graphics.LoadImage(path); !ok {
	// 	s.lnHead[2] = defaultSkin.lnHead[2]
	// }
	// filename =
	// 	path = filepath.Join(skinPath, filename)
	// if s.lnHead[3], ok = graphics.LoadImage(path); !ok {
	// 	s.lnHead[3] = defaultSkin.lnHead[3]
	// }

	for i := range s.lnBody {
		s.lnBody[i] = make([]*ebiten.Image, 1)
	}
	filename = "mania-note1L-0.png"
	path = filepath.Join(skinPath, filename)
	if s.lnBody[0][0], ok = mode.LoadImage(path); !ok {
		s.lnBody[0][0] = defaultSkin.lnBody[0][0]
	}
	filename = "mania-note2L-0.png"
	path = filepath.Join(skinPath, filename)
	if s.lnBody[1][0], ok = mode.LoadImage(path); !ok {
		s.lnBody[1][0] = defaultSkin.lnBody[1][0]
	}
	filename = "mania-noteSL-0.png"
	path = filepath.Join(skinPath, filename)
	if s.lnBody[2][0], ok = mode.LoadImage(path); !ok {
		s.lnBody[2][0] = defaultSkin.lnBody[2][0]
	}
	filename = "mania-noteSCL-0.png"
	path = filepath.Join(skinPath, filename)
	if s.lnBody[3][0], ok = mode.LoadImage(path); !ok {
		s.lnBody[3][0] = defaultSkin.lnBody[3][0]

	}
	// skin 에서는 setting 상태와 관계 없이 있는 대로 이미지를 불러온다
	s.lnTail = s.lnHead
	// filename =
	// 	path = filepath.Join(skinPath, filename)
	// if s.lnTail[0], ok = graphics.LoadImage(path); !ok {
	// 	s.lnTail[0] = defaultSkin.lnTail[0]
	// }
	// filename =
	// 	path = filepath.Join(skinPath, filename)
	// if s.lnTail[1], ok = graphics.LoadImage(path); !ok {
	// 	s.lnTail[1] = defaultSkin.lnTail[1]
	// }
	// filename =
	// 	path = filepath.Join(skinPath, filename)
	// if s.lnTail[2], ok = graphics.LoadImage(path); !ok {
	// 	s.lnTail[2] = defaultSkin.lnTail[2]
	// }
	// filename =
	// 	path = filepath.Join(skinPath, filename)
	// if s.lnTail[3], ok = graphics.LoadImage(path); !ok {
	// 	s.lnTail[3] = defaultSkin.lnTail[3]
	// }

	filename = "mania-key1.png"
	path = filepath.Join(skinPath, filename)
	if s.keyButton[0], ok = mode.LoadImage(path); !ok {
		s.keyButton[0] = defaultSkin.keyButton[0]
	}
	filename = "mania-key2.png"
	path = filepath.Join(skinPath, filename)
	if s.keyButton[1], ok = mode.LoadImage(path); !ok {
		s.keyButton[1] = defaultSkin.keyButton[1]
	}
	filename = "mania-keyS.png"
	path = filepath.Join(skinPath, filename)
	if s.keyButton[2], ok = mode.LoadImage(path); !ok {
		s.keyButton[2] = defaultSkin.keyButton[2]
	}
	filename = "mania-keyS.png" // todo: 4th image
	path = filepath.Join(skinPath, filename)
	if s.keyButton[3], ok = mode.LoadImage(path); !ok {
		s.keyButton[3] = defaultSkin.keyButton[3]
	}
	filename = "mania-key1D.png"
	path = filepath.Join(skinPath, filename)
	if s.keyButtonPressed[0], ok = mode.LoadImage(path); !ok {
		s.keyButtonPressed[0] = defaultSkin.keyButtonPressed[0]
	}
	filename = "mania-key2D.png"
	path = filepath.Join(skinPath, filename)
	if s.keyButtonPressed[1], ok = mode.LoadImage(path); !ok {
		s.keyButtonPressed[1] = defaultSkin.keyButtonPressed[1]
	}
	filename = "mania-keySD.png"
	path = filepath.Join(skinPath, filename)
	if s.keyButtonPressed[2], ok = mode.LoadImage(path); !ok {
		s.keyButtonPressed[2] = defaultSkin.keyButtonPressed[2]
	}
	filename = "mania-keySD.png" // todo: 4th image
	path = filepath.Join(skinPath, filename)
	if s.keyButtonPressed[3], ok = mode.LoadImage(path); !ok {
		s.keyButtonPressed[3] = defaultSkin.keyButtonPressed[3]
	}
	filename = "mania-hit300g.png"
	path = filepath.Join(skinPath, filename)
	if s.hitResults[0], ok = mode.LoadImage(path); !ok {
		s.hitResults[0] = defaultSkin.hitResults[0]
	}
	filename = "mania-hit300.png"
	path = filepath.Join(skinPath, filename)
	if s.hitResults[1], ok = mode.LoadImage(path); !ok {
		s.hitResults[1] = defaultSkin.hitResults[1]
	}
	filename = "mania-hit200.png"
	path = filepath.Join(skinPath, filename)
	if s.hitResults[2], ok = mode.LoadImage(path); !ok {
		s.hitResults[2] = defaultSkin.hitResults[2]
	}
	filename = "mania-hit50.png"
	path = filepath.Join(skinPath, filename)
	if s.hitResults[3], ok = mode.LoadImage(path); !ok {
		s.hitResults[3] = defaultSkin.hitResults[3]
	}
	filename = "mania-hit0.png"
	path = filepath.Join(skinPath, filename)
	if s.hitResults[4], ok = mode.LoadImage(path); !ok {
		s.hitResults[4] = defaultSkin.hitResults[4]
	}
	s.noteLighting = make([]*ebiten.Image, 1)
	filename = "lightingN-3.png"
	path = filepath.Join(skinPath, filename)
	if s.noteLighting[0], ok = mode.LoadImage(path); !ok {
		s.noteLighting[0] = defaultSkin.noteLighting[0]
	}
	s.lnLighting = make([]*ebiten.Image, 1)
	filename = "lightingL-0.png"
	path = filepath.Join(skinPath, filename)
	if s.lnLighting[0], ok = mode.LoadImage(path); !ok {
		s.lnLighting[0] = defaultSkin.lnLighting[0]
	}
	filename = "mania-stage-right.png"
	path = filepath.Join(skinPath, filename)
	if s.stageRight, ok = mode.LoadImage(path); !ok {
		s.stageRight = defaultSkin.stageRight
	}
	filename = "mania-stage-bottom.png"
	path = filepath.Join(skinPath, filename)
	if s.stageBottom, ok = mode.LoadImage(path); !ok {
		s.stageBottom = defaultSkin.stageBottom
	}
	filename = "mania-stage-hint.png"
	path = filepath.Join(skinPath, filename)
	if s.stageHint, ok = mode.LoadImage(path); !ok {
		s.stageHint = defaultSkin.stageHint
	}
}
