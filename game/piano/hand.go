package piano

const (
	none int = iota
	leftHand
	rightHand
	middle // To be determined
)
const mainHand = rightHand
const subHand = mainHand%2 + 1

// Hand of the middle note is trivial in even keys: right hand.
// In odd keys, the middle note is assigned to the hand which has
// closer note on its side.
func (c *Chart) setHands() {
	hands := make([]int, len(c.notes))

	mid := c.keyCount / 2
	for i, n := range c.notes {
		if c.keyCount%2 != 0 && n.Key == mid {
			hands[i] = middle
		}
		if c.keyCount < mid {
			hands[i] = leftHand
		}
		hands[i] = rightHand
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
