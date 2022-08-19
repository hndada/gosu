package ctrl

const defaultTPS = 1000

const (
	transDuration = 400 // In milliseconds.
	longDuration  = 200
	shortDuration = 80
)

var (
	transCountdown = int(float64(transDuration) / 1000 * defaultTPS)
	longCountdown  = int(float64(longDuration) / 1000 * defaultTPS)
	shortCountdown = int(float64(shortDuration) / 1000 * defaultTPS)
)

func SetTPS(v float64) {
	transCountdown = int(float64(transDuration) / 1000 * v)
	longCountdown = int(float64(longDuration) / 1000 * v)
	shortCountdown = int(float64(shortDuration) / 1000 * v)
}
