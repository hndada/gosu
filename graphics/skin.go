package graphics

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"image"
	"os"
	"path/filepath"
)

// todo: backward-compatibility하면 게임 스킨 파싱 버전 안 적어도 되나?
// case1: 클라가 최신, 스킨이 구식 버전 (보통)
// case2: 스킨은 최신, 클라가 구식
// 버전 업데이트되는 경우는 기능의 추가 혹은 악용의 가능성으로 인한 삭제
// 추가야 backward 지원되고 삭제는 반영되어야 하는 부분이므로 업데이트가 강제되므로 스킨 버전은 따로 필요 없을 듯

// 스킨은 말그대로 이미지 struct
// skin은 only for asset/resources; 설정은 전적으로 user에게 맡긴다
// author나 license는 skin 폴더에 별도 텍스트 파일로 저장

// todo: 스킨 표준
// score, combo: 0~9을 하나의 이미지로
// 4개짜리 이미지, 하나로 뭉치기
// 길이/4, 높이에 해당하는 거 만큼 롱노트 SubImage
type skin struct {
	name            string
	score           [10]*ebiten.Image
	combo           [10]*ebiten.Image
	hpBarFrame      *ebiten.Image
	hpBarColor      *ebiten.Image
	boxLeft         *ebiten.Image
	boxMiddle       *ebiten.Image
	boxRight        *ebiten.Image
	chartPanelFrame *ebiten.Image

	mania maniaSkin
}

