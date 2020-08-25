package gosu

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/mode/mania"
)


// todo: even, odd 대신 one, two 로 바꾸자
type NoteSprite struct {
	kind     Kind // even, odd, space, pinky
	noteType int  // todo: mania.NoteType
	x        float64
	y        float64
	op       ebiten.DrawImageOptions
}
type LNSprite struct {
	head   *NoteSprite
	tail   *NoteSprite
	bodyop ebiten.DrawImageOptions
}

func (s LNSprite) height() float64 {
	return s.tail.y - s.head.y
}
// op은 노트와 일대일 대응되므로 별도의 타입이나 idx 필요 없음
func (s *SceneMania) setNoteSprites() {
	// speed는 나중에 적용
	// hitPosition도 나중에 적용
	{
		s.notes = make([]NoteSprite, len(s.chart.Notes))
		var i int
		var n mania.Note
		var offset float64
		sfactors := s.chart.TimingPoints.SpeedFactors
		for si, sp := range sfactors {
			for n.Time < sp.Time {
				// kind, notetYpe, x 설정
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
			case mania.LNHead:
				lastLNHeads[n.Key] = i
			case mania.LNTail:
				var ls LNSprite
				ls.head = &s.notes[lastLNHeads[n.Key]]
				ls.tail = &s.notes[i]
			}
		}
	}
	// noteimgs := config.NoteImages(7)
}

// op에 값 적용하는 함수
func (s *SceneMania) applySpeed(speed float64) {
	s.speed = speed
	var hitPosition float64
	for i, n := range s.notes {
		// todo: scaled x 구하기
		var x float64
		y := (-(n.y - s.viewport) + hitPosition) * speed
		s.notes[i].op.GeoM.Reset()
		s.notes[i].op.GeoM.Translate(x, y)
	}
	for i, n := range s.lnotes {
		// 이미지 불러오기
		// n.height() 활용해서 스케일 구하기
		s.lnotes[i].bodyop.GeoM.Reset()
		s.lnotes[i].bodyop.GeoM.Translate(n.tail.x, n.tail.y)
		s.lnotes[i].bodyop.GeoM.Scale() // body 생성
	}
}