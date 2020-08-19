package mania

import (
	"github.com/hndada/gosu/internal/tools"
	"github.com/hndada/gosu/mode/mania/lv"
)

const thumb = 0
const (
	left  = -1
	alter = 0
	right = 1
)

const defaultHand = right

func finger(key int, keymode int) int {
	switch keymode {
	case 8: // based on left-scratch
		if key == 0 {
			return 4
		}
		return finger(key-1, 7)
	case 2, 4, 6: // no thumb
		if key >= keymode/2 {
			return key - keymode/2 + 1
		}
	case 10: // no alter; use both thumbs
		if key < keymode/2 {
			return keymode/2 - key - 1
		}
	}
	return tools.AbsInt(key - keymode/2)
}

// supposed comparing keys are in same hand
func isHoldOuter(holdKey, key, keymode int) bool {
	h := hand(holdKey, keymode)
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
	h := hand(holdKey, keymode)
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
	h := hand(holdKey, keymode)
	switch h {
	case left:
		return holdKey == key+1
	case right:
		return holdKey == key-1
	default: // thumb
		return false
	}
}

func isSameHand(h1, h2 int) bool {
	return tools.IsIntSameSign(h1, h2)
}

func hand(ns []Note, keymode int) {
	for i, n := range ns { // naive hand
		ns[i].hand = func(key, keymode int) int {
			// currently no considering right-scratched 8K map, but its fine since we can consider MR-ed one
			switch {
			case key < keymode/2:
				return left
			case key > keymode/2:
				return right
			default: // key==keymode/2
				switch {
				case keymode == 8:
					return alter // 8K mode is actually 7+1 key mode
				case keymode%2 == 0:
					return right // even keymode use no alterable thumb
				default:
					return alter // odd keymode use thumb, which is alterable
				}
			}
		}(n.Key, keymode)
	}

	for i, n := range ns {
		// affect idx has already been calculated
		if n.hand != alter {
			continue
		}

		// rule 1: use default hand if there is a note very next to alterable note
		if n.chord[n.Key+defaultHand] != lv.noFound {
			ns[i].hand = defaultHand
			continue
		}

		// rule 2: the hand which has more notes in the chord
		// rule 3: default hand if each hand has same number of notes
		leftCount, rightCount := 0, 0
		for key := n.Key - 1; key >= 0; key-- {
			if n.chord[key] <= lv.noFound {
				break
			}
			leftCount++
		}
		for key := n.Key + 1; key < len(n.chord); key++ {
			if n.chord[key] <= lv.noFound {
				break
			}
			rightCount++
		}

		switch {
		case leftCount > rightCount:
			ns[i].hand = left
		case leftCount < rightCount:
			ns[i].hand = right
		default: // if two counts are same
			ns[i].hand = defaultHand
		}
	}
}
