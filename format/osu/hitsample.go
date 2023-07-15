package osu

import (
	"strconv"
	"strings"
)

const (
	HitSoundNormal = 1 << iota
	HitSoundWhistle
	HitSoundFinish
	HitSoundClap
)

type HitSample struct { // delimiter:
	NormalSet   int    // Sample set of the normal sound.
	AdditionSet int    // Sample set of the whistle, finish, and clap sounds.
	Index       int    // If this is 0, the timing point's sample index will be used.
	Volume      int    // If this is 0, the timing point's volume will be used.
	Filename    string // Custom filename of the addition sound.
}

// There are valid hit samples which length is fewer than 5.
// I suppose zero values are used for the rest of the fields.
// Reference: 193127 fripSide - only my railgun (TV Size) (DJPop) [4K HD].osu
// v11: 192,320,8391,128,2,9860:1:0:2
func newHitSample(s string) (hs HitSample, err error) {
	// normalSet:additionSet:index:volume:filename
	vs := strings.Split(s, `:`)

	if len(vs) >= 1 {
		if hs.NormalSet, err = parseInt(vs[0]); err != nil {
			return
		}
	}
	if len(vs) >= 2 {
		if hs.AdditionSet, err = parseInt(vs[1]); err != nil {
			return
		}
	}
	if len(vs) >= 3 {
		if hs.Index, err = parseInt(vs[2]); err != nil {
			return
		}
	}
	if len(vs) >= 4 {
		if hs.Volume, err = parseInt(vs[3]); err != nil {
			return
		}
	}
	hs.Filename = vs[4]

	return
}

func (h HitObject) SampleFilename() string {
	if h.HitSample.Filename != "" {
		return h.HitSample.Filename
	}

	// default sample
	// Todo: test
	sampleSetName := h.sampleSetName()
	if sampleSetName == "" {
		return ""
	}

	// <sampleSet>-hit<hitSound><index>.wav
	var b strings.Builder
	b.WriteString(sampleSetName)
	b.WriteString("-hit")
	b.WriteString(h.hitSoundName())

	// Index is omitted if it is 0 or 1.
	if h.HitSample.Index >= 2 {
		b.WriteString(strconv.Itoa(h.HitSample.Index))
	}

	// Seems in old maps, .ogg are also supported.
	// Reference: 63089 fripSide - only my railgun (TV Size)
	b.WriteString(".wav")
	return b.String()
}

func (h HitObject) isHitSoundNormal() bool { return h.HitSound <= HitSoundNormal }

var sampleSetNames = [...]string{"", "normal", "soft", "drum"}

func (h HitObject) sampleSetName() string {
	i := h.HitSample.NormalSet
	if !h.isHitSoundNormal() && h.HitSample.AdditionSet != 0 {
		i = h.HitSample.AdditionSet
	}
	return sampleSetNames[i]
}

var hitSoundNames = [...]string{"normal", "whistle", "finish", "clap"}

// Todo: need a test when both bits are set such as Finish and Clap.
func (h HitObject) hitSoundName() string {
	for i := 3; i >= 0; i-- {
		v := 1 << i // 8, 4, 2, 1
		if h.HitSound&v != 0 {
			return hitSoundNames[i]
		}
	}
	return hitSoundNames[0]
}
