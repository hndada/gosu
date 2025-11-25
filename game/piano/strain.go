package piano

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/mat"
)

const maxStepTimeWindow = 30

var fingerIndicesMap = make(map[int][]int)
var leftKeyIndicesMap = make(map[int][5]int)
var rightKeyIndicesMap = make(map[int][5]int)

func init() {
	initIndices()
}
func initIndices() {
	// Left-scratch is default in 8 key mode
	isLeftScratchModes := map[int]bool{8: true}

	// 1. fingerIndicesMap
	// 7:  {3, 2, 1, 0, 1, 2, 3},
	// 8:  {4, 3, 2, 1, 0, 1, 2, 3}, // Left-scratch
	// 9:  {4, 3, 2, 1, 0, 1, 2, 3, 4},
	for keyCount := 1; keyCount <= 10; keyCount++ {
		if isLeftScratchModes[keyCount] {
			fingerIndicesMap[keyCount] = append([]int{4},
				fingerIndicesMap[keyCount-1]...)
			continue
		}

		fins := make([]int, keyCount)
		mid := keyCount / 2
		for k := 0; k < keyCount; k++ {
			if k < mid {
				fins[k] = mid - k
			} else {
				if keyCount%2 == 0 {
					fins[k] = k - mid + 1 // skip thumb
				} else {
					fins[k] = k - mid
				}
			}
		}
		fingerIndicesMap[keyCount] = fins
	}

	// 2. keyIndicesMap
	//	6: [2][5]int{
	//		{-1, 2, 1, 0, -1},
	//		{-1, 3, 4, 5, -1},
	//	},
	//	7: [2][5]int{
	//		[5]int{3, 2, 1, 0, -1},
	//		[5]int{3, 4, 5, 6, -1},
	//	},
	for keyCount := 1; keyCount <= 9; keyCount++ {
		var lkis, rkis [5]int
		if isLeftScratchModes[keyCount] {
			prevKis := leftKeyIndicesMap[keyCount-1]
			for f := range prevKis {
				if f == 0 {
					lkis[f] = 4
					continue
				}
				lkis[f] = prevKis[f-1]
			}
			rkis = rightKeyIndicesMap[keyCount-1]
			leftKeyIndicesMap[keyCount] = lkis
			rightKeyIndicesMap[keyCount] = rkis
			continue
		}

		for i := 0; i < 5; i++ {
			lkis[i] = -1
			rkis[i] = -1
		}
		mid := keyCount / 2
		for k, fi := range fingerIndicesMap[keyCount] {
			if k < mid {
				lkis[fi] = k
			} else {
				rkis[fi] = k
			}
		}
		// map the thumb on odd key mode
		if keyCount%2 == 1 {
			lkis[0] = rkis[0]
		}
		leftKeyIndicesMap[keyCount] = lkis
		rightKeyIndicesMap[keyCount] = rkis
	}
}

const (
	leftHand = iota
	rightHand
	middle // To be determined
)
const mainHand = rightHand
const subHand = (mainHand + 1) % 2

// Hand of the middle note is trivial in even keys: right hand.
// In odd keys, the middle note is assigned to the hand which has
// closer note on its side.
func (c *Chart) setHands() {
	hands := make([]int, len(c.notes))

	mid := c.keyCount / 2
	for i, n := range c.notes {
		switch {
		case n.Key < mid:
			hands[i] = leftHand
		case c.keyCount%2 != 0 && n.Key == mid:
			hands[i] = middle
		default:
			hands[i] = rightHand
		}
	}

	// Determine 'middle' hand
	for i, h := range hands {
		if h != middle {
			continue
		}

		if i == 0 || i == len(hands)-1 {
			hands[i] = mainHand
			continue
		}

		prevHand := hands[i-1]
		nextHand := hands[i+1]
		if prevHand != subHand || nextHand != subHand {
			hands[i] = mainHand
			continue
		}

		prevNote := c.notes[i-1]
		currNote := c.notes[i]
		nextNote := c.notes[i+1]
		pdt := currNote.Time - prevNote.Time // prev delta time
		ndt := nextNote.Time - currNote.Time
		if pdt == ndt {
			hands[i] = mainHand
		} else if pdt < ndt {
			hands[i] = leftHand
		} else {
			hands[i] = rightHand
		}
	}

	for i, h := range hands {
		c.notes[i].hand = h
	}
}

