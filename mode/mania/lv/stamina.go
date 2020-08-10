package lv

import (
	"github.com/hndada/gosu/internal/tools"
	"github.com/hndada/gosu/mode/mania"
)

const (
	instantScale = 0.025
	instantDecay = 0.1
	gradualScale = 0.005
	gradualDecay = 0.5
)

func CalcStamina(ns []mania.Note, keymode int) {
	var instant, gradual float64
	var elapsedTime float64
	prevTimes := tools.GetIntSlice(keymode, 0)
	for i, n := range ns {
		elapsedTime = float64(n.Time - prevTimes[n.Key])
		instant = instant * tools.DecayFactor(instantDecay, elapsedTime)
		instant += instantScale * n.Strain
		gradual = gradual * tools.DecayFactor(gradualDecay, elapsedTime)
		gradual += gradualScale * n.Strain
		if instant < 0 || gradual < 0 {
			panic("negative stamina")
		}
		ns[i].Stamina = instant + gradual
		prevTimes[n.Key] = n.Time
	}
}
