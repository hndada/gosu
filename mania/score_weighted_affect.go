package mania

const noFound = -1

func (c *Chart) markAffect() {
	// scan the notes whether its affectable, chordable or not
	// further lanes which cause *miss* when be hit at the same time goes 'chuck cutter'
	prevs := make([]int, c.KeyCount)
	for k := range prevs {
		prevs[k] = -1
	}
	for i, n := range c.Notes {
		c.markPrevAffect(i, prevs)
		c.markNextAffect(i)
		prevs[n.Key] = i
	}
}

func (c *Chart) markPrevAffect(i int, prevs []int) {
	n := c.Notes[i]
	for prevKey, prevIdx := range prevs {
		if prevIdx == noFound {
			continue
		}
		prevNote := c.Notes[prevIdx]
		time := n.Time - prevNote.Time
		switch prevNote.Key == n.Key {
		case true: // jack
			if time <= maxDeltaJack {
				c.Notes[i].trillJack[prevKey] = prevIdx
			}
		default:
			if time <= maxDeltaTrill {
				if time <= maxDeltaChord { // chord
					c.Notes[i].chord[prevKey] = prevIdx
					c.Notes[i].trillJack[prevKey] = noFound
				} else { // trill
					c.Notes[i].trillJack[prevKey] = prevIdx
				}
			}
		}
	}
	c.Notes[i].chord[n.Key] = i // putting note itself to chord
}

func (c *Chart) markNextAffect(i int) {
	n := c.Notes[i]
	nextIdx := i + 1
	for nextIdx < len(c.Notes) {
		nextNote := c.Notes[nextIdx]
		time := nextNote.Time - n.Time
		if time > maxDeltaTrill {
			break
		}
		if nextNote.Type != TypeLNTail &&
			nextNote.Key != n.Key && // jack is not relevant
			c.Notes[i].chord[nextNote.Key] == noFound { // prev notes is prior to next notes
			switch {
			case time <= maxDeltaChord:
				c.Notes[i].chord[nextNote.Key] = nextIdx
				// default:
				// 	ns[i].chord[nextNote.Key] = cut
			}
		}
		nextIdx++
	}
}
