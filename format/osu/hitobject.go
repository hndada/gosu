package osu

import (
	"fmt"
	"strconv"
	"strings"
)

type HitObject struct { // delimiter,
	X            int          `json:"x"`
	Y            int          `json:"y"`
	Time         int          `json:"time"`
	NoteType     int          `json:"noteType"` // nofloat
	HitSound     int          `json:"hitSound"` // nofloat
	EndTime      int          `json:"endTime"`  // optional
	SliderParams SliderParams // optional
	HitSample    HitSample    // optional
}
type SliderParams struct { // delimiter,
	CurveType   string   `json:"curveType"` // one letter
	CurvePoints [][2]int // delimiter| // delimiter: // slice of paired integers
	Slides      int      `json:"slides"`
	Length      float64  `json:"length"`
	EdgeSounds  []int    // delimiter|
	EdgeSets    [][2]int // delimiter| // delimiter:
}
type HitSample struct { // delimiter:
	NormalSet   int    `json:"normalSet"`   // nofloat
	AdditionSet int    `json:"additionSet"` // nofloat
	Index       int    `json:"index"`       // nofloat
	Volume      int    `json:"volume"`
	Filename    string `json:"filename"`
}

func newHitObject(line string) (HitObject, error) {
	var ho HitObject
	vs := strings.SplitN(line, `,`, 6)
	if len(vs) < 5 {
		return ho, fmt.Errorf("invalid hit object: %v (not enough length; requires 6)", vs)
	}
	x, err := strconv.ParseFloat(vs[0], 64)
	if err != nil {
		return ho, err
	}
	ho.X = int(x)

	y, err := strconv.ParseFloat(vs[1], 64)
	if err != nil {
		return ho, err
	}
	ho.Y = int(y)

	time, err := strconv.ParseFloat(vs[2], 64)
	if err != nil {
		return ho, err
	}
	ho.Time = int(time)

	noteType, err := strconv.Atoi(vs[3])
	if err != nil {
		return ho, err
	}
	ho.NoteType = noteType

	hitSound, err := strconv.Atoi(vs[4])
	if err != nil {
		return ho, err
	}
	ho.HitSound = hitSound
	if len(vs) == 5 {
		if ho.NoteType != HitTypeNote {
			return ho, fmt.Errorf("invalid hit object: %v; not enough length for non-normal note; requires 6)", vs)
		} else {
			return ho, nil
		}
	}
	var vs2 []string
	switch ho.NoteType & ComboMask {
	case HitTypeNote, HitTypeSlider, HitTypeSpinner:
		vs2 = strings.Split(vs[5], `,`)
	case HitTypeHoldNote:
		vs[5] = strings.Replace(vs[5], `,`, `:`, 1) // osu format v12 backward-compatibility
		vs2 = strings.SplitN(vs[5], `:`, 2)
	default:
		return ho, fmt.Errorf("error at %s's vs2: %s (not reach)", line, vs2)
	}

	var hsSkip bool
	switch ho.NoteType & ComboMask {
	case HitTypeNote:
		// hitSample

	case HitTypeSlider:
		// curveType|curvePoints,slides,length,edgeSounds,edgeSets,hitSample
		var vs3 string
		if strings.Contains(vs2[len(vs2)-1], `:`) {
			vs3 = strings.Join(vs2[:len(vs2)-1], `,`)
		} else { // no hit sample
			vs3 = strings.Join(vs2, `,`)
			hsSkip = true
		}
		sp, err := newSliderParams(vs3)
		if err != nil {
			return ho, err
		}
		ho.SliderParams = sp

	case HitTypeSpinner, HitTypeHoldNote:
		// endTime,hitSample
		// endTime:hitSample
		f, err := strconv.ParseFloat(vs2[0], 64)
		if err != nil {
			return ho, err
		}
		ho.EndTime = int(f)

	default:
		return ho, fmt.Errorf("invalid hit object: error at %s; invalid note type %d", line, ho.NoteType&ComboMask)
	}
	if hsSkip {
		return ho, nil
	}
	hs, err := newHitSample(vs2[len(vs2)-1])
	if err != nil {
		return ho, err
	}
	ho.HitSample = hs
	return ho, nil
}
func newSliderParams(s string) (SliderParams, error) {
	// curveType|curvePoints,slides,length,edgeSounds,edgeSets
	var sp SliderParams
	vs := strings.Split(s, `,`)
	if len(vs) < 3 {
		return sp, fmt.Errorf("invalid hit object: error at slider parameter %s; no enough length at %v", s, vs)
	}
	{
		// example: B|200:200|250:200
		vs2 := strings.Split(vs[0], `|`)
		sp.CurveType = vs2[0]
		sp.CurvePoints = make([][2]int, len(vs2)-1)
		for i, v := range vs2[1:] {
			var point [2]int
			for j, v2 := range strings.Split(v, `:`) {
				f, err := strconv.ParseFloat(v2, 64)
				if err != nil {
					return sp, err
				}
				point[j] = int(f)
			}
			sp.CurvePoints[i] = point
		}
	}
	{
		f, err := strconv.ParseFloat(vs[1], 64)
		if err != nil {
			return sp, err
		}
		sp.Slides = int(f)
	}
	{
		f, err := strconv.ParseFloat(vs[2], 64)
		if err != nil {
			return sp, err
		}
		sp.Length = f
	}
	if len(vs) == 3 {
		return sp, nil
	}
	if len(vs) < 5 {
		return sp, fmt.Errorf("invalid hit object: error at slider parameter %s; "+
			"no enough length at edge sound samples in slider parameter (%v)", s, vs)
	}
	{
		// example: 2|1|2
		vs2 := strings.Split(vs[3], `|`)
		sp.EdgeSounds = make([]int, len(vs2))
		for i := 0; i < len(vs2); i++ {
			f, err := strconv.ParseFloat(vs2[i], 64)
			if err != nil {
				return sp, err
			}
			sp.EdgeSounds[i] = int(f)
		}
	}
	{
		// example: 0:0|0:0|0:2
		vs2 := strings.Split(vs[4], `|`)
		sp.EdgeSets = make([][2]int, len(vs2))
		for i, v := range vs2 {
			var pair [2]int
			for j, v2 := range strings.Split(v, `:`) {
				f, err := strconv.ParseFloat(v2, 64)
				if err != nil {
					return sp, err
				}
				pair[j] = int(f)
			}
			sp.EdgeSets[i] = pair
		}
	}
	return sp, nil
}

