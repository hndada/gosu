package gosu

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/config"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/mode/mania"
)

type NoteSprite struct {
	noteType mode.NoteType
	kind     config.Kind
	x        float64
	y        float64
	op       ebiten.DrawImageOptions
}

type LNSprite struct {
	kind   config.Kind
	head   *NoteSprite
	tail   *NoteSprite
	bodyop ebiten.DrawImageOptions
}

func (s LNSprite) height() float64 {
	return s.tail.y - s.head.y
}

func (s *SceneMania) setNoteSprites() {
	{
		s.notes = make([]NoteSprite, len(s.chart.Notes))
		var i int
		var n mania.Note
		var offset float64
		sfactors := s.chart.TimingPoints.SpeedFactors
		for si, sp := range sfactors {
			for n.Time < sp.Time {
				// kind, noteType, x 설정
				s.notes[i].y = float64(n.Time-sp.Time)*sp.Factor + offset
				i++
				n = s.chart.Notes[i]
			}
			if si < len(sfactors) {
				offset += float64(sfactors[si+1].Time-sp.Time) * sp.Factor
			}
		}
	}
	{
		s.lnotes = make([]LNSprite, s.chart.NumLN())
		lastLNHeads := make([]int, s.chart.Keys)
		for i, n := range s.chart.Notes {
			switch n.Type {
			case mania.TypeLNHead:
				lastLNHeads[n.Key] = i
			case mania.TypeLNTail:
				var ls LNSprite
				ls.head = &s.notes[lastLNHeads[n.Key]]
				ls.tail = &s.notes[i]
				ls.kind = ls.head.kind
			}
		}
	}
}

// op에 값 적용하는 함수
// speed와 hitPosition을 적용
func (s *SceneMania) applySpeed(speed float64) {
	s.speed = speed
	var hitPosition float64
	for i, n := range s.notes {
		// todo: scaled x 구하기
		var x float64
		y := (-(n.y - s.progress) + hitPosition) * speed
		s.notes[i].op.GeoM.Reset()
		s.notes[i].op.GeoM.Translate(x, y)
	}
	for i, n := range s.lnotes {
		// 이미지 불러오기
		// todo: n.height() 활용해서 스케일 구하기
		s.lnotes[i].bodyop.GeoM.Reset()
		s.lnotes[i].bodyop.GeoM.Translate(n.tail.x, n.tail.y)
		s.lnotes[i].bodyop.GeoM.Scale() // body 생성
	}
}
