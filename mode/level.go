package mode

import "sort"

func Level(c Chart) float64 {
	const decayFactor = 0.95

	ds := c.Difficulties()
	sort.Slice(ds, func(i, j int) bool { return ds[i] > ds[j] })

	sum, weight := 0.0, 1.0
	for _, term := range ds {
		sum += weight * term
		weight *= decayFactor
	}

	// No additional Math.Pow; it would make a little change.
	return sum
}

func DifficultyPieceTimes(dys []*Dynamic, chartDuration int64) (times []int64, durations []int64) {
	const (
		minDuration = 400  // 400ms. 2 beats with 300 BPM
		maxDuration = 1000 // 1000ms. 2 beats with 120 BPM
	)
	times = make([]int64, 0, 300)

	const meter = 2
	beatTimes := BeatTimes(dys, chartDuration, meter)

	var accDuration int64 // accumulated duration
	for i, time := range beatTimes[1:] {
		var prevTime int64
		if i == 0 {
			prevTime = beatTimes[0]
		} else {
			prevTime = times[len(times)-1]
		}
		duration := time - prevTime
		accDuration += duration
		switch {
		case accDuration < minDuration:
			continue

		// Todo: not tested
		// However, this case is not likely to happen.
		case accDuration > maxDuration:
			accDuration -= duration // go back
			unit := float64(duration)
			for accDuration+int64(unit) > maxDuration {
				unit /= 2
			}
			for t := float64(prevTime) + unit; int64(t+0.1) < time; t += unit {
				times = append(times, int64(t))
			}
			// for d := float64(accDuration) + unit; d+unit < maxDuration; d += unit {
			// 	times = append(times, prevTime+int64(d))
			// }
			accDuration = 0
		default:
			times = append(times, time)
			accDuration = 0
		}
	}

	durations = make([]int64, 0, len(times))
	for i, time := range times {
		var d int64
		if i == 0 {
			d = beatTimes[0] - time
		} else {
			d = time - times[i-1]
		}
		durations = append(durations, d)
	}
	return
}
