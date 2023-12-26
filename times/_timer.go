package times

import "time"

type Timer struct {
	StartTime time.Time
}

// Since returns the time elapsed since t, considering playback rates.
// since() was not exported directly to prevent confusion.
func (t Timer) Since() time.Duration {
	return since(t.StartTime)
}
