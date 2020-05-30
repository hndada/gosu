package mania

import "github.com/hndada/gosu/game/tools"

const thumb = 0
const (
	left  = -1
	alter = 0
	right = 1
)

// currently no considering right-scratched 8K map
func hand(key, keymode int) int {
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
}
func finger(key int, keymode int) int {
	switch keymode {
	case 8: // based on left-scratch
		if key == 0 {
			return 4
		}
		return finger(key-1, 7)
	case 2, 4, 6: //no thumb
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
