package mode

import "github.com/hajimehoshi/ebiten/v2"

// TPS affects only on Update(), not on Draw().
var TPS = float64(ebiten.TPS())

func ToTick(ms int32) int       { return int(TPS * float64(ms) / 1000) }
func ToTime(tick int) int32     { return int32(float64(tick) / TPS * 1000) }
func ToSecond(ms int32) float64 { return float64(ms) / 1000 }
