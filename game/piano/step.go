package piano

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

const maxStepTimeWindow = 30

// Required preprocess:
// 1. Sort by time and key
// 2. Set prev and next note
func (c Chart) setStepIDs() {
	c.setHands()
	pn := c.notes[0] // pivot note
	exists := make([]bool, c.keyCount)
	exists[pn.Key] = true
	for i, n := range c.notes {
		if i == 0 {
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

var FingersMap = map[int][]int{
	1:  {0},
	2:  {1, 1},
	3:  {1, 0, 1},
	4:  {2, 1, 1, 2},
	5:  {2, 1, 0, 1, 2},
	6:  {3, 2, 1, 1, 2, 3},
	7:  {3, 2, 1, 0, 1, 2, 3},
	8:  {4, 3, 2, 1, 0, 1, 2, 3}, // Left-scratch
	9:  {4, 3, 2, 1, 0, 1, 2, 3, 4},
	10: {4, 3, 2, 1, 0, 0, 1, 2, 3, 4},
}

func (c *Chart) calcStrains() {
	var left, right Hand
	fingers := FingersMap[c.keyCount]
	stepNotes := make([]Note, 0, 5)
	for i, n := range c.notes {
		// Check if prev step is complete
		if i != 0 && n.step != stepNotes[0].step {
			var pressings [5]bool
			for _, n := range stepNotes {
				fi := fingers[n.Key]
				pressings[fi] = true
			}

			hand := stepNotes[0].hand
			if hand == leftHand {
				strains := left.calcStrain()
			} else {
				strains := right.calcStrain()
			}
			stepNotes = make([]Note, 0, 5)
		}
		stepNotes = append(stepNotes, n)
	}
}

// var FromHandFingerToKeyMap = map[int][2][]int{
// 	7: [2]int{
// 		{3, 2, 1, 0, -1},
// 		{3, 4, 5, 6, -1},
// 	},
// }

// func d(keyCount, hand int, strain [5]float64) int {
// 	if hand == leftHand {

// 	} else {
// 		keyCount / 2
// 	}
// }

type Hand struct {
	Positions    [5]float64
	PrevHoldings [5]bool
	Pressings    [5]bool
}

// TODO: 100ms 넘었으면 pressed 건은 pos 리셋

// To press the note, it should starts at the top (0)
// and keep pressing toward the bottom (1). Hence, we will
// correct the position first before press the notes.
// Meanwhile, no need to correct the pos for releasing as
// they are always at the bottom.
func (h *Hand) calcStrain() [5]float64 {
	// 1. Release all non-holding keys for preparing pressing.
	// 2. Press the keys
	var strains [5]float64
	for _, ps := range [2][5]bool{h.PrevHoldings, h.Pressings} {
		reqMoves := calcReqMoves(h.Positions, ps)
		partialStrains := calcStrains(reqMoves)
		for i, s := range partialStrains {
			strains[i] += s
		}
		respMoves := calcMoves(partialStrains)
		for i, m := range respMoves {
			h.Positions[i] += m
		}
	}
	return strains
}

func calcReqMoves(poses [5]float64, ps [5]bool) [5]float64 {
	var reqMoves [5]float64
	for i, pos := range poses {
		if ps[i] { // pressed
			reqMoves[i] = 1.0 - pos
		} else {
			reqMoves[i] = 0.0 - pos
		}
	}
	return reqMoves
}

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
func calcStrains(dps [5]float64) [5]float64 {
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
		strains[fin] = xdata[i]
	}
	return strains
}

func calcMoves(strains [5]float64) [5]float64 {
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
