package osu

import (
	"fmt"
	"strings"
)

const (
	HitTypeNote    = 1 << iota
	HitTypeSlider  // 1 << 1 = 2
	NewCombo       // Appears as 5 = 1 + 4 at first note.
	HitTypeSpinner // 1 << 3 = 8

	ComboColourSkip1
	ComboColourSkip2
	ComboColourSkip3
	HitTypeHoldNote // 1 << 7 = 128
)

type HitObject struct { // delimiter,
	X            int
	Y            int
	Time         int
	NoteType     int
	HitSound     int
	EndTime      int
	SliderParams SliderParams
	HitSample    HitSample
}

const comboMask = ^(NewCombo + ComboColourSkip1 + ComboColourSkip2 + ComboColourSkip3)

func newHitObject(line string) (ho HitObject, err error) {
	// x,y,time,type,hitSound,objectParams,hitSample
	vs := strings.SplitN(line, `,`, 6)
	if len(vs) < 5 {
		return ho, fmt.Errorf("hit object has not enough length: %v", vs)
	}

	if ho.X, err = parseInt(vs[0]); err != nil {
		return
	}
	if ho.Y, err = parseInt(vs[1]); err != nil {
		return
	}
	if ho.Time, err = parseInt(vs[2]); err != nil {
		return
	}
	if ho.NoteType, err = parseInt(vs[3]); err != nil {
		return
	}
	if ho.HitSound, err = parseInt(vs[4]); err != nil {
		return
	}

	if len(vs) == 5 {
		// x,y,time,type,hitSound
		if ho.NoteType == HitTypeNote {
			return
		}
		return ho, fmt.Errorf("hit object has not enough length: %v", vs)
	}

	// The first comment of each case is the format of sub.
	sub := vs[5] // A remained substring.
	var hitSampleData string
	switch ho.NoteType & comboMask {
	case HitTypeNote:
		// hitSample
		hitSampleData = sub

	case HitTypeSlider:
		// curveType|curvePoints,slides,length,edgeSounds,edgeSets,hitSample

		// It is fine to use same variable name in different scopes
		// even if they are in same function.
		// https://go.dev/play/p/EqVlYivsxOE

		vs := strings.Split(sub, `,`)
		if strings.Contains(vs[len(vs)-1], `:`) {
			hitSampleData = vs[len(vs)-1]
			vs = vs[:len(vs)-1] // drop the last element
		}

		objectParamsData := strings.Join(vs, `,`)
		if ho.SliderParams, err = newSliderParams(objectParamsData); err != nil {
			return
		}

	case HitTypeSpinner:
		// endTime,hitSample
		vs := strings.SplitN(sub, `,`, 2)
		if ho.EndTime, err = parseInt(vs[0]); err != nil {
			return
		}
		hitSampleData = vs[1]

	case HitTypeHoldNote:
		// endTime:hitSample
		vs := strings.SplitN(sub, `:`, 2)
		if ho.EndTime, err = parseInt(vs[0]); err != nil {
			return
		}
		hitSampleData = vs[1]

	default:
		return ho, fmt.Errorf("hit object has invalid note type: %v", vs)
	}

	ho.HitSample, err = newHitSample(hitSampleData)
	return
}

// Column returns index of column at osu!mania playfield.
func (ho HitObject) Column(columnCount int) int { return columnCount * ho.X / 512 }

const (
	taikoKatMask = HitSoundWhistle | HitSoundClap
	taikoBigMask = HitSoundFinish
)

func (ho HitObject) IsDon() bool { return ho.HitSound&taikoKatMask == 0 }
func (ho HitObject) IsKat() bool { return ho.HitSound&taikoKatMask != 0 }
func (ho HitObject) IsBig() bool { return ho.HitSound&taikoBigMask != 0 }

// A strange code written in parsing HitTypeHoldNote.
// vs[5] = strings.Replace(vs[5], `,`, `:`, 1) // osu format v12 backward-compatibility
