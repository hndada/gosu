package gosu

// op에 값 적용하는 함수
// hitPosition은 settings 단계에서 미리 적용하고 옴
// func (s *SceneMania) applySpeed(speed float64) {
// 	s.speed = speed
// 	for i, n := range s.notes {
// 		y := (n.y - s.progress) * speed
// 		// s.notes[i].y = y
// 		var sprite graphics.Sprite
// 		switch n.type_ {
// 		case mania.TypeNote:
// 			sprite = s.stage.Notes[n.key]
// 		case mania.TypeLNHead:
// 			sprite = s.stage.LNHeads[n.key]
// 		case mania.TypeLNTail:
// 			sprite = s.stage.LNTails[n.key]
// 		}
// 		sprite.ResetPosition(n.op)
// 		s.notes[i].op.GeoM.Translate(0, y)
// 	}
// 	for i, n := range s.lnotes {
// 		y := n.tail.y
// 		h := n.height() * speed
// 		// if h > 32000 { // todo: 32768 넘어가면 cut
// 		// 	fmt.Println("too long")
// 		// 	h = 32000
// 		// }
// 		s.lnotes[i].i = s.stage.LNBodys[n.key][0].Image(h)
// 		s.stage.LNBodys[n.key][0].ResetPosition(n.bodyop)
// 		s.lnotes[i].bodyop.GeoM.Translate(0, y)
// 	}
// }
