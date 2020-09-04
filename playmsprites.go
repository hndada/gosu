package gosu

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/graphics"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/mode/mania"
)

// type, key, x and y are for "backup" to re-fetch options or redraw
type NoteSprite struct {
	type_ mode.NoteType
	key   int
	x, y  float64 // mania doesnt use x
	i     *ebiten.Image
	op    *ebiten.DrawImageOptions
}

type LNSprite struct {
	key    int
	head   *NoteSprite
	tail   *NoteSprite
	i      *ebiten.Image
	bodyop *ebiten.DrawImageOptions
}

func (s LNSprite) height() float64 {
	return s.head.y - s.tail.y // y는 음수
}

func (s *SceneMania) setNoteSprites() {
	s.notes = make([]NoteSprite, len(s.chart.Notes))
	for i, n := range s.chart.Notes {
		var ns NoteSprite
		ns.type_ = n.Type
		ns.key = n.Key
		switch n.Type {
		case mania.TypeNote:
			ns.i = s.stage.Notes[n.Key].Image()
		case mania.TypeLNHead:
			ns.i = s.stage.LNHeads[n.Key].Image()
		case mania.TypeLNTail:
			ns.i = s.stage.LNTails[n.Key].Image()
		}
		ns.op = &ebiten.DrawImageOptions{}
		s.notes[i] = ns
	}
	{
		var ni int
		var offset float64
	outer:
		for si, sp := range s.sfactors {
			if si == len(s.sfactors)-1 {
				for i, n := range s.chart.Notes[ni:] {
					s.notes[i].y = -float64(n.Time-sp.Time)*sp.Factor + offset
				}
			} else {
				nextsp := s.sfactors[si+1]
				var n mania.Note
				for n.Time < nextsp.Time {
					s.notes[ni].y = -float64(n.Time-sp.Time)*sp.Factor + offset
					ni++
					if ni >= len(s.chart.Notes) { // 어느 speed factor point가 마지막 노트보다 뒤에 있을 때
						break outer
					}
					n = s.chart.Notes[ni]
				}
				offset += -float64(nextsp.Time-sp.Time) * sp.Factor
			}
		}
	}
	s.lnotes = make([]LNSprite, 0, s.chart.NumLN())
	lastLNHeads := make([]int, s.chart.Keys)
	for i, n := range s.chart.Notes {
		switch n.Type {
		case mania.TypeLNHead:
			lastLNHeads[n.Key] = i
		case mania.TypeLNTail:
			var ls LNSprite
			ls.key = n.Key
			ls.head = &s.notes[lastLNHeads[n.Key]]
			ls.tail = &s.notes[i]
			ls.bodyop = &ebiten.DrawImageOptions{}
			s.lnotes = append(s.lnotes, ls)
		}
	}
}

// op에 값 적용하는 함수
// hitPosition은 settings 단계에서 미리 적용하고 옴
// todo: 판정선 가운데에 노트 가운데가 맞을 때 Max가 뜨게
func (s *SceneMania) applySpeed(speed float64) {
	s.speed = speed
	for i, n := range s.notes {
		y := (n.y - s.progress) * speed
		// s.notes[i].y = y
		var sprite graphics.Sprite
		switch n.type_ {
		case mania.TypeNote:
			sprite = s.stage.Notes[n.key]
		case mania.TypeLNHead:
			sprite = s.stage.LNHeads[n.key]
		case mania.TypeLNTail:
			sprite = s.stage.LNTails[n.key]
		}
		sprite.ResetPosition(n.op)
		s.notes[i].op.GeoM.Translate(0, y)
	}
	// todo: animation
	for i, n := range s.lnotes {
		y := n.tail.y
		h := n.height() * speed
		s.lnotes[i].i = s.stage.LNBodys[n.key][0].Image(h)
		s.stage.LNBodys[n.key][0].ResetPosition(n.bodyop)
		s.lnotes[i].bodyop.GeoM.Translate(0, y)
	}
}
