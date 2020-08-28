package config

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"image"
	_ "image/png"
	"os"
	"path/filepath"
)

// todo: score, combo: 0~9을 하나의 이미지로
// todo: backward-compatibility하면 게임 스킨 파싱 버전 안 적어도 되나?
// case1: 클라가 최신, 스킨이 구식 버전 (보통)
// case2: 스킨은 최신, 클라가 구식

// skin은 only for asset/resources; 설정은 전적으로 user에게 맡긴다
// author나 license는 skin 폴더에 별도 텍스트 파일로 저장

// todo: 괜히 여기서 op 다 줄려 하지 말고, default x와 y값을 주자
// 여러 단계로 나누어서 x,y 진행
const (
	ScoreComma = iota + 10
	ScoreDot
	ScorePercent
)

// 중간에 render 필요한 애들만 image.Image로 불러오자
type Skin struct {
	Name            string
	Score           [13]*ebiten.Image // unscaled
	Combo           [10]*ebiten.Image // unscaled
	HPBarFrame      *ebiten.Image     // 폭맞춤x, screenHeigth
	HPBarColor      *ebiten.Image     // 폭맞춤x, screenHeigth
	BoxLeft         *ebiten.Image     // unscaled
	BoxMiddle       *ebiten.Image     // unscaled
	BoxRight        *ebiten.Image     // unscaled
	ChartPanelFrame *ebiten.Image     // unscaled

	Mania ManiaSkin
}

// screen에 처음부터 fixed 된 상태로 그려질 애들
// settings 건들 때마다 갱신
// 먼저 100 스케일로 그리고 확대하면 깨지니까
// todo: 판정선 가운데에 노트 가운데가 맞을 때 Max가 뜨게

type ManiaSkin struct {
	ManiaSkinImage
	ManiaSkinImageOp
}
type ManiaSkinImage struct {
	Note             [4]*ebiten.Image
	LNHead           [4]*ebiten.Image
	LNBody           [4][]*ebiten.Image
	LNTail           [4]*ebiten.Image
	KeyButton        [4]*ebiten.Image
	KeyButtonPressed [4]*ebiten.Image

	HitResults   [5]*ebiten.Image // unscaled
	NoteLighting []*ebiten.Image
	LNLighting   []*ebiten.Image
	StageRight   *ebiten.Image // 폭맞춤x, screenHeigth
	StageBottom  *ebiten.Image // fieldWidth, 폭맞춤 y
	StageHint    *ebiten.Image // fieldWidth, 설정값 ('노트와 동일한 높이로' 옵션 추가)

	Stage *ebiten.Image
	// StageOp *ebiten.DrawImageOptions
}

// s.noteWidths[7] = [4]float64{4.5, 4, 5, 5.5}

// ebiten.Image가 본체
// [4]*ebiten.Image
// [7][]*ebiten.DrawingImageOptions, 길이 7

// Info를 따로 분리하지 말고
// 로컬로 다룬 다음에 바로 Option 생성으로.

// 7키 이미지 (및 옵션) 저장
// noteImages
// lnHeadImages
// lnBodyImages -> 얘는 길이가 달라져야 하니까 음좀더 고민
// lnTailImages

// 실제 Draw할때는 screen.DrawImage(noteImage[note.Key])
// 실제 Draw할때는 noteImages[n.Key].DrawTo(screen)
type Image struct {
	raw *ebiten.Image
	op  *ebiten.DrawImageOptions
}

func (i *Image) DrawTo(dst *ebiten.Image) {
	dst.DrawImage(i.raw, i.op)
}

// 업데이트에서 dy만큼 이동
// 아니지 직접건들면 스킨값이 바뀌잖아
//
// x는 고정 가능
// y는 매번 업데이트 되어야함
// h는 처음에 정해지고 안바뀜
//
// 옵션은 값복사를 해야겠네

// 여기선 무대 오프셋은 적용 안함
// 맵 불러올때마다 하든 스킨 폴더에서 한 단계 더 거쳐서 하든지 하기로
func Op() ebiten.DrawImageOptions {
	var op ebiten.DrawImageOptions
	op.GeoM.Scale()
	op.GeoM.Translate()
	return op
}

// todo: ebiten.Image의 pointer 지우고 함수 만들기
type ManiaSkinImageOp struct {
	Note map[int][4]*ebiten.DrawImageOptions
}

func (s *Settings) SetNoteOp() {
	ops := make(map[int][4]*ebiten.DrawImageOptions)
	s.noteHeigth
	for key, ws := range s.noteWidths {
		scale := float64(s.ScreenSize().Y) / 100
		var op2 [4]*ebiten.DrawImageOptions
		op := &ebiten.DrawImageOptions{}
		for i, w := range ws {
			sw, sh := s.Skin.Mania.ManiaSkinImage.Note[i].Size()
			op.GeoM.Scale(w*scale/float64(sw), s.noteHeigth/float64(sh))
		}
		ops[i]
	}
	s.noteWidths
	s.Skin.Mania.ManiaSkinImageOp.Note = ops
}

// 이미지 소스 (ebiten.Image; 포인터 여부는 상관 없게. 메소드로 한번 감쌀 예정)
// 옵션
// 하나의 이미지에서 키마다 다른 옵션
func OpSetSize(p image.Point) *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	// 이미지 원래 크기를 알아야 하니까 메소드가 불가피
	// op.GeoM.Scale(w*scale/float64(sw), s.noteHeigth/float64(sh))

}
func (s *Settings) RenderManiaStage() {
	// 필드의 중앙이 스크린의 중앙에 오게 op.GeoM.Translate(dx, 0)
	// main *ebiten.image // fieldWidth, screenHeight (generated)
	op := &ebiten.DrawImageOptions{}
	stageCenter := float64(s.ScreenSize().X) * s.StagePosition / 100
	var fieldWidth float64

	op.GeoM.Translate(stageCenter, 0)
	op.GeoM.Translate(-fieldWidth/2, 0)
	s.Skin.Mania.Stage.DrawImage(i, op)
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
		s.LNBody[i] = make([]*ebiten.Image, 1, 1)
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

func loadSkinImage(filename string) (*ebiten.Image, error) {
	var skinPath = "C:\\Users\\hndada\\Documents\\GitHub\\hndada\\gosu\\asset\\Skin\\"
	path := filepath.Join(skinPath, filename)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	src, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	img, _ := ebiten.NewImageFromImage(src, ebiten.FilterDefault)
	return img, nil
}
