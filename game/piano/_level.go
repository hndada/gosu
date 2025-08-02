package piano

import (
	"sort"

	"github.com/hndada/gosu/plays"
)

type Level struct {
	Level float64
	steps []step
}

// Threshold that determines whether a note is in a step or not.
const inStepThreshold = 30

// Strain is for calculating difficulty.
var chordStrain = func(len float64) float64 { return 1.0/len + 0.1*(len-1) }
var jackStrain = plays.LinearInterpolate([]float64{0, 200}, []float64{1.5, 0})
var bombStrain = plays.LinearInterpolate([]float64{0, 200}, []float64{0.75, 0})

const tailHoldBonus = 1.40
const normalHoldBonus = 1.25

// weight is for Tail's variadic weight based on its length.
// For example, short long note does not require much strain to release.
var tailStrain = plays.LinearInterpolate(
	[]float64{0, 50, 200, 800}, []float64{0.4, 0.1, 0.1, 0.7})

type step struct {
	time  int32 // base time: average or min of notes' time
	notes []*Note

	hands        []int     // from notes; used at chordStrains
	holdings     []bool    // from staged notes; used at chordStrains
	chordStrains []float64 // from current step.notes
	jackStrains  []float64 // from each note.Prev
	bombStrains  []float64 // from staged notes
	weights      []float64 // from each note.Weight()
}

func (c *Chart) setSteps() {
	var st = step{
		time:  c.Notes[0].Time,
		notes: make([]*Note, c.KeyCount()),
	}
	// It is guaranteed that n is in stagedListNotes since it is sorted by time.
	stagedList := c.newStagedNotes()

	for _, n := range c.Notes {
		// Start with new step if the note is too far or the lane has occupied.
		if n.Time-st.time > inStepThreshold || st.notes[n.Key] != nil {
			// calculate strains of the step
			st.setHands()
			st.setHoldings(stagedList)
			st.setChordStrains()
			st.setJackStrains()
			st.setBombStrains(stagedList)
			st.setWeights()

			// append
			c.steps = append(c.steps, st)
			st = step{
				time:  n.Time,
				notes: make([]*Note, c.KeyCount()),
			}
		}
		// set to the step
		st.notes[n.Key] = n
		stagedList[n.Key] = n.Next
	}
	c.steps = append(c.steps, st)
}

const (
	none = iota
	leftHand
	rightHand
)

// Hand of the middle note is trivial in even keys: right hand.
// In odd keys, the middle note is assigned to the hand with more notes.
// Todo: handle scratch
func (st *step) setHands() {
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

func (st *step) setHoldings(stagedList []*Note) {
	st.holdings = make([]bool, len(st.notes))
	for k, sn := range stagedList {
		if sn == nil {
			continue
		}
		if sn.Type != Tail {
			continue
		}
		if n := st.notes[k]; n != nil {
			continue
		}
		// Check the remaining duration of the long note is long enough.
		st.holdings[k] = sn.Time-st.time > inStepThreshold
	}
}

// Condition of chord:
// 1. same hand
// 2. Tail / non-Tail (Tail is for release)
// 3. long notes cut the chord
type chordKey struct {
	hand         int
	tail         bool
	holdingIndex int
}

func (st *step) setChordStrains() {
	chords := make(map[chordKey][]*Note)
	var holdingIndex = 0
	for k, n := range st.notes {
		if st.holdings[k] {
			holdingIndex++
		}
		if n == nil {
			continue
		}
		ck := chordKey{
			hand:         st.hands[k],
			tail:         n.Type == Tail,
			holdingIndex: holdingIndex,
		}
		chords[ck] = append(chords[ck], n)
	}

	st.chordStrains = make([]float64, len(st.notes))
	for _, ns := range chords {
		strain := chordStrain(float64(len(ns)))
		for _, n := range ns {
			st.chordStrains[n.Key] = strain
		}
	}
}

func (st *step) setJackStrains() {
	st.jackStrains = make([]float64, len(st.notes))
	for k, n := range st.notes {
		if n == nil {
			continue
		}
		if n.Prev == nil {
			continue
		}
		// Long note itself has no jack strain.
		if n.Type == Tail {
			continue
		}
		x := float64(n.Time - n.Prev.Time)
		st.jackStrains[k] = jackStrain(x)
	}
}

// Bomb a virtual note that is not in a step, but should not be pressed.
// If a bomb is pressed, it will judge the staged note with poor judgment.
func (st *step) setBombStrains(stagedList []*Note) {
	st.bombStrains = make([]float64, len(st.notes))
	for k, sn := range stagedList {
		n := st.notes[k]
		if sn == nil {
			continue
		}
		if n != nil {
			continue
		}
		if sn.Type == Tail {
			continue
		}
		x := float64(sn.Time - st.time)
		st.bombStrains[k] = bombStrain(x)
	}
}

func (st *step) setWeights() {
	// Holding bonus is enabled if the hand is holding at least one long note.
	var (
		leftHandHoldBonus  bool
		rightHandHoldBonus bool
	)
	for k, h := range st.holdings {
		if h {
			switch st.hands[k] {
			case leftHand:
				leftHandHoldBonus = true
			case rightHand:
				rightHandHoldBonus = true
			}
		}
	}

	st.weights = make([]float64, len(st.notes))
	for k, n := range st.notes {
		if n == nil {
			continue
		}

		switch n.Type {
		case Tail:
			// Todo: put Duration on Tail
			head := n.Prev
			x := float64(head.Duration)
			st.weights[k] = tailStrain(x)
		default:
			st.weights[k] = 1
		}

		if st.hands[k] == leftHand && leftHandHoldBonus ||
			st.hands[k] == rightHand && rightHandHoldBonus {
			if n.Type == Tail {
				st.weights[k] *= tailHoldBonus
			} else {
				st.weights[k] *= normalHoldBonus
			}
		}
	}
}

func (st *step) strain() float64 {
	var strain float64
	for k, w := range st.weights {
		base := st.chordStrains[k] + st.jackStrains[k] + st.bombStrains[k]
		strain += base * w
	}
	return strain
}

const unitDuration = 800 // 800ms. 2 beats with 150 BPM
const decayFactor = 0.95
const levelScale = 0.05

// Different BPM make duration of diff different.
// However, it looks fine not to scale each diff based on its duration
// and using the same size of duration on each piece.
// They will be alleviated into diffs.
func (c *Chart) setLevel() {
	// times, durations := plays.DifficultyPieceTimes(c.Dynamics, c.Duration())
	// scale := standardDuration / float64(durations[i])
	// diffs = append(diffs, diff*scale)
	diffs := make([]float64, 0, len(c.steps))

	// chart has at least one step.
	endTime := c.steps[0].time + unitDuration
	var diff float64
	for _, st := range c.steps {
		if st.time > endTime {
			diffs = append(diffs, diff)
			diff = 0
			endTime += unitDuration
		}
		diff += st.strain()
	}

	sort.Slice(diffs, func(i, j int) bool { return diffs[i] > diffs[j] })
	difficulty := plays.WeightedSum(diffs, decayFactor)

	// No additional Math.Pow; it would make a little change.
	c.Level.Level = difficulty * levelScale
}

// Todo: debug level calculation
// color each note based on its strain
// with printing strain value.

// type Pattern []step // for readability