type maniaSkin struct {
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

// 게임 켜질 때
// 세팅의 skin id 에 맞추어서 메모리에 올리기 아니면 gob로 저장
// skinPathList = map[int]string
// todo: filename list, .toml로 저장
func (s *skin) LoadSkin(skinPath string) {
	var filename, path string
	var ok bool
	for i := 0; i < 10; i++ {
		filename = fmt.Sprintf("score-%d.png", i)
		path = filepath.Join(skinPath, filename)
		if s.score[i], ok = LoadImage(path); !ok {
			s.score[i] = defaultSkin.score[i]
		}
	}
	for i := 0; i < 10; i++ {
		filename = fmt.Sprintf("combo-%d.png", i)
		path = filepath.Join(skinPath, filename)
		if s.combo[i], ok = LoadImage(path); !ok {
			s.combo[i] = defaultSkin.combo[i]
		}
	}
	filename = "scorebar-bg.png"
	path = filepath.Join(skinPath, filename)
	if s.hpBarFrame, ok = LoadImage(path); !ok {
		s.hpBarFrame = defaultSkin.hpBarFrame
	}
	filename = "scorebar-colour.png"
	path = filepath.Join(skinPath, filename)
	if s.hpBarColor, ok = LoadImage(path); !ok {
		s.hpBarColor = defaultSkin.hpBarColor
	}
	filename = "button-left.png"
	path = filepath.Join(skinPath, filename)
	if s.boxLeft, ok = LoadImage(path); !ok {
		s.boxLeft = defaultSkin.boxLeft
	}
	filename = "button-middle.png"
	path = filepath.Join(skinPath, filename)
	if s.boxMiddle, ok = LoadImage(path); !ok {
		s.boxMiddle = defaultSkin.boxMiddle
	}
	filename = "button-right.png"
	path = filepath.Join(skinPath, filename)
	if s.boxRight, ok = LoadImage(path); !ok {
		s.boxRight = defaultSkin.boxRight
	}
	filename = "menu-button-background.png"
	path = filepath.Join(skinPath, filename)
	if s.chartPanelFrame, ok = LoadImage(path); !ok {
		s.chartPanelFrame = defaultSkin.chartPanelFrame
	}
	s.mania.load(skinPath)
}

func (s *maniaSkin) load(skinPath string) {
	var filename, path string
	var ok bool
	filename = "mania-note1.png"
	path = filepath.Join(skinPath, filename)
	if s.note[0], ok = LoadImage(path); !ok {
		s.note[0] = defaultSkin.mania.note[0]
	}
	filename = "mania-note2.png"
	path = filepath.Join(skinPath, filename)
	if s.note[1], ok = LoadImage(path); !ok {
		s.note[1] = defaultSkin.mania.note[1]
	}
	filename = "mania-noteS.png"
	path = filepath.Join(skinPath, filename)
	if s.note[2], ok = LoadImage(path); !ok {
		s.note[2] = defaultSkin.mania.note[2]
	}
	filename = "mania-noteSC.png"
	path = filepath.Join(skinPath, filename)
	if s.note[3], ok = LoadImage(path); !ok {
		s.note[3] = defaultSkin.mania.note[3]
	}

	s.lnHead = s.note
	// filename =
	// 	path = filepath.Join(skinPath, filename)
	// if s.lnHead[0], ok = LoadImage(path); !ok {
	// 	s.lnHead[0] = defaultSkin.mania.lnHead[0]
	// }
	// filename =
	// 	path = filepath.Join(skinPath, filename)
	// if s.lnHead[1], ok = LoadImage(path); !ok {
	// 	s.lnHead[1] = defaultSkin.mania.lnHead[1]
	// }
	// filename =
	// 	path = filepath.Join(skinPath, filename)
	// if s.lnHead[2], ok = LoadImage(path); !ok {
	// 	s.lnHead[2] = defaultSkin.mania.lnHead[2]
	// }
	// filename =
	// 	path = filepath.Join(skinPath, filename)
	// if s.lnHead[3], ok = LoadImage(path); !ok {
	// 	s.lnHead[3] = defaultSkin.mania.lnHead[3]
	// }

	for i := range s.lnBody {
		s.lnBody[i] = make([]*ebiten.Image, 1)
	}
	filename = "mania-note1L-0"
	path = filepath.Join(skinPath, filename)
	if s.lnBody[0][0], ok = LoadImage(path); !ok {
		s.lnBody[0][0] = defaultSkin.mania.lnBody[0][0]
	}
	filename = "mania-note2L-0"
	path = filepath.Join(skinPath, filename)
	if s.lnBody[1][0], ok = LoadImage(path); !ok {
		s.lnBody[1][0] = defaultSkin.mania.lnBody[1][0]
	}
	filename = "mania-noteSL-0"
	path = filepath.Join(skinPath, filename)
	if s.lnBody[2][0], ok = LoadImage(path); !ok {
		s.lnBody[2][0] = defaultSkin.mania.lnBody[2][0]
	}
	filename = "mania-noteSCL-0"
	path = filepath.Join(skinPath, filename)
	if s.lnBody[3][0], ok = LoadImage(path); !ok {
		s.lnBody[3][0] = defaultSkin.mania.lnBody[3][0]

	}
	// skin 에서는 setting 상태와 관계 없이 있는 대로 이미지를 불러온다
	s.lnTail = s.lnHead
	// filename =
	// 	path = filepath.Join(skinPath, filename)
	// if s.lnTail[0], ok = LoadImage(path); !ok {
	// 	s.lnTail[0] = defaultSkin.mania.lnTail[0]
	// }
	// filename =
	// 	path = filepath.Join(skinPath, filename)
	// if s.lnTail[1], ok = LoadImage(path); !ok {
	// 	s.lnTail[1] = defaultSkin.mania.lnTail[1]
	// }
	// filename =
	// 	path = filepath.Join(skinPath, filename)
	// if s.lnTail[2], ok = LoadImage(path); !ok {
	// 	s.lnTail[2] = defaultSkin.mania.lnTail[2]
	// }
	// filename =
	// 	path = filepath.Join(skinPath, filename)
	// if s.lnTail[3], ok = LoadImage(path); !ok {
	// 	s.lnTail[3] = defaultSkin.mania.lnTail[3]
	// }

	filename = "mania-key1.png"
	path = filepath.Join(skinPath, filename)
	if s.keyButton[0], ok = LoadImage(path); !ok {
		s.keyButton[0] = defaultSkin.mania.keyButton[0]
	}
	filename = "mania-key2.png"
	path = filepath.Join(skinPath, filename)
	if s.keyButton[1], ok = LoadImage(path); !ok {
		s.keyButton[1] = defaultSkin.mania.keyButton[1]
	}
	filename = "mania-keyS.png"
	path = filepath.Join(skinPath, filename)
	if s.keyButton[2], ok = LoadImage(path); !ok {
		s.keyButton[2] = defaultSkin.mania.keyButton[2]
	}
	filename = "mania-keyS.png" // todo: 4th image
	path = filepath.Join(skinPath, filename)
	if s.keyButton[3], ok = LoadImage(path); !ok {
		s.keyButton[3] = defaultSkin.mania.keyButton[3]
	}
	filename = "mania-key1D.png"
	path = filepath.Join(skinPath, filename)
	if s.keyButtonPressed[0], ok = LoadImage(path); !ok {
		s.keyButtonPressed[0] = defaultSkin.mania.keyButtonPressed[0]
	}
	filename = "mania-key2D.png"
	path = filepath.Join(skinPath, filename)
	if s.keyButtonPressed[1], ok = LoadImage(path); !ok {
		s.keyButtonPressed[1] = defaultSkin.mania.keyButtonPressed[1]
	}
	filename = "mania-keySD.png"
	path = filepath.Join(skinPath, filename)
	if s.keyButtonPressed[2], ok = LoadImage(path); !ok {
		s.keyButtonPressed[2] = defaultSkin.mania.keyButtonPressed[2]
	}
	filename = "mania-keySD.png" // todo: 4th image
	path = filepath.Join(skinPath, filename)
	if s.keyButtonPressed[3], ok = LoadImage(path); !ok {
		s.keyButtonPressed[3] = defaultSkin.mania.keyButtonPressed[3]
	}
	filename = "mania-hit300g.png"
	path = filepath.Join(skinPath, filename)
	if s.hitResults[0], ok = LoadImage(path); !ok {
		s.hitResults[0] = defaultSkin.mania.hitResults[0]
	}
	filename = "mania-hit300.png"
	path = filepath.Join(skinPath, filename)
	if s.hitResults[1], ok = LoadImage(path); !ok {
		s.hitResults[1] = defaultSkin.mania.hitResults[1]
	}
	filename = "mania-hit200.png"
	path = filepath.Join(skinPath, filename)
	if s.hitResults[2], ok = LoadImage(path); !ok {
		s.hitResults[2] = defaultSkin.mania.hitResults[2]
	}
	filename = "mania-hit50.png"
	path = filepath.Join(skinPath, filename)
	if s.hitResults[3], ok = LoadImage(path); !ok {
		s.hitResults[3] = defaultSkin.mania.hitResults[3]
	}
	filename = "mania-hit0.png"
	path = filepath.Join(skinPath, filename)
	if s.hitResults[4], ok = LoadImage(path); !ok {
		s.hitResults[4] = defaultSkin.mania.hitResults[4]
	}
	s.noteLighting = make([]*ebiten.Image, 1)
	filename = "lightingN-3.png"
	path = filepath.Join(skinPath, filename)
	if s.noteLighting[0], ok = LoadImage(path); !ok {
		s.noteLighting[0] = defaultSkin.mania.noteLighting[0]
	}
	s.lnLighting = make([]*ebiten.Image, 1)
	filename = "lightingL-0.png"
	path = filepath.Join(skinPath, filename)
	if s.lnLighting[0], ok = LoadImage(path); !ok {
		s.lnLighting[0] = defaultSkin.mania.lnLighting[0]
	}
	filename = "mania-stage-right.png"
	path = filepath.Join(skinPath, filename)
	if s.stageRight, ok = LoadImage(path); !ok {
		s.stageRight = defaultSkin.mania.stageRight
	}
	filename = "mania-stage-bottom.png"
	path = filepath.Join(skinPath, filename)
	if s.stageBottom, ok = LoadImage(path); !ok {
		s.stageBottom = defaultSkin.mania.stageBottom
	}
	filename = "mania-stage-hint.png"
	path = filepath.Join(skinPath, filename)
	if s.stageHint, ok = LoadImage(path); !ok {
		s.stageHint = defaultSkin.mania.stageHint
	}
}

// loadSkinImage로 한번에 표시하려면 reflect 써야함
func LoadImage(path string) (*ebiten.Image, bool) {
	empty, _ := ebiten.NewImage(0, 0, ebiten.FilterDefault)
	f, err := os.Open(path)
	if err != nil {
		return empty, false
	}
	defer f.Close()
	src, _, err := image.Decode(f)
	if err != nil {
		return empty, false
	}
	img, _ := ebiten.NewImageFromImage(src, ebiten.FilterDefault)
	return img, true
}

// 정확도를 표시하지 않을 거므로 필요 없음
// const (
// 	ScoreComma = iota + 10
// 	ScoreDot
// 	ScorePercent
// )
