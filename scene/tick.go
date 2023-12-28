package scene

var (
	shortTicks = 5
	longTicks  = 20
)

var lastTPS float64 = 60 // default value of ebiten

// https://go.dev/play/p/NgTdSwjyCXC
func UpdateTPS(tps float64) {
	newTPS := float64(tps)
	scale := newTPS / lastTPS
	shortTicks = int(float64(shortTicks) * scale)
	longTicks = int(float64(longTicks) * scale)
	lastTPS = newTPS
}
