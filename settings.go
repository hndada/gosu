package gosu

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/mode/mania"
	"image"
	"image/color"
)

// 매냐: 세로 100 기준으로 가로 얼마 만큼 쓸래
type Settings struct {
	// Screen
	MaxTPS       int
	ScreenWidth  int
	ScreenHeight int
	DimValue     int

	// Sound
	VolumeSFX int
	VolumeBGM int

	// Mania
	ScrollSpeed       float64
	ManiaKeyLayout    map[int][]ebiten.Key
	HitPosition       float64 // object which is now set at 'options'
	ComboPosition     float64
	HitResultPosition float64

	NoteWidth       map[int][4]float64 // 키마다 width 설정
	NoteHeigth      float64            // 두께; 키 관계없이 동일
	LNHeadCustom    bool               // if false, head uses normal note image.
	LNTailMode      uint8              // 0: Tail=Head 1: Tail=Body 2: Custom
	SpotlightColor  [4]color.RGBA
	LineInJudgeLine bool

	SplitGap            float64
	UpsideDown          bool
	FieldCenterPosition float64 // -100 ~ 100; 0 is a center
	ColumnDivisionWidth float64 // default is 0.
}

// 로컬db 만들고 나면 loading 구현
func LoadSettings() Settings {
	return DefaultSettings()
}

func DefaultSettings() Settings {
	s := Settings{
		MaxTPS:            240,
		ScreenWidth:       1600,
		ScreenHeight:      900,
		DimValue:          25,
		VolumeSFX:         50,
		VolumeBGM:         50,
		ScrollSpeed:       1.33,
		ManiaKeyLayout:    make(map[int][]ebiten.Key),
		HitPosition:       70,
		ComboPosition:     50,
		HitResultPosition: 60,
		NoteWidth:         make(map[int][4]float64),
		NoteHeigth:        3,
		LNHeadCustom:      false,
		LNTailMode:        0,
		SpotlightColor: [4]color.RGBA{
			{64, 0, 0, 64},
			{0, 0, 64, 64},
			{64, 48, 0, 64},
			{40, 0, 40, 64},
		},
		LineInJudgeLine:     true,
		SplitGap:            0,
		UpsideDown:          false,
		FieldCenterPosition: 0,
		ColumnDivisionWidth: 0,
	}
	s.ManiaKeyLayout[4] = []ebiten.Key{
		ebiten.KeyD, ebiten.KeyF, ebiten.KeyJ, ebiten.KeyK,
	}
	s.ManiaKeyLayout[7] = []ebiten.Key{
		ebiten.KeyS, ebiten.KeyD, ebiten.KeyF,
		ebiten.KeySpace, ebiten.KeyJ, ebiten.KeyK, ebiten.KeyL,
	}

	s.NoteWidth[4] = [4]float64{10, 9, 11, 12}
	s.NoteWidth[7] = [4]float64{4.5, 4, 5, 5.5} // [4]float64{4.67, 3.83, 5.5, 5.5}
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
)

var KeyButtonType = make(map[int][]int)

func init() {
	KeyButtonType[0] = []int{}
	KeyButtonType[1] = []int{middle}
	KeyButtonType[2] = []int{even, even}
	KeyButtonType[3] = []int{even, middle, even}
	KeyButtonType[4] = []int{even, odd, odd, even}
	KeyButtonType[5] = []int{even, odd, middle, odd, even}
	KeyButtonType[6] = []int{even, odd, even, even, odd, even}
	KeyButtonType[7] = []int{even, odd, even, middle, even, odd, even}
	KeyButtonType[8] = []int{pinky, even, odd, even, even, odd, even, pinky}
	KeyButtonType[9] = []int{pinky, even, odd, even, middle, even, odd, even, pinky}
	KeyButtonType[10] = []int{pinky, even, odd, even, middle, middle, even, odd, even, pinky}

	for i := 1; i <= 8; i++ { // 정말 잘 짠듯
		KeyButtonType[i|scratchLeft] = append([]int{pinky}, KeyButtonType[i-1]...)
		KeyButtonType[i|scratchRight] = append(KeyButtonType[i-1], pinky)
	}
}

// 전체 리스트를 보낼지, 4개만 알짜로 보낼지.
// 4개만 보내면 키버튼 타입 레인지 돌면서 옵션으로 스케일 걸고 그리게 될 것
// 아니면 얘도 옵션 리로드할 때만 호출하고 계속 재활용하자
func (g *Game) NoteSizes(keys int) []image.Point {
	ps := make([]image.Point, keys)
	samples := g.Settings.NoteWidth[keys]
	scale := float64(g.ScreenHeight) / 100
	for i, t := range KeyButtonType[keys] {
		w := int(scale * samples[t])
		h := int(scale * g.NoteHeigth)
		ps[i] = image.Pt(w, h)
	}
	return ps
}

// temp
func noteColor(n mania.Note, keys int) color.RGBA {
	switch n.Key {
	case 0, 2, 4, 6:
		return color.RGBA{239, 243, 247, 0xff} // white
	case 1, 5:
		return color.RGBA{66, 211, 247, 0xff} // blue
	case 3:
		return color.RGBA{255, 203, 82, 0xff} // yellow
	}
	panic("not reach")
}
