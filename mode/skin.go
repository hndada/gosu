package mode

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
)

// todo: backward-compatibility하면 게임 스킨 파싱 버전 안 적어도 되나?
// case1: 클라가 최신, 스킨이 구식 버전 (보통)
// case2: 스킨은 최신, 클라가 구식
// 버전 업데이트되는 경우는 기능의 추가 혹은 악용의 가능성으로 인한 삭제
// 추가야 backward 지원되고 삭제는 반영되어야 하는 부분이므로 업데이트가 강제되므로 스킨 버전은 따로 필요 없을 듯

// 스킨은 말그대로 이미지 struct; option 적용 안 돼있음
// Skin은 only for asset/resources; 설정은 전적으로 user에게 맡긴다
// author나 license는 Skin 폴더에 별도 텍스트 파일로 저장

// todo: 스킨 표준 세우기:
// score, combo: 0~9을 하나의 이미지로
// 4개짜리 이미지, 하나로 뭉치기
// 길이/4, 높이에 해당하는 거 만큼 롱노트 SubImage
type SkinTemplate struct {
	name       string
	score      [10]*ebiten.Image
	combo      [10]*ebiten.Image
	hpBarFrame *ebiten.Image
	hpBarColor *ebiten.Image
	// boxLeft         *ebiten.Image
	// boxMiddle       *ebiten.Image
	// boxRight        *ebiten.Image
	// chartPanelFrame *ebiten.Image
}

var Skin SkinTemplate

var defaultSkin SkinTemplate

// 게임 켜질 때
// 세팅의 Skin id 에 맞추어서 메모리에 올리기 (먼 훗날 gob로 저장)
// SkinPathList = map[int]string
// todo: filename list, .toml로 저장

func LoadSkin(skinPath string) {
	var filename, path string
	var ok bool
	for i := 0; i < 10; i++ {
		filename = fmt.Sprintf("score-%d.png", i)
		path = filepath.Join(skinPath, filename)
		if Skin.score[i], ok = LoadImage(path); !ok {
			Skin.score[i] = defaultSkin.score[i]
		}
	}
	for i := 0; i < 10; i++ {
		filename = fmt.Sprintf("combo-%d.png", i)
		path = filepath.Join(skinPath, filename)
		if Skin.combo[i], ok = LoadImage(path); !ok {
			Skin.combo[i] = defaultSkin.combo[i]
		}
	}
	filename = "scorebar-bg.png"
	path = filepath.Join(skinPath, filename)
	if Skin.hpBarFrame, ok = LoadImage(path); !ok {
		Skin.hpBarFrame = defaultSkin.hpBarFrame
	}
	filename = "scorebar-colour.png"
	path = filepath.Join(skinPath, filename)
	if Skin.hpBarColor, ok = LoadImage(path); !ok {
		Skin.hpBarColor = defaultSkin.hpBarColor
	}
	// filename = "button-left.png"
	// path = filepath.Join(skinPath, filename)
	// if Skin.boxLeft, ok = LoadImage(path); !ok {
	// 	Skin.boxLeft = defaultSkin.boxLeft
	// }
	// filename = "button-middle.png"
	// path = filepath.Join(skinPath, filename)
	// if Skin.boxMiddle, ok = LoadImage(path); !ok {
	// 	Skin.boxMiddle = defaultSkin.boxMiddle
	// }
	// filename = "button-right.png"
	// path = filepath.Join(skinPath, filename)
	// if Skin.boxRight, ok = LoadImage(path); !ok {
	// 	Skin.boxRight = defaultSkin.boxRight
	// }
	// filename = "menu-button-background.png"
	// path = filepath.Join(skinPath, filename)
	// if Skin.chartPanelFrame, ok = LoadImage(path); !ok {
	// 	Skin.chartPanelFrame = defaultSkin.chartPanelFrame
	// }
}

// todo: graphic로 이동
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
