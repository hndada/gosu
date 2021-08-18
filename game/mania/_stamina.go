package mania

import (
	"github.com/hndada/gosu/game"
)

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
		time := n.Time - prevTimes[n.Key]
		instant *= game.DecayFactor(instantDecay, time)
		instant += instantScale * n.strain
		gradual *= game.DecayFactor(gradualDecay, time)
		gradual += gradualScale * n.strain
		if instant < 0 || gradual < 0 {
			panic("negative stamina")
		}
		c.Notes[i].stamina = instant + gradual
		prevTimes[n.Key] = n.Time
	}
}
