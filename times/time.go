package times

// Time is a point of time.
// Duration a length of time.
func ToTick(time int64, tps int) int { return int(float64(time) / 1000 * float64(tps)) }
func ToTime(tick int, tps int) int64 { return int64(float64(tick) / float64(tps) * 1000) }
