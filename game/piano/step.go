package piano

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
