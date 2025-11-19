package piano

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

const maxStepTimeWindow = 30

// Required preprocess:
// 1. Sort by time and key
// 2. Set prev and next note
func (c Chart) calcSteps() {
	c.calcHands()
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

	x := new(mat.VecDense)
	a := mat.NewDense(rank, rank, infl)
	b := mat.NewVecDense(rank, dps2)
	x.SolveVec(a, b) // solves a*x = b
	xdata := x.RawVector().Data

	var strains [5]float64
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

// To press the note, it should starts at the top (0)
// and keep pressing toward the bottom (1). Hence,
// we will correct the position first before press the notes.
// Meanwhile, no need to correct the pos for releasing as
// they are always at the bottom.
func a() {
	var poses [5]float64
	// 1. Release all non-holding keys for preparing pressing.
	// 2. Press the keys
	var holdings [5]bool
	var pressings [5]bool
	for _, ps := range [2][5]bool{holdings, pressings} {
		reqMoves := calcReqMoves(poses, ps)
		strains := calcStrains(reqMoves) // TODO: each note
		respMoves := calcMoves(strains)
		for i, m := range respMoves {
			poses[i] += m
		}
	}
}
