package osu

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

type TimingPoint struct { // delimiter,
	Time        int
	BeatLength  float64
	Meter       int
	SampleSet   int
	SampleIndex int
	Volume      int
	Uninherited bool
	Effects     int
}

// issue: The code `if err != nil { return }` cannot be omitted
// in current Go version. There is a proposal of new error handling, though.
func newTimingPoint(line string) (tp TimingPoint, err error) {
	// time,beatLength,meter,sampleSet,sampleIndex,volume,uninherited,effects
	vs := strings.Split(line, `,`)
	if len(vs) < 8 {
		return tp, errors.New("invalid timing point: not enough length")
	}

	if tp.Time, err = parseInt(vs[0]); err != nil {
		return tp, fmt.Errorf("%s: %w", line, err)
	}
	if tp.BeatLength, err = parseFloat(vs[1]); err != nil {
		switch vs[1] {
		case "∞":
			tp.BeatLength = math.Inf(1)
		case "-∞":
			tp.BeatLength = math.Inf(-1)
		default:
			return tp, fmt.Errorf("%s: %w", line, err)
		}
	}
	if tp.Meter, err = parseInt(vs[2]); err != nil {
		return tp, fmt.Errorf("%s: %w", line, err)
	}
	if tp.SampleSet, err = parseInt(vs[3]); err != nil {
		return tp, fmt.Errorf("%s: %w", line, err)
	}
	if tp.SampleIndex, err = parseInt(vs[4]); err != nil {
		return tp, fmt.Errorf("%s: %w", line, err)
	}
	if tp.Volume, err = parseInt(vs[5]); err != nil {
		return tp, fmt.Errorf("%s: %w", line, err)
	}
	if tp.Uninherited, err = parseBool(vs[6]); err != nil {
		return tp, fmt.Errorf("%s: %w", line, err)
	}
	if tp.Effects, err = parseInt(vs[7]); err != nil {
		return tp, fmt.Errorf("%s: %w", line, err)
	}
	return tp, nil
}

// Inherited means some of its values are derived from the base timing point.
// The base timing point is called Uninherited.
// Todo: parent timing point or previous timing point?
func (tp TimingPoint) IsInherited() bool { return !tp.Uninherited }

// BPM supposes the given timing point is Uninherited.
func (tp TimingPoint) BPM() float64 { return 1000 * 60 / tp.BeatLength }

// BeatLengthScale returns a beat scale, aka speed factor. The standard value is 1.
// BeatLengthScale supposes the given timing point is Inherited.
func (tp TimingPoint) BeatLengthScale() float64 { return 100 / (-tp.BeatLength) }

// Kiai infers to highlight.
func (tp TimingPoint) IsKiai() bool            { return tp.Effects&1 != 0 }
func (tp TimingPoint) IsFirstBarOmitted() bool { return tp.Effects&(1<<3) != 0 }