//	if len(vs) < 5 {
//		return hs, errors.New("invalid hit sample: not enough length at hit sample")
//	}
func newHitSample(s string) (HitSample, error) {
	// normalSet:additionSet:index:volume:filename
	var hs HitSample
	vs := strings.Split(s, `:`)
	normalSet, err := strconv.Atoi(vs[0])
	if err != nil {
		return hs, err
	}
	hs.NormalSet = normalSet
	if len(vs) == 1 {
		return hs, nil // Todo: need to test whether this is not an error actually
	}

	additionSet, err := strconv.Atoi(vs[1])
	if err != nil {
		return hs, err
	}
	hs.AdditionSet = additionSet
	if len(vs) == 2 {
		return hs, nil // Todo: need to test whether this is not an error actually
	}

	index, err := strconv.Atoi(vs[2])
	if err != nil {
		return hs, err
	}
	hs.Index = index
	if len(vs) == 3 {
		return hs, nil // Todo: need to test whether this is not an error actually
	}

	volume, err := strconv.ParseFloat(vs[3], 64)
	if err != nil {
		return hs, err
	}
	hs.Volume = int(volume)
	if len(vs) == 4 {
		return hs, nil // Todo: need to test whether this is not an error actually
	}

	hs.Filename = vs[4]
	return hs, nil
}

// type HitType int

const (
	HitTypeNote = 1 << iota
	HitTypeSlider
	NewCombo
	HitTypeSpinner
	ComboColourSkip1
	ComboColourSkip2
	ComboColourSkip3
	HitTypeHoldNote
)
const ComboMask = ^(NewCombo + ComboColourSkip1 + ComboColourSkip2 + ComboColourSkip3)

// var ActualHitType = [...]HitType{HitTypeNote, HitTypeSlider, HitTypeSpinner, HitTypeHoldNote}

const (
	HitSoundNormal = 1 << iota
	HitSoundWhistle
	HitSoundFinish
	HitSoundClap
)

var HitSounds = []string{"normal", "whistle", "finish", "clap"}
var SampleSets = []string{"x", "normal", "soft", "drum"}
var YetDetermined = "?"

// Default sound sample would not even used by most users.
// func (h HitObject) SampleFilename() string {
// 	hs := h.HitSample
// 	set := hs.NormalSet
// 	if h.HitSound >= 2 {
// 		set = hs.AdditionSet
// 	}
// 	return fmt.Sprintf("%s-hit%s%d.wav", SampleSets[set], HitSounds[h.HitSound], hs.Index)
// }

const (
	TaikoKatMask = HitSoundWhistle | HitSoundClap
	TaikoBigMask = HitSoundFinish
)

// Column returns index of column at mania playfield
func (ho HitObject) Column(columnCount int) int {
	return columnCount * ho.X / 512
}
func IsKat(ho HitObject) bool { return ho.HitSound&TaikoKatMask != 0 }
func IsDon(ho HitObject) bool { return ho.HitSound&TaikoKatMask == 0 }
func IsBig(ho HitObject) bool { return ho.HitSound&TaikoBigMask != 0 }

func (h HitObject) SliderDuration(speed float64) int {
	if h.NoteType&HitTypeSlider == 0 {
		return 0
	}
	hs := h.SliderParams
	length := float64(hs.Slides) * hs.Length
	// speed := (bpm / 60000) * beatScale * (multiplier * 100)
	return int(length / speed)
}
