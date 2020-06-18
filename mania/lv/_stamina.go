package mania

import (
	"github.com/hndada/gosu/game/tools"
)

const (
	instantScale = 0.025
	instantDecay = 0.1
	gradualScale = 0.005
	gradualDecay = 0.5
)

func (beatmap *ManiaBeatmap) CalcStamina() {
	beatmap.setStamina()
}

func (beatmap *ManiaBeatmap) setStamina() {
	var instant, gradual float64
	var elapsedTime float64
	prevTimes := tools.GetIntSlice(beatmap.Keymode, 0)
	for i, note := range beatmap.Notes {
		elapsedTime = float64(note.Time - prevTimes[note.Key])
		instant = instant * tools.DecayFactor(instantDecay, elapsedTime)
		instant += instantScale * note.Strain
		gradual = gradual * tools.DecayFactor(gradualDecay, elapsedTime)
		gradual += gradualScale * note.Strain
		if instant < 0 {
			panic(&tools.ValError{"instantStamina", tools.Ftoa(instant), tools.ErrFlow})
		}
		if gradual < 0 {
			panic(&tools.ValError{"gradualStamina", tools.Ftoa(gradual), tools.ErrFlow})
		}
		beatmap.Notes[i].Stamina = instant + gradual
		prevTimes[note.Key] = note.Time
	}
}
