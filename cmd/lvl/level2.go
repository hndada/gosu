package main

import "github.com/hndada/gosu/game/piano"

var ChordStrain = func(n int) float64 { return 1.0/float64(n) + 0.1*float64(n-1) }
var JackWeight = SimpleDecay(1.0, 200, 1.25)
var TailStrain = LinearFunc(
	[]float64{0, 50, 200, 800},
	[]float64{0.4, 0.1, 0.1, 0.7},
)

// Bonus: plus(+) operator
// Weight: times(x) operator
// Factor: times(x) operator with constant operand
const (
	LNReleasingBonusFactor = 1.40 // Releasing tail
	LNHoldingBonusFactor   = 1.25 // Holding LN
)

type Note = piano.Note

// Step is convenient for handling "Bomb" note
type Step struct {
	time  int32
	notes []Note

	hands              []int
	baseStrains        []float64 // based on finger factors
	chordFactors       []float64
	isSameHandHoldings []bool
	// jack
	// bomb

	next *Step
	prev *Step
}

func (st Step) strain() float64 {
	var strain float64
	for k, w := range st.weights {
		base := st.chordStrains[k] + st.jackStrains[k] + st.bombStrains[k]
		strain += base * w
	}
	return strain
}

const maxStepTimeWindow = 30

func (st Step) accepts(target Note) bool {
	if len(st.notes) < target.Key {
		panic("Step.accepts(): not enough notes length")
	}

	if !st.notes[target.Key].IsBlank() {
		return false
	}

	first := st.notes[0]
	return target.Time-first.Time <= maxStepTimeWindow
}

// Required preprocess:
// 1. Sort by time and key
// 2. Set prev and next note
func NewSteps(ns []Note, keyCount int) []Step {
	if len(ns) == 0 {
		return nil
	}

	stagedNotes := make([]Note, keyCount)
	for k := 0; k < keyCount; k++ {
		for _, n := range ns {
			if n.Key != k {
				continue
			}
			stagedNotes[k] = n
			break
		}
	}

	steps := make([]Step, 0, len(ns))

	// The first step.
	st := Step{
		time:  ns[0].Time,
		notes: make([]Note, keyCount),
	}

	for _, n := range ns {
		if !st.accepts(n) {
			// TODO: calculate strains

			steps = append(steps, st)
			st = Step{
				time:  n.Time,
				notes: make([]Note, keyCount),
			}
		}
		st.notes[n.Key] = n

	}
	return steps
}

const (
	none = iota
	leftHand
	rightHand
)

// Hand of the middle note is trivial in even keys: right hand.
// In odd keys, the middle note is assigned to the hand with more notes.
// Todo: handle scratch
func (st Step) setHands() {
	st.hands = make([]int, len(st.notes))
	leftCount, rightCount := 0, 0

	for k, n := range st.notes {
		if n == nil {
			continue
		}
		switch {
		case k < len(st.notes)/2:
			st.hands[k] = leftHand
			leftCount++
		case k > len(st.notes)/2:
			st.hands[k] = rightHand
			rightCount++
		}
	}

	middle := len(st.notes) / 2
	if st.notes[middle] == nil {
		return
	}

	if len(st.notes)%2 == 0 {
		st.hands[middle] = rightHand
	} else {
		if leftCount < rightCount {
			st.hands[middle] = rightHand
		} else {
			st.hands[middle] = leftHand
		}
	}
}
