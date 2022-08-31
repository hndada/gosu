package piano

// PlayNote is for in-game. Handled by pointers to modify its fields easily.
// type Note struct {
// 	gosu.Note
// 	Key int
// 	// Marked   bool
// 	// NextTail *PlayNote // For performance of DrawLongNoteBodies()
// }

// NotePosition calculates position, the centered y-axis value.
// y = position - h/2

// func NewPlayNotes(c *Chart) ([]*PlayNote, []*PlayNote, []*PlayNote, float64) {
// 	playNotes := make([]*PlayNote, 0, len(c.Notes))
// 	firstStagedNotes := make([]*PlayNote, c.KeyCount)
// 	firstLowestTails := make([]*PlayNote, c.KeyCount)
// 	prevs := make([]*PlayNote, c.KeyCount)
// 	prevTails := make([]*PlayNote, c.KeyCount)
// 	var weights float64
// 	for _, n := range c.Notes {
// 		prev := prevs[n.Key]
// 		pn := &PlayNote{
// 			Note: n,
// 			Prev: prev,
// 		}
// 		if prev != nil { // Next value is set later.
// 			prev.Next = pn
// 		}
// 		prevs[n.Key] = pn
// 		if firstStagedNotes[n.Key] == nil {
// 			firstStagedNotes[n.Key] = pn
// 		}
// 		if n.Type == Tail {
// 			if prevTails[n.Key] != nil {
// 				prevTails[n.Key].NextTail = pn
// 			}
// 			prevTails[n.Key] = pn
// 			if firstLowestTails[n.Key] == nil {
// 				firstLowestTails[n.Key] = pn
// 			}
// 		}
// 		weights += pn.Weight()
// 		playNotes = append(playNotes, pn)
// 	}
// 	return playNotes, firstStagedNotes, firstLowestTails, weights
// }

// top returns top position of long note body.
// func (s ScenePlay) Top(tail *PlayNote) int {
// 	return int(s.Position(tail.Time)+s.TailSprites[tail.Key].H/2) - 1 // Extra 1 pixel for compensating round-down
// }

// // bottom returns bottom position of long note body.
// func (s ScenePlay) Bottom(tail *PlayNote) int {
// 	return int(s.Position(tail.Prev.Time)-s.HeadSprites[tail.Key].H/2) + 1 // Extra 1 pixel for compensating round-down
// }

// DrawLongNoteBodies draws long sprite with Binary-building method, instead of SubImage.
// DrawLongNoteBodies draws long note before drawing Head or Tail.
// DrawLongNoteBodies just draws sub image of long note body.
// func (s *ScenePlay) DrawLongNoteBodies(screen *ebiten.Image) {
// 	for k, n0 := range s.LowestTails {
// 		for n := n0; n != nil && s.Top(n) >= screenSizeY; n = n.NextTail {
// 			s.LowestTails[k] = n
// 		}
// 		for n := s.LowestTails[k]; n != nil && s.Bottom(n) >= 0; n = n.NextTail {
// 			top := s.Top(n)
// 			bottom := s.Bottom(n)
// 			if top < 0 {
// 				top = 0
// 			}
// 			if bottom > screenSizeY {
// 				bottom = screenSizeY
// 			}

// 			var pow int
// 			y := float64(top)
// 			for length := bottom - top; length > 0; length /= 2 {
// 				if length%2 == 0 {
// 					pow++
// 					continue
// 				}
// 				sprite := s.BodySprites[k][pow]
// 				sprite.Y = y
// 				op := sprite.Op()
// 				if n.Marked {
// 					op.ColorM.ChangeHSV(0, 0.3, 0.3)
// 				}
// 				screen.DrawImage(sprite.I, op)
// 				y += sprite.H // 1 << pow
// 				pow++
// 			}

// 			// Test case: draw long note by SubImage.
// 			{
// 				// sprite := s.BodySpritesTest[k]
// 				// sprite.Y = float64(top)
// 				// op := sprite.Op()
// 				// if n.Marked {
// 				// 	op.ColorM.ChangeHSV(0, 0.3, 0.3)
// 				// }
// 				// rect := sprite.I.Bounds()
// 				// rect.Max.Y = bottom - top

// 				// screen.DrawImage(sprite.I.SubImage(rect).(*ebiten.Image), op)
// 			}
// 		}
// 	}
// }

// Time bound for drawing notes in milliseconds.
// The tight values can be calculated by similar way of NotePosition() does.
// const (
// 	up   = 10 * 1000
// 	down = -2 * 1000
// )

// func (s *ScenePlay) DrawNotes(screen *ebiten.Image) {
// 	for _, n := range s.PlayNotes {
// 		td := n.Time - s.Time()
// 		if td > up || td < down {
// 			continue
// 		}
// 		var sprite draws.Sprite
// 		switch n.Type {
// 		case Head:
// 			sprite = s.HeadSprites[n.Key]
// 		case Tail:
// 			sprite = s.TailSprites[n.Key]
// 		default:
// 			sprite = s.NoteSprites[n.Key]
// 		}
// 		sprite.Y = s.Position(n.Time) - sprite.H/2
// 		op := sprite.Op()
// 		if n.Type == Head {
// 			if n.Next.Marked {
// 				op.ColorM.ChangeHSV(0, 0.3, 0.3)
// 			}
// 		} else {
// 			if n.Marked {
// 				op.ColorM.ChangeHSV(0, 0.3, 0.3)
// 			}
// 		}
// 		screen.DrawImage(sprite.I, op)
// 	}
// }

// Body should be drawn as if one big stretched image.
// How to draw Body:
// 1. Find the section from start to end.
// 2. Stretch the equivalent part of Body sprite.