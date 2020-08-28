package config

import (
	"github.com/hajimehoshi/ebiten"
	"image"
	"image/color"
)

// 매냐: 세로 100 기준으로 가로 얼마 만큼 쓸래
// EachDimness map[[16]byte]uint8 -> 얘는 toml 등으로 저장
// EachSpeed map[[16]byte]float64 -> 얘는 toml 등으로 저장
type ManiaSettings struct {
	Display      ManiaDisplaySettings // interface
	KeyLayout    map[int][]ebiten.Key // interface // todo: 무결성 검사, 겹치는거 있는지 매번 확인
	GeneralSpeed float64              // todo: fixed decimal?
	GroupSpeeds  []float64
}

type ManiaDisplaySettings struct {
	HitPosition       float64 // object which is now set at 'options'
	ComboPosition     float64
	HitResultPosition float64
	StagePosition     float64 // 0 ~ 100; 50 is a center

	NoteWidths     map[int][4]float64 // 키마다 width 설정
	NoteHeigth     float64            // 두께; 키 관계없이 동일
	lnHeadCustom   bool               // if false, head uses normal note image.
	lnTailMode     uint8              // 0: Tail=Head 1: Tail=Body 2: Custom
	lineInHint     bool
	SpotlightColor [4]color.RGBA

	// SplitGap            float64
	// UpsideDown          bool
	// ColumnDivisionWidth float64
}

// reset    save cancel
// 설정 켜면 임시 매니아디스플레이세팅 Struct 값복사로 생성
// 설정값 바꾸면 임시 세팅값이 바뀜, 실시간으로 보여주기-> refresh
// save 누르면 main 세팅 struct로 값복사
// refresh를 옵션 창에서 강제할까, struct 상에서 강제할까
func (s *ManiaDisplaySettings) refresh() {
}

const (
	lnTailModeHead = iota
	lnTailModeBody
	lnTailModeCustom
)

func (s *Settings) SetLNTailMode(mode uint8) {
	switch mode {
	case lnTailModeHead, lnTailModeBody, lnTailModeCustom:
		s.lnTailMode = mode
	default:
		return
	}
	switch mode {
	case lnTailModeHead:
		// 헤드 이미지 로드
	case lnTailModeBody:
		// 바디 이미지 로드
	case lnTailModeCustom:
	// 커스텀 이미지 로드
	default:
		return
	}
}

func (s *Settings) SetLNHeadCustom(set bool) {
	s.lnHeadCustom = set
	// 헤드 이미지 로드
}

type Kind int8

const (
	one Kind = iota
	two
	middle
	pinky
)

// applied at keys
// example: 40 = 32 + 8 = Left-scratching 8 Key
const (
	scratchLeft  = 1 << 5 // 32
	scratchRight = 1 << 6 // 64
	scratchMask  = ^(scratchLeft | scratchRight)
)

// 다른 패키지에서 수정하면 안됨
var NoteKinds = make(map[int][]Kind)

func init() {
	NoteKinds[0] = []Kind{}
	NoteKinds[1] = []Kind{middle}
	NoteKinds[2] = []Kind{one, one}
	NoteKinds[3] = []Kind{one, middle, one}
	NoteKinds[4] = []Kind{one, two, two, one}
	NoteKinds[5] = []Kind{one, two, middle, two, one}
	NoteKinds[6] = []Kind{one, two, one, one, two, one}
	NoteKinds[7] = []Kind{one, two, one, middle, one, two, one}
	NoteKinds[8] = []Kind{pinky, one, two, one, one, two, one, pinky}
	NoteKinds[9] = []Kind{pinky, one, two, one, middle, one, two, one, pinky}
	NoteKinds[10] = []Kind{pinky, one, two, one, middle, middle, one, two, one, pinky}

	for i := 1; i <= 8; i++ { // 정말 잘 짠듯
		NoteKinds[i|scratchLeft] = append([]Kind{pinky}, NoteKinds[i-1]...)
		NoteKinds[i|scratchRight] = append(NoteKinds[i-1], pinky)
	}
}

func (s *Settings) SetNoteWidths(keys int, vs [4]float64) {
	s.noteWidths[keys] = vs
	s.setNoteSizes()
}

func (s *Settings) SetNoteHeight(h float64) {
	s.noteHeigth = h
	s.setNoteSizes()
}

var noteSizeSamples = make(map[int][4]image.Point)

func (s *Settings) setNoteSizeSamples() {
	scale := float64(s.screenHeight) / 100
	for key, ws := range s.noteWidths {
		var samples [4]image.Point
		for kind := 0; kind < 4; kind++ {
			w := int(scale * ws[kind])
			h := int(scale * s.noteHeigth)
			samples[kind] = image.Pt(w, h)
		}
		noteSizeSamples[key] = samples
	}
}

// func NoteImages(keys int) []ebiten.Image {
// 	var s Skin
// 	imgs := make([]ebiten.Image, keys&scratchMask)
//
// 	for key, kind := range NoteKinds[keys] {
// 		src, _ := ebiten.NewImageFromImage(s.Mania.Note[kind], ebiten.FilterDefault)
// 		dstSize := noteSizeSamples[keys][kind]
// 		op := &ebiten.DrawImageOptions{}
// 		op.GeoM.Scale(ratio(dstSize, image.Pt(src.Size())))
// 		imgs[key] = *NewImage(dstSize.X, dstSize.Y)
// 		imgs[key].DrawImage(src, op)
// 	}
// 	return imgs
// }
//
// func NewImage(width, height int) *ebiten.Image {
// 	img, _ := ebiten.NewImage(width, height, ebiten.FilterDefault)
// 	return img
// }
//
// func ratio(dst, src image.Point) (float64, float64) {
// 	return float64(dst.X) / float64(src.X),
// 		float64(dst.Y) / float64(src.Y)
// }
