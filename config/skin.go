package config

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
	"path/filepath"
)

// todo: backward-compatibility하면 게임 스킨 파싱 버전 안 적어도 되나?
// case1: 클라가 최신, 스킨이 구식 버전 (보통)
// case2: 스킨은 최신, 클라가 구식

// skin은 only for asset/resources; 설정은 전적으로 user에게 맡긴다
// author나 license는 skin 폴더에 별도 텍스트 파일로 저장

const (
	ScoreComma = iota + 10
	ScoreDot
	ScorePercent
)

// todo: score, combo: 0~9을 하나의 이미지로
// todo: image.Image는 결국 필요 없는거 아닌가? - 중간에 render 필요한 애들만 image로 불러오자
type Skin struct {
	Name            string
	Score           [13]image.Image // unscaled
	Combo           [10]image.Image // unscaled
	HPBarFrame      image.Image
	HPBarColor      image.Image
	BoxLeft         image.Image // unscaled
	BoxMiddle       image.Image // unscaled
	BoxRight        image.Image // unscaled
	ChartPanelFrame image.Image // unscaled

	Mania ManiaSkin
}

type ManiaSkin struct {
	Note             [4]image.Image
	LNHead           [4]image.Image
	LNBody           [4][]image.Image
	LNTail           [4]image.Image
	KeyButton        [4]image.Image
	KeyButtonPressed [4]image.Image

	HitResults   [5]image.Image
	NoteLighting []image.Image
	LNLighting   []image.Image
	StageRight   image.Image
	StageBottom  image.Image
	StageHint    image.Image

	Stage ManiaStage
}

// todo: TOML-ize
func LoadSkin() Skin {
	var s Skin
	var err error

	for i := range s.Score {
		var word string
		switch i {
		case ScoreComma:
			word = "comma"
		case ScoreDot:
			word = "dot"
		case ScorePercent:
			word = "percent"
		default:
			word = fmt.Sprint(i)
		}
		filename := fmt.Sprintf("score-%s.png", word)
		if s.Score[i], err = loadSkinImage(filename); err != nil {
			continue
		}
	}
	for i := range s.Combo {
		filename := fmt.Sprintf("combo-%d.png", i)
		if s.Combo[i], err = loadSkinImage(filename); err != nil {
			continue
		}
	}
	s.HPBarFrame, _ = loadSkinImage("scorebar-bg.png")
	s.HPBarColor, _ = loadSkinImage("scorebar-colour.png")
	s.BoxLeft, _ = loadSkinImage("button-left.png")
	s.BoxMiddle, _ = loadSkinImage("button-middle.png")
	s.BoxRight, _ = loadSkinImage("button-right.png")
	s.ChartPanelFrame, _ = loadSkinImage("menu-button-background.png")

	s.Mania = loadManiaSkin()
	return s
}

func loadManiaSkin() ManiaSkin {
	var s ManiaSkin
	s.Note[0], _ = loadSkinImage("mania-note1.png")
	s.Note[1], _ = loadSkinImage("mania-note2.png")
	s.Note[2], _ = loadSkinImage("mania-noteS.png")
	s.Note[3], _ = loadSkinImage("mania-noteSC.png")

	s.LNHead = s.Note
	// s.LNHead[0] = s.Note[0]
	// s.LNHead[1] = s.Note[1]
	// s.LNHead[2] = s.Note[2]
	// s.LNHead[3] = s.Note[3]

	for i := range s.LNBody {
		s.LNBody[i] = make([]image.Image, 1, 1)
	}
	s.LNBody[0][0], _ = loadSkinImage("mania-note1.png")
	s.LNBody[1][0], _ = loadSkinImage("mania-note2.png")
	s.LNBody[2][0], _ = loadSkinImage("mania-noteS.png")
	s.LNBody[3][0], _ = loadSkinImage("mania-noteSC.png")

	s.LNTail = s.LNHead

	s.KeyButton[0], _ = loadSkinImage("mania-key1.png")
	s.KeyButton[1], _ = loadSkinImage("mania-key2.png")
	s.KeyButton[2], _ = loadSkinImage("mania-keyS.png")
	s.KeyButtonPressed[0], _ = loadSkinImage("mania-key1D.png")
	s.KeyButtonPressed[1], _ = loadSkinImage("mania-key2D.png")
	s.KeyButtonPressed[2], _ = loadSkinImage("mania-keySD.png")

	s.HitResults[0], _ = loadSkinImage("mania-hit300g.png")
	s.HitResults[1], _ = loadSkinImage("mania-hit300.png")
	s.HitResults[2], _ = loadSkinImage("mania-hit200.png")
	s.HitResults[3], _ = loadSkinImage("mania-hit50.png")
	s.HitResults[4], _ = loadSkinImage("mania-hit0.png")

	s.StageRight, _ = loadSkinImage("mania-stage-right.png")
	s.StageBottom, _ = loadSkinImage("mania-stage-bottom.png")
	s.StageHint, _ = loadSkinImage("mania-stage-hint.png")

	return s
}

func loadSkinImage(filename string) (image.Image, error) {
	var skinPath = "C:\\Users\\hndada\\Documents\\GitHub\\hndada\\gosu\\asset\\Skin\\"
	path := filepath.Join(skinPath, filename)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, nil
}
