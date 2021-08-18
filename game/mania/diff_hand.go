package mania

const (
	thumb = 0

	left  = -1
	alter = 0
	right = 1
)

const defaultHand = right

// only for init
func finger(keyCount, key int) int {
	actualKeyCount := keyCount & ScratchMask
	switch {
	case keyCount&leftScratch != 0:
		if key == 0 {
			return finger(actualKeyCount-1, 0) + 1
		}
		return finger(actualKeyCount-1, key-1)
	case keyCount&rightScratch != 0:
		if key == actualKeyCount-1 {
			return finger(actualKeyCount-1, actualKeyCount-2) + 1
		}
		return finger(actualKeyCount-1, key)
	default:
		if keyCount%2 == 0 {
			var v int
			if key >= keyCount/2 {
				v = key - keyCount/2 + 1
			} else {
				v = keyCount/2 - key
			}
			if keyCount == 10 {
				v--
			}
			return v
		} else {
			v := key - keyCount/2
			if v < 0 {
				return -v
			}
			return v
		}
	}
}

func hand(keyCount, key int) int {
	actualKeyCount := keyCount & ScratchMask
	switch {
	case keyCount&leftScratch != 0:
		if key == 0 {
			return left
		}
		return hand(actualKeyCount-1, key-1)
	case keyCount&rightScratch != 0:
		if key == actualKeyCount-1 {
			return right
		}
		return hand(actualKeyCount-1, key)
	default:
		switch {
		case key < keyCount/2:
			return left
		case key > keyCount/2:
			return right
		default: // key == keyCount/2:
			if keyCount%2 == 0 {
				return right
			}
			return alter // odd key count use thumb, which is alterable
		}
	}
}
func (n *Note) settleAlterHand() {
	// affect idx has already been calculated
	if n.hand != alter {
		return
	}
	// rule 1: use default hand if there is a note very next to alterable note
	if n.chord[n.Key+defaultHand] != -1 {
		n.hand = defaultHand
		return
	}

	// rule 2: the hand which has more notes in the chord
	// rule 3: default hand if each hand has same number of notes
	leftCount, rightCount := 0, 0
	for k := n.Key - 1; k >= 0; k-- {
		if n.chord[k] <= -1 {
			break
		}
		leftCount++
	}
	for k := n.Key + 1; k < len(n.chord); k++ {
		if n.chord[k] <= -1 {
			break
		}
		rightCount++
	}

	switch {
	case leftCount > rightCount:
		n.hand = left
	case leftCount < rightCount:
		n.hand = right
	default: // if two counts are same
		n.hand = defaultHand
	}
}

const (
	outer = iota + 1
	inner
	innerAdj
)

// supposed comparing keys are in same hand
func isHoldOuter(holdKey, key, keymode int) bool {
	h := finger(holdKey, keymode)
	switch h {
	case left:
		return holdKey < key
	case right:
		return key < holdKey
	default: // h is a thumb, which is always excluded
		return false
	}
}

// supposed comparing keys are in same hand
func isHoldInner(holdKey, key, keymode int) bool {
	h := finger(holdKey, keymode)
	switch h {
	case left:
		return key < holdKey
	case right:
		return holdKey < key
	default: // h is a thumb, which is always included
		return true
	}
}

func isHoldInnerAdj(holdKey, key, keymode int) bool {
	// hold note hitting with thumb does not afford adjacent bonus
	h := finger(holdKey, keymode)
	switch h {
	case left:
		return holdKey == key+1
	case right:
		return holdKey == key-1
	default: // thumb
		return false
	}
}
func sameHand(h1, h2 float64) bool {
	if h1 == 0 || h2 == 0 {
		panic("alter yet")
	}
	return h1*h2 > 0
}
