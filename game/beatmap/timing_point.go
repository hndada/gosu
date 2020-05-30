package beatmap

import (
	"strings"

	"github.com/hndada/gosu/game/tools"
)

type TimingPoint struct {
	Time        int
	Bpm         float64
	SpeedScale  float64
	Uninherited bool
	Kiai        bool
}

func parseTimingPoint(line string) TimingPoint {
	splitValues := strings.Split(line, `,`)
	timingPoint := TimingPoint{
		Time:        tools.Atoi(splitValues[0]),
		Uninherited: tools.Atoi(splitValues[6]) != 0,
		Kiai:        tools.Atoi(splitValues[7])&1 != 0,
	}
	beatLength := tools.Atof(splitValues[1])
	v := timingPoint.Uninherited
	switch v {
	case true:
		timingPoint.Bpm = 1000 * 60 / beatLength
	case false:
		timingPoint.SpeedScale = 100 / (-beatLength)
	}
	return timingPoint
}
