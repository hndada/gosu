package gosu

// func BarPositions(transPoints []*TransPoint, endTime int64) []float64 {
// 	ps := make([]float64, 0)
// 	first := transPoints[0]
// 	first = first.FetchPresent()
// 	var margin int64 = 5000
// 	if margin > first.Time {
// 		margin = first.Time
// 	}
// 	speed := first.BPM * first.BeatLengthScale
// 	for t := float64(first.Time); t >= float64(-margin); t -= float64(first.Meter) * 60000 / first.BPM {
// 		p := speed * (t - 0)
// 		ps = append([]float64{p}, ps...)
// 	}

// 	ps = ps[:len(ps)-1] // Drop for avoiding duplicated
// 	for tp := first; tp != nil; tp = tp.NextBPMPoint.FetchPresent() {
// 		nextTime := endTime + margin
// 		if tp.NextBPMPoint != nil {
// 			nextTime = tp.NextBPMPoint.Time
// 		}
// 		speed := tp.BPM * tp.BeatLengthScale
// 		unit := float64(tp.Meter) * 60000 / tp.BPM
// 		for t := float64(tp.Time); t < float64(nextTime); t += unit {
// 			p := ps[len(ps)-1] + speed*unit
// 			ps = append(ps, p)
// 		}
// 	}
// 	return ps
// }

// func BarTimes(transPoints []*TransPoint, endTime int64) []int64 {
// 	ts := make([]int64, 0)
// 	first := transPoints[0]
// 	first.FetchPresent()
// 	var margin int64 = 5000
// 	if margin > first.Time {
// 		margin = first.Time
// 	}
// 	for t := float64(first.Time); t >= float64(-margin); t -= float64(first.Meter) * 60000 / first.BPM {
// 		ts = append([]int64{int64(t)}, ts...)
// 	}
// 	// Bar for the first TransPoint has already appended.
// 	for tp := first.NextBPMPoint; tp != nil; tp = tp.NextBPMPoint {
// 		next := endTime + margin
// 		if tp.NextBPMPoint != nil {
// 			next = tp.NextBPMPoint.Time
// 		}
// 		unit := float64(tp.Meter) * 60000 / tp.BPM
// 		for t := float64(tp.Time); t < float64(next); t += unit {
// 			ts = append(ts, int64(t))
// 		}
// 	}
// 	return ts
// }
