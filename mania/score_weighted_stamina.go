package mania

import "math"

const (
	instantScale = 0.025
	instantDecay = 0.1
	gradualScale = 0.005
	gradualDecay = 0.5
)

func (c *Chart) calcStamina() {
	var instant, gradual float64
	prevTimes := make([]int64, c.KeyCount)
	for i, n := range c.Notes {
		time := n.Time - prevTimes[n.key]
		instant *= DecayFactor(instantDecay, time)
		instant += instantScale * n.strain
		gradual *= DecayFactor(gradualDecay, time)
		gradual += gradualScale * n.strain
		if instant < 0 || gradual < 0 {
			panic("negative stamina")
		}
		c.Notes[i].stamina = instant + gradual
		prevTimes[n.key] = n.Time
	}
}

// Difficulty relating
func DecayFactor(decayBase float64, time int64) float64 {
	return math.Pow(decayBase, float64(time)/1000)
}
