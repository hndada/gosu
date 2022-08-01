package osu

import (
	"errors"
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
	if len(vs) < 6 {
		return ho, errors.New("invalid hit object: not enough length")
	}
	{
		f, err := strconv.ParseFloat(vs[0], 64)
		if err != nil {
			return ho, err
		}
		ho.X = int(f)
	}
	{
		f, err := strconv.ParseFloat(vs[1], 64)
		if err != nil {
			return ho, err
		}
		ho.Y = int(f)
	}
	{
		f, err := strconv.ParseFloat(vs[2], 64)
		if err != nil {
			return ho, err
		}
		ho.Time = int(f)
	}
	{
		i, err := strconv.Atoi(vs[3])
		if err != nil {
			return ho, err
		}
		ho.NoteType = i
	}
	{
		i, err := strconv.Atoi(vs[4])
		if err != nil {
			return ho, err
		}
		ho.HitSound = i
	}

	var vs2 []string
	switch ho.NoteType & ComboMask {
	case HitTypeNote, HitTypeSlider, HitTypeSpinner:
		vs2 = strings.Split(vs[5], `,`)
	case HitTypeHoldNote:
		vs2 = strings.SplitN(vs[5], `:`, 2)
	default:
		return ho, errors.New("not reach")
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
		return ho, errors.New("invalid hit object data")
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
		return sp, errors.New("invalid hit object: not enough length at slider parameter")
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
		return sp, errors.New("invalid hit object: not enough length at edge sound samples in slider parameter ")
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

func newHitSample(s string) (HitSample, error) {
	// normalSet:additionSet:index:volume:filename
	var hs HitSample
	vs := strings.Split(s, `:`)
	if len(vs) < 5 {
		return hs, errors.New("invalid hit sample: not enough length")
	}
	{
		i, err := strconv.Atoi(vs[0])
		if err != nil {
			return hs, err
		}
		hs.NormalSet = i
	}
	{
		i, err := strconv.Atoi(vs[1])
		if err != nil {
			return hs, err
		}
		hs.AdditionSet = i
	}
	{
		i, err := strconv.Atoi(vs[2])
		if err != nil {
			return hs, err
		}
		hs.Index = i
	}
	{
		f, err := strconv.ParseFloat(vs[3], 64)
		if err != nil {
			return hs, err
		}
		hs.Volume = int(f)
	}
	{
		hs.Filename = vs[4]
	}
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
	HitSoundNormal = iota
	HitSoundWhistle
	HitSoundFinish
	HitSoundClap
)

// Todo: test
// Assumes whether normal or additional sample set is input in every call
func (hs HitSample) SampleFilename(sampleSet, hitSound int) string {
	if hs.Filename != "" {
		return hs.Filename
	}
	var sampleSetName, hitSoundName, index string
	switch sampleSet {
	case 1:
		sampleSetName = "normal"
	case 2:
		sampleSetName = "soft"
	case 3:
		sampleSetName = "drum"
	}
	switch hitSound {
	case HitSoundNormal:
		hitSoundName = "normal"
	case HitSoundWhistle:
		hitSoundName = "whistle"
	case HitSoundFinish:
		hitSoundName = "finish"
	case HitSoundClap:
		hitSoundName = "clap"
	}
	index = strconv.Itoa(hs.Index)
	return fmt.Sprintf("%s-hit%s%s.wav", sampleSetName, hitSoundName, index)
}

// Todo: test
func (ho HitObject) SliderDuration(tps []TimingPoint, multiplier float64) int {
	if ho.NoteType != HitTypeSlider {
		return 0
	}
	for j := len(tps) - 1; j >= 0; j-- {
		tp := tps[j]
		if tp.Time > ho.Time || tp.Uninherited {
			continue
		}
		duration := ho.SliderParams.Length / (multiplier * 100) * (-tp.BeatLength)
		return int(duration)
	}
	return 0
}

const (
	TaikoKatMask = HitSoundWhistle | HitSoundClap
	TaikoDonMask = ^TaikoKatMask // TODO: need a test
	TaikoBigMask = HitSoundFinish
)

// Column returns index of column at mania playfield
func (ho HitObject) Column(columnCount int) int {
	return columnCount * ho.X / 512
}
