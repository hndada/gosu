package draws

// // Move to the next Tween if the current one is finished
// switch {
// case !tws.yoyo:
// 	// Standard behavior: increment tick until maxTick
// 	if tws.index < len(tws.Tweens)-1 {
// 		tws.index++
// 	} else {
// 		tws.loop++
// 		if tws.loop < tws.maxLoop {
// 			tws.index = 0
// 		}
// 	}
// case tws.yoyo && !tws.backward:
// 	// Yoyo mode - increasing tick
// 	if tws.index < len(tws.Tweens)-1 {
// 		tws.index++
// 	} else {
// 		tws.backward = true
// 	}
// case tws.yoyo && tws.backward:
// 	// Yoyo mode - decreasing tick
// 	if tws.index > 0 {
// 		tws.index--
// 	} else {
// 		tws.backward = false
// 		tws.loop++
// 	}
// }