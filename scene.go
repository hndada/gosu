package main

func NewScene() {
	// Notes 준비 단계
	prevs := make([]int, c.KeyCount)
	for k := range prevs {
		prevs[k] = -1 // no found
	}
	for next, n := range c.Notes {
		prev := prevs[n.Key]
		c.Notes[next].prev = prev
		if prev != -1 {
			c.Notes[prev].next = next
		}
		prevs[n.Key] = next
	}
	for _, lastIdx := range prevs {
		c.Notes[lastIdx].next = -1
	}
}
