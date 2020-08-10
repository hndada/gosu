package osu

import (
	"errors"
	"strconv"
	"strings"

	"github.com/hndada/gosu/internal/tools"
)

type HitObject struct {
	X         int
	Y         int
	StartTime int
	NoteType  int
	HitSound  int
	EndTime   int
	*SliderParams
	*HitSamples
}

const (
	NtNote = 1 << iota
	NtSlider
	NewCombo
	NtSpinner
	ColorSkip1
	ColorSkip2
	ColorSkip3
	NtHoldNote
	LastNoteType
)
const ComboMask = ^(NewCombo + ColorSkip1 + ColorSkip2 + ColorSkip3)

var NtArray = [...]int{NtNote, NtSlider, NtSpinner, NtHoldNote}

// todo: should I make custom type 'HitSound'?
const (
	Normal = iota
	Whistle
	Finish
	Clap
)

// todo: custom type of [2]int to Pair?
type SliderParams struct {
	CurveType   string   // one letter
	CurvePoints [][2]int // slice of paired integers
	Slides      int
	Length      float64
	EdgeSounds  [2]int
	EdgeSets    [2][2]int
}

type HitSamples struct {
	NormalSet   int
	AdditionSet int
	Index       int
	Volume      int
	Filename    string
}

// todo: should I make custom type 'SoundSet'?
const (
	Default = iota
	NormalSet
	SoftSet
	DrumSet
)

func parseHitObject(line string) (HitObject, error) {
	var hitObject HitObject
	var hitSamples HitSamples
	vs := strings.Split(line, `,`)

	// commons: x,y,time,type,hitSound
	x, err := tools.Atoi(vs[0])
	if err != nil {
		return hitObject, err
	}
	hitObject.X = x

	y, err := tools.Atoi(vs[1])
	if err != nil {
		return hitObject, err
	}
	hitObject.Y = y

	startTime, err := tools.Atoi(vs[2])
	if err != nil {
		return hitObject, err
	}
	hitObject.StartTime = startTime

	noteType, err := tools.Atoi(vs[3])
	if err != nil {
		return hitObject, err
	}
	noteType, err = parseNoteType(noteType)
	if err != nil {
		return hitObject, err
	}
	hitObject.NoteType = noteType

	hitSound, err := tools.Atoi(vs[4])
	if err != nil {
		return hitObject, err
	}
	hitObject.HitSound = hitSound

	switch hitObject.NoteType {
	case NtNote:
		// hitSample
		hitSamples, err = parseHitSamples(vs[5])
		if err != nil {
			return hitObject, err
		}
	case NtSlider:
		// curveType|curvePoints,slides,length,edgeSounds,edgeSets,hitSample
		var params SliderParams
		curveValues := strings.Split(vs[5], `|`)
		if !strings.ContainsAny(curveValues[0], "BCLP") && len(curveValues[0]) != 1 {
			return hitObject, errors.New("invalid curve type")
		}
		params.CurveType = curveValues[0]
		params.CurvePoints = make([][2]int, len(curveValues)-1)
		for i, s := range curveValues[1:] {
			p, err := tools.PairInt(s, `:`)
			if err != nil {
				return hitObject, err
			}
			params.CurvePoints[i] = p
		}

		slides, err := tools.Atoi(vs[6])
		if err != nil {
			return hitObject, err
		}
		params.Slides = slides

		length, err := strconv.ParseFloat(vs[7], 64)
		if err != nil {
			return hitObject, err
		}
		params.Length = length

		edgeSounds := strings.Split(vs[8], `|`)
		if len(edgeSounds) != 2 {
			return hitObject, errors.New("invalid edgeSounds")
		}
		for i := 0; i < 2; i++ {
			d, err := tools.Atoi(edgeSounds[i])
			if err != nil {
				return hitObject, err
			}
			params.EdgeSounds[i] = d
		}

		edgeSets := strings.Split(vs[9], `|`)
		if len(edgeSets) != 2 {
			return hitObject, errors.New("invalid edgeSets")
		}
		for i := 0; i < 2; i++ {
			p, err := tools.PairInt(edgeSets[i], `:`)
			if err != nil {
				return hitObject, err
			}
			params.EdgeSets[i] = p
		}

		hitObject.SliderParams = &params

		hitSamples, err = parseHitSamples(vs[10])
		if err != nil {
			return hitObject, err
		}
	case NtSpinner:
		// endTime,hitSample
		endTime, err := tools.Atoi(vs[5])
		if err != nil {
			return hitObject, err
		}
		hitObject.EndTime = endTime

		hitSamples, err = parseHitSamples(vs[6])
		if err != nil {
			return hitObject, err
		}
	case NtHoldNote:
		// endTime:hitSample
		vs2 := strings.SplitN(vs[5], `:`, 2)

		endTime, err := tools.Atoi(vs2[0])
		if err != nil {
			return hitObject, err
		}
		hitObject.EndTime = endTime

		hitSamples, err = parseHitSamples(vs2[1])
		if err != nil {
			return hitObject, err
		}
	default:
		return hitObject, errors.New("cannot reach")
	}

	hitObject.HitSamples = &hitSamples
	return hitObject, nil
}

func parseNoteType(v int) (int, error) {
	nt := v & ComboMask
	for _, v := range NtArray {
		if nt == v {
			return nt, nil
		}
	}
	return -1, errors.New("invalid note type")
}

func parseHitSamples(s string) (HitSamples, error) {
	// normalSet:additionSet:index:volume:filename
	var hitSamples HitSamples
	vs := strings.Split(s, `:`)

	normalSet, err := tools.Atoi(vs[0])
	if err != nil {
		return hitSamples, err
	}
	hitSamples.NormalSet = normalSet

	additionSet, err := tools.Atoi(vs[1])
	if err != nil {
		return hitSamples, err
	}
	hitSamples.AdditionSet = additionSet

	index, err := tools.Atoi(vs[2])
	if err != nil {
		return hitSamples, err
	}
	hitSamples.Index = index

	volume, err := tools.Atoi(vs[3])
	if err != nil {
		return hitSamples, err
	}
	hitSamples.Volume = volume

	hitSamples.Filename = vs[4]

	return hitSamples, nil
}