// Required preprocess:
// 1. Sort by time and key
// 2. Set prev and next note
func (c *Chart) setStepIDs() {
	c.setHands()
	pn := c.notes[0] // pivot note
	exists := make([]bool, c.keyCount)
	exists[pn.Key] = true
	for i, n := range c.notes {
		if i == 0 {
			c.notes[0].step = 0
			continue
		}
		prev := c.notes[i-1]

		isExist := exists[n.Key]
		sameHand := pn.hand == n.hand
		inTime := n.Time-pn.Time <= maxStepTimeWindow
		if !isExist && sameHand && inTime {
			c.notes[i].step = prev.step
		} else {
			pn = c.notes[i]
			exists = make([]bool, c.keyCount)
			c.notes[i].step = prev.step + 1
		}
		exists[n.Key] = true
	}
}

func (c *Chart) calcStrains() {
	c.setStepIDs()

	var (
		fis = fingerIndicesMap[c.keyCount]
		// lrkis: left right key indices
		lrkis = [2][5]int{
			leftKeyIndicesMap[c.keyCount],
			rightKeyIndicesMap[c.keyCount],
		}
		// lrhs: left right hand states
		lrhs = [2]handState{
			{handKind: leftHand,
				fis:       fis,
				kis:       lrkis[leftHand],
				prevNotes: make([]*Note, c.keyCount)},
			{handKind: rightHand,
				fis:       fis,
				kis:       lrkis[rightHand],
				prevNotes: make([]*Note, c.keyCount)},
		}
		sns     = make([]*Note, 0, c.keyCount) // step notes
		strains [5]float64
	)

	for i, n := range c.notes {
		// Check if prev step is complete
		if i != 0 && n.step != sns[0].step {
			hand := sns[0].hand
			hs := &lrhs[hand]
			strains = hs.calcStrain(sns)
			for _, sn := range sns {
				fin := hs.fis[sn.Key]
				sn.strain = strains[fin]
			}
			sns = make([]*Note, 0, 5)
		}
		// Beware not to use a pointer to local variable 'n'
		sns = append(sns, &c.notes[i])
	}
	minStrain := 10000.0
	maxStrain := -1000.0
	for _, n := range c.notes {
		if minStrain > n.strain {
			minStrain = n.strain
		} else if maxStrain < n.strain {
			maxStrain = n.strain
		}
	}
	fmt.Printf("min: %.2f	max: %.2f\n", minStrain, maxStrain)
}

type handState struct {
	handKind int
	fis      []int
	kis      [5]int

	time      int32
	positions [5]float64
	prevNotes []*Note
}

// It is guaranteed that step notes contains non-nil elements only.
func stepTime(step []*Note) int32 {
	var sumTime int32
	for _, n := range step {
		sumTime += n.Time
	}
	return sumTime / int32(len(step))
}

