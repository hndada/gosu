package config

import (
	"github.com/hajimehoshi/ebiten"
	"image"
	"image/color"
)

// 매냐: 세로 100 기준으로 가로 얼마 만큼 쓸래
// 값 변경과 동시에 실행되어야 하는 함수가 있는 경우 private/method로 set하는 방법으로 변경
type Settings struct {
	// General - screen
	maxTPS int
	// screenWidth  int
	// screenHeight int
	screenSize image.Point

	// General - sound
	volumeMaster uint8
	volumeBGM    uint8
	volumeSFX    uint8

	// Play - screen
	DimValue    uint8
	ScrollSpeed float64

	// Mania
	ManiaKeyLayout    map[int][]ebiten.Key
	HitPosition       float64 // object which is now set at 'options'
	ComboPosition     float64
	HitResultPosition float64

	noteWidths      map[int][4]float64 // 키마다 width 설정
	noteHeigth      float64            // 두께; 키 관계없이 동일
	lnHeadCustom    bool               // if false, head uses normal note image.
	lnTailMode      uint8              // 0: Tail=Head 1: Tail=Body 2: Custom
	lineInJudgeLine bool
	SpotlightColor  [4]color.RGBA

	SplitGap            float64
	UpsideDown          bool
	StagePosition       float64 // -100 ~ 100; 0 is a center
	ColumnDivisionWidth float64 // default is 0.
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

func (s *Settings) SetLineInJudgeLine(set bool) {
	s.lineInJudgeLine = set
	// judgeline에 흰선 그리기
}

func (s *Settings) SetLNHeadCustom(set bool) {
	s.lnHeadCustom = set
	// 헤드 이미지 로드
}

func percent(v uint8) float64 { return float64(v) / 100 }

// Streamer 2개 만들기: BGM, SFX
func (s *Settings) SetDownVolumeMaster() {
	switch {
	case s.volumeMaster <= 0:
		return
	case s.volumeMaster <= 5:
		s.volumeMaster -= 1
	case s.volumeMaster <= 100:
		s.volumeMaster -= 5
	}
	// Streamer에 vol 꽂기
}

// func (s *Settings) SetVolumeBGM(vol uint8) {
// 	s.volumeBGM = vol
// 	// percent(s.volumeBGM) * percent(s.volumeMaster)
// }
//
// func (s *Settings) SetVolumeSFX(vol uint8) {
// 	s.volumeSFX = vol
// 	// Streamer에 vol 꽂기
// 	// percent(s.volumeSFX) * percent(s.volumeMaster)
// }

func (s *Settings) ApplyDim(op *ebiten.DrawImageOptions) {
	op.ColorM.ChangeHSV(0, 1, float64(s.DimValue))
}

// 배경 스케일 자동 맞춤
// func (s *Settings) AdjustBg(op *ebiten.DrawImageOptions) {
// }

// 로컬db 만들고 나면 loading 구현
func LoadSettings() Settings {
	return newSettings()
}

func newSettings() Settings {
	s := Settings{
		maxTPS:     240,
		screenSize: image.Pt(1600, 900),
		// screenWidth:       1600,
		// screenHeight:      900,
		DimValue:          25,
		volumeMaster:      100,
		volumeBGM:         50,
		volumeSFX:         50,
		ScrollSpeed:       1.33,
		ManiaKeyLayout:    make(map[int][]ebiten.Key),
		HitPosition:       70,
		ComboPosition:     50,
		HitResultPosition: 60,
		noteWidths:        make(map[int][4]float64),
		noteHeigth:        3,
		lnHeadCustom:      false,
		lnTailMode:        0,
		SpotlightColor: [4]color.RGBA{
			{64, 0, 0, 64},
			{0, 0, 64, 64},
			{64, 48, 0, 64},
			{40, 0, 40, 64},
		},
		lineInJudgeLine:     true,
		SplitGap:            0,
		UpsideDown:          false,
		StagePosition:       0,
		ColumnDivisionWidth: 0,
	}
	s.ManiaKeyLayout[4] = []ebiten.Key{
		ebiten.KeyD, ebiten.KeyF, ebiten.KeyJ, ebiten.KeyK,
	}
	s.ManiaKeyLayout[7] = []ebiten.Key{
		ebiten.KeyS, ebiten.KeyD, ebiten.KeyF,
		ebiten.KeySpace, ebiten.KeyJ, ebiten.KeyK, ebiten.KeyL,
	}

	s.noteWidths[4] = [4]float64{10, 9, 11, 12}
	s.noteWidths[7] = [4]float64{4.5, 4, 5, 5.5} // [4]float64{4.67, 3.83, 5.5, 5.5}
	return s
}

const (
	even int = iota
	odd
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
var NoteKinds = make(map[int][]int)

func init() {
	NoteKinds[0] = []int{}
	NoteKinds[1] = []int{middle}
	NoteKinds[2] = []int{even, even}
	NoteKinds[3] = []int{even, middle, even}
	NoteKinds[4] = []int{even, odd, odd, even}
	NoteKinds[5] = []int{even, odd, middle, odd, even}
	NoteKinds[6] = []int{even, odd, even, even, odd, even}
	NoteKinds[7] = []int{even, odd, even, middle, even, odd, even}
	NoteKinds[8] = []int{pinky, even, odd, even, even, odd, even, pinky}
	NoteKinds[9] = []int{pinky, even, odd, even, middle, even, odd, even, pinky}
	NoteKinds[10] = []int{pinky, even, odd, even, middle, middle, even, odd, even, pinky}

	for i := 1; i <= 8; i++ { // 정말 잘 짠듯
		NoteKinds[i|scratchLeft] = append([]int{pinky}, NoteKinds[i-1]...)
		NoteKinds[i|scratchRight] = append(NoteKinds[i-1], pinky)
	}
}

func (s *Settings) SetMaxTPS(tps int) {
	s.maxTPS = tps
	ebiten.SetMaxTPS(tps)
}

func (s *Settings) MaxTPS() int   { return s.maxTPS }
func (s *Settings) Tick() float64 { return 1000 / float64(s.maxTPS) } // 1000ms
func (s *Settings) SetScreenSize(p image.Point) {
	// s.screenWidth = w
	// s.screenHeight = h
	s.screenSize = p
	s.setNoteSizes()
	ebiten.SetWindowSize(p.X, p.Y)
}
func (s *Settings) ScreenSize() image.Point { return s.screenSize }

func (s *Settings) SetNoteWidths(keys int, vs [4]float64) {
	s.noteWidths[keys] = vs
	s.setNoteSizes()
}

func (s *Settings) SetNoteHeight(h float64) {
	s.noteHeigth = h
	s.setNoteSizes()
}

// 4개만 보내면 키버튼 타입 레인지 돌면서 옵션으로 스케일 걸고 그리게 될 것
// 옵션 리로드할 때만 호출
// todo: test
// todo: should be field of settings
// 얘도 4개씩만 저장해야됨
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

func NoteImages(keys int) []ebiten.Image {
	var s Skin
	imgs := make([]ebiten.Image, keys&scratchMask)

	for key, kind := range NoteKinds[keys] {
		src, _ := ebiten.NewImageFromImage(s.Mania.Note[kind], ebiten.FilterDefault)
		dstSize := noteSizeSamples[keys][kind]
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(ratio(dstSize, image.Pt(src.Size())))
		imgs[key] = *NewImage(dstSize.X, dstSize.Y)
		imgs[key].DrawImage(src, op)
	}
	return imgs
}

func NewImage(width, height int) *ebiten.Image {
	img, _ := ebiten.NewImage(width, height, ebiten.FilterDefault)
	return img
}

func ratio(dst, src image.Point) (float64, float64) {
	return float64(dst.X) / float64(src.X),
		float64(dst.Y) / float64(src.Y)
}
