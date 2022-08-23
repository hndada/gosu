package drum

import "github.com/hndada/gosu/format/osu"

type NoteType int

const (
	don = 1 << iota
	kat
	head // Head of roll note, aka slider.
	tail // Tail of roll note, aka slider.
)
const (
	normal = 8 << iota
	big
	shake // Also known as spinner.
)
const (
	Don     NoteType = normal + don
	Kat     NoteType = normal + kat
	Head    NoteType = normal + head
	Tail    NoteType = normal + tail
	BigDon  NoteType = big + don
	BigKat  NoteType = big + kat
	BigHead NoteType = big + head
	BigTail NoteType = big + tail
	Shake   NoteType = shake
)

type Note struct {
	Type      NoteType
	Time      int64
	Time2     int64   // Time of opposite note's time, if existed.
	Speed     float64 // Each note has own speed.
	ScaledBPM float64 // For calculating Roll tick density.
	// BeatScale      float64
	SampleFilename string
	SampleVolume   float64 // Range is 0 to 1.
}

// EndTime does not work in Roll.
func NewNote(f any, speed, scaledBPM float64) []Note {
	ns := make([]Note, 0, 2)
	switch f := f.(type) {
	case osu.HitObject:
		n := Note{
			Type:           DrumNoteTypeFromOsu(f),
			Time:           int64(f.Time),
			Time2:          int64(f.Time), // Should not rely on f.EndTime
			Speed:          speed,
			ScaledBPM:      scaledBPM,
			SampleFilename: f.HitSample.Filename,
			SampleVolume:   float64(f.HitSample.Volume) / 100,
		}
		if IsRoll(f) {
			sp := f.SliderParams
			length := float64(sp.Slides) * sp.Length
			n.Time2 = n.Time + int64(length/speed)
			n2 := n
			if n.Type == Head {
				n2.Type = Tail
			} else {
				n2.Type = BigTail
			}
			n2.Time = n.Time2
			n2.Time2 = n.Time
			n2.SampleFilename = "" // Tail has no sample sound.
			if n.Type == BigHead {
				n2.Type = BigTail
			}
			ns = append(ns, n, n2)
			// n2 := Note{
			// 	Type:  Tail,
			// 	Time:  n.Time2,
			// 	Time2: n.Time,
			// }
		} else {
			ns = append(ns, n)
		}
	}
	return ns
}

func DrumNoteTypeFromOsu(h osu.HitObject) NoteType {
	if h.NoteType&osu.HitTypeSpinner != 0 {
		return Shake
	}
	if osu.IsBig(h) {
		switch {
		case IsRoll(h):
			return BigHead
		case osu.IsDon(h):
			return BigDon
		case osu.IsKat(h):
			return BigKat
		}
	}
	switch {
	case IsRoll(h):
		return Head
	case osu.IsDon(h):
		return Don
	case osu.IsKat(h):
		return Kat
	default:
		return Don
	}
}
func IsRoll(h osu.HitObject) bool { return h.NoteType&osu.HitTypeSlider != 0 }

// func IsRoll2(h osu.HitObject) bool { return h.NoteType&osu.ComboMask == osu.HitTypeSlider }

const (
	MaxScaledBPM = 256
	MinScaledBPM = 128
)

// It is proved that all BPMs are set into [128, 256) by v*2 or v/2.
func ScaledBPM(bpm float64) float64 {
	if bpm < 0 {
		bpm = -bpm
	}
	switch {
	case bpm >= 256:
		for bpm >= 256 {
			bpm /= 2
		}
	case bpm >= 128:
		return bpm
	case bpm < 128:
		for bpm < 128 {
			bpm *= 2
		}
	}
	return bpm
}

func (n Note) IsBig() bool { return n.Type&big != 0 }
func (n Note) NoteKind() int { // Todo: NoteKind -> NoteType?
	if n.IsBig() {
		return BigNote
	}
	return NormalNote
}

// func RollDuration(hs osu.SliderParams, bpm float64, multiplier float64) int64 {
// 	length := float64(hs.Slides) * hs.Length
// 	speed := (bpm / 60000) * (multiplier * 100) // Unit: amount of osu!pixel per 100ms
// 	return int64(length / speed)
// }

// func RollDuration(h osu.HitObject, speed float64) int {
// 	if h.NoteType&osu.HitTypeSlider == 0 {
// 		return 0
// 	}
// 	hs := h.SliderParams
// 	length := float64(hs.Slides) * hs.Length
// 	// speed := (bpm / 60000) * beatScale * (multiplier * 100)
// 	return int(length / speed)
// }
