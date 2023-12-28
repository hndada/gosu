package times

import "time"

type playbackRateLog struct {
	time    time.Time
	rate    float64
	elapsed time.Duration
}

// Set init time as standard: duration is 0.
var playbackRateLogs = []playbackRateLog{
	{time.Now(), 1.0, 0},
}

// ClearPlaybackRateLogs is proper to be called
// when transitioning to a new scene.
func ClearPlaybackRateLogs() {
	playbackRateLogs = []playbackRateLog{
		{time.Now(), 1.0, 0},
	}
}

func PlaybackRate() float64 {
	return playbackRateLogs[len(playbackRateLogs)-1].rate
}

// time.Now().Sub(t) is not identical with Since(t)
// because Sub does not consider playback rates.
func Now() time.Time {
	log := playbackRateLogs[len(playbackRateLogs)-1]
	d := log.rate * float64(time.Since(log.time))
	return log.time.Add(time.Duration(d))
}

// Since returns the time elapsed since t, considering playback rates.
func Since(t time.Time) time.Duration {
	for i := len(playbackRateLogs) - 1; i >= 0; i-- {
		log := playbackRateLogs[i]
		if log.time.After(t) {
			continue
		}
		// Same analogy as speed, time, and distance:
		// duration = rate * (time difference)
		// As well as calculating position of Dynamics:
		// prev.Position + prev.Speed * float64(d.Time-prev.Time)

		// Multiplying 1.0 does not harm precision.
		// https://go.dev/play/p/5HRusSw8qtP
		td := t.Sub(log.time) // Time difference
		scaled := log.rate * float64(td)
		return log.elapsed + time.Duration(scaled)
	}

	// If the given time is before the first log, return
	// the normal time difference. It is not likely to happen.
	return time.Since(t)
}

func SetPlaybackRate(newRate float64) {
	d := Since(time.Now())
	log := playbackRateLog{time.Now(), newRate, d}
	playbackRateLogs = append(playbackRateLogs, log)
}
