package mania

import (
	"errors"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/rg-parser/osugame/osu"
)

const (
	typeNote mode.NoteType = 1 << iota
	typeReleaseNote
	typeLongNote
)
const (
	TypeNote   = typeNote
	TypeLNHead = typeLongNote | typeNote
	TypeLNTail = typeLongNote | typeReleaseNote
)

// 난이도 및 점수 관련은 나중에
// 아래와 같이 난이도 계산에만 쓰이는 값들은 unexported로 할듯
type Note struct {
	mode.BaseNote
	Key int
	// Strain
	// Read
	// Stamina

	// hand        int
	// chord       []int
	// trillJack   []int
	// holdImpacts []float64
	//
	// strainBase   float64
	// chordPenalty float64
	// trillBonus   float64
	// jackBonus    float64
	// HoldBonus    float64 // score 필요함
}

var ErrDuration = errors.New("invalid duration: not a positive value")

func newNote(ho osu.HitObject, keys int) ([]Note, error) {
	ns := make([]Note, 0, 2)
	var n Note
	switch ho.NoteType & osu.ComboMask {
	case osu.HitTypeHoldNote:
		n.Type = typeLongNote
	case osu.HitTypeNote:
		n.Type = TypeNote
	default:
		return ns, errors.New("invalid hit object")
	}
	n.Key = ho.Column(keys)
	// n.Time = int64(float64(ho.Time) / mods.TimeRate)
	n.Time = int64(ho.Time)
	n.SampleFilename = ho.HitSample.Filename
	n.SampleVolume = uint8(ho.HitSample.Volume)

	if n.Type == typeLongNote {
		n.Type = TypeLNHead
		// n.Time2 = int64(float64(ho.EndTime) / mods.TimeRate)
		n.Time2 = int64(ho.EndTime)
		ns = append(ns, n)
		if n.Time2-n.Time <= 0 {
			return ns, ErrDuration
		}

		var n2 Note
		n2.Type = TypeLNTail
		n2.Key = n.Key
		n2.Time = n.Time2
		n2.Time2 = n.Time
		ns = append(ns, n2)
	} else {
		ns = append(ns, n)
	}
	return ns, nil
}

// func SortNotes(ns []Note) {
// 	sort.Slice(ns, func(i, j int) bool {
// 		if ns[i].Time == ns[j].Time {
// 			return ns[i].Key < ns[j].Key
// 		}
// 		return ns[i].Time < ns[j].Time
// 	})
// }

// func (n *Note) initSlices(keymode int) {
//	n.trillJack = tools.GetIntSlice(keymode, noFound)
//	n.chord = tools.GetIntSlice(keymode, noFound)
//	n.holdImpacts = make([]float64, keymode)
// }

// These values are applied at keys
// Example: 40 = 32 + 8 = Left-scratching 8 Key
const (
	ScratchLeft  = 1 << 5 // 32
	ScratchRight = 1 << 6 // 64
)
const ScratchMask = ^(ScratchLeft | ScratchRight)

type keyKind uint8

const (
	one keyKind = iota
	two
	middle
	pinky
)

var keyKinds = make(map[int][]keyKind)

func init() {
	keyKinds[0] = []keyKind{}
	keyKinds[1] = []keyKind{middle}
	keyKinds[2] = []keyKind{one, one}
	keyKinds[3] = []keyKind{one, middle, one}
	keyKinds[4] = []keyKind{one, two, two, one}
	keyKinds[5] = []keyKind{one, two, middle, two, one}
	keyKinds[6] = []keyKind{one, two, one, one, two, one}
	keyKinds[7] = []keyKind{one, two, one, middle, one, two, one}
	keyKinds[8] = []keyKind{pinky, one, two, one, one, two, one, pinky}
	keyKinds[9] = []keyKind{pinky, one, two, one, middle, one, two, one, pinky}
	keyKinds[10] = []keyKind{pinky, one, two, one, middle, middle, one, two, one, pinky}

	for i := 1; i <= 8; i++ { // 정말 잘 짠듯
		keyKinds[i|ScratchLeft] = append([]keyKind{pinky}, keyKinds[i-1]...)
		keyKinds[i|ScratchRight] = append(keyKinds[i-1], pinky)
	}
}