// To press the note, it should starts at the top (0)
// and keep pressing toward the bottom (1). Hence, we will
// correct the position first before press the notes.
// Meanwhile, no need to correct the pos for releasing as
// they are always at the bottom.
func (hs *handState) calcStrain(step []*Note) [5]float64 {
	const normalKeyStrokeDuration = 100 // ms
	// Adjust positions for non-holding keystroke.
	// Normal keystroke takes 80~120ms.
	st := stepTime(step)
	dt := st - hs.time // delta time
	for fin, pos := range hs.positions {
		k := hs.kis[fin]
		if k == -1 {
			continue
		}
		pn := hs.prevNotes[k]
		if pn == nil || pn.Kind == Head {
			continue // should be 1.0
		}
		pos2 := pos - float64(dt)/normalKeyStrokeDuration
		if pos2 < 0 {
			pos2 = 0
		}
		hs.positions[fin] = pos2
	}

	var (
		reqMoves1 [5]float64
		reqMoves2 [5]float64
		strains   [5]float64
	)
	// 1. Release all non-holding keys for preparing pressing.
	// 2. Press the keys
	for _, sn := range step {
		fin := hs.fis[sn.Key]
		switch sn.Kind {
		case Normal, Head:
			reqMoves1[fin] = 0.0
			reqMoves2[fin] = 1.0
		case Tail:
			reqMoves1[fin] = 1.0
			reqMoves2[fin] = 0.0
		}
	}

	for _, reqMoves := range [2][5]float64{reqMoves1, reqMoves2} {
		partialStrains := hs.calcPartialStrains(reqMoves)
		for i, s := range partialStrains {
			strains[i] += s
		}
		respMoves := hs.calcMoves(partialStrains)
		for i, m := range respMoves {
			hs.positions[i] += m
		}
	}

	// Store for next calculation
	hs.time = st
	for _, sn := range step {
		hs.prevNotes[sn.Key] = sn
	}
	return strains
}

// func (h handState) calcReqMoves(poses [5]float64, ps [5]bool) [5]float64 {
// 	var reqMoves [5]float64
// 	for i, pos := range poses {
// 		if ps[i] { // pressed
// 			reqMoves[i] = 1.0 - pos
// 		} else {
// 			reqMoves[i] = 0.0 - pos
// 		}
// 	}
// 	return reqMoves
// }

// How the influence matrix is derived:
// 1. Degree of dependence among fingers.
// Neighbor finger is affected strongly, vice versa.
// Index is the most independent, pinky is the opposite.
// Matrix is symmetric before scaled.
// var nonNeighborInfluence = 0.05
// var neighborInfluences = [4]float64{0.05, 0.1, 0.2, 0.3}
//
// 2. base strain: Index ≈ Middle > Ring > Pinky
// Each row is inversely scaled by the following values:
// var baseStrains [5]float64{1.1, 1.0, 1.05, 1.1, 1.2}
var influences = [25]float64{
	0.90909, 0.04545, 0.04545, 0.04545, 0.04545,
	0.05000, 1.00000, 0.10000, 0.05000, 0.05000,
	0.04762, 0.09524, 0.95238, 0.19048, 0.04762,
	0.04545, 0.04545, 0.18182, 0.90909, 0.27273,
	0.04167, 0.04167, 0.04167, 0.25000, 0.83333,
}

// calcStrains calculates strain with the following process.
// 1. Extract active variables only
// 2. Solve linear equations
// 3. Convert solutions into [5]float64
func (h handState) calcPartialStrains(dps [5]float64) [5]float64 {
	var strains [5]float64

	activeFingers := make([]int, 0, 5)
	dps2 := make([]float64, 0, 5)
	for i, dp := range dps {
		if math.Abs(dp) > 1e-9 { // check non-zero
			activeFingers = append(activeFingers, i)
			dps2 = append(dps2, dp)
		}
	}

	rank := len(activeFingers)
	infl := make([]float64, rank*rank)
	for i, fin1 := range activeFingers {
		for j, fin2 := range activeFingers {
			infl[i*rank+j] = influences[5*fin1+fin2]
		}
	}
	if rank == 0 {
		return strains
	}

	x := new(mat.VecDense)
	a := mat.NewDense(rank, rank, infl)
	b := mat.NewVecDense(rank, dps2)
	x.SolveVec(a, b) // solves a*x = b
	xdata := x.RawVector().Data

	for i, fin := range activeFingers {
		strains[fin] = math.Abs(xdata[i])
	}
	return strains
}

func (h handState) calcMoves(strains [5]float64) [5]float64 {
	var moves [5]float64
	for i := range strains {
		move := 0.0
		for j := 0; j < 5; j++ {
			move += influences[i*5+j] * strains[j]
		}
		moves[i] = move
	}
	return moves
}
