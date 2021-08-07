package mania

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hndada/gosu/game"

	"image"
	_ "image/jpeg"
)

// lnhead와 lntail 분리 유지
// keyboardChannel (kbChan) 삭제
// scene이 최종 패키지에 전부 import 됨 (no cycle dependency)
// 내 생각에 지금 느린건 (혹은 느리다고 보이는건) audio가 Time을 제대로 안내주기 때문인거 같음 -> 맞음
// 이전 값에 상관없이 언제나 다시 그리므로 applySpeed()가 따로 필요 없음
// 최종 이미지는 언제나 사이즈가 int, int이므로 image.Point로 다뤄도 됨
// todo: timing points, (decending/ascending) order로 sort -> rg-parser에서
type Scene struct { // aka Clavier
	game.PlayScene
	mods         Mods
	chart        *Chart
	speedFactors []game.SpeedFactorPoint
	stamps       []timeStamp
	stage        Stage // for quick access
	notes        []NoteSprite
	lnotes       []LNSprite // 롱노트 특성상, 2개로 나누는 게 불가피해보임

	bg   *ebiten.Image
	bgop *ebiten.DrawImageOptions

	speed        float64
	hitPosition  float64
	displayScale float64

	audioPlayer *game.AudioPlayer
	// layout      []types.VKCode
	endTime int64

	score    float64
	karma    float64
	hp       float64
	combo    int32
	stampIdx int

	lastPressed []bool
	staged      []int
	done        bool
}

type timeStamp struct {
	time     int64
	nextTime int64
	position float64
	factor   float64
}

func NewScene(c *Chart, mods Mods) *Scene {
	s := &Scene{}
	s.mods = mods
	s.chart = c.ApplyMods(s.mods)
	// todo: 노트가 언제나 양수 시간에 있다고 상정; 실제로는 노트가 BufferTime보다 뒤에 있을 수 있음
	initSpeedFactor := game.SpeedFactorPoint{Time: 0, Factor: 1}
	s.speedFactors = append([]game.SpeedFactorPoint{initSpeedFactor}, s.chart.TimingPoints.SpeedFactors...)
	s.stamps = make([]timeStamp, len(s.speedFactors))
	var position float64
	for i, sf := range s.speedFactors {
		var stamp timeStamp
		stamp.time = sf.Time
		stamp.position = position
		stamp.factor = sf.Factor
		if i < len(s.speedFactors)-1 {
			nextTime := s.speedFactors[i+1].Time
			stamp.nextTime = nextTime
			position += float64(nextTime-sf.Time) * sf.Factor
		} else {
			stamp.nextTime = 9223372036854775807 // max int64
		}
		s.stamps[i] = stamp
	}
	s.stage = SpriteMap.Stages[s.chart.Keys]
	s.notes = make([]NoteSprite, len(s.chart.Notes))
	for i, n := range s.chart.Notes {
		var ns NoteSprite
		var sprite game.Sprite
		switch n.Type {
		case TypeNote:
			sprite = s.stage.Notes[n.Key]
		case TypeLNHead:
			sprite = s.stage.LNHeads[n.Key]
		case TypeLNTail:
			sprite = s.stage.LNTails[n.Key]
		}
		ns.x = float64(sprite.Position().X)
		ns.i = sprite.Image()
		_, ns.h = ns.i.Size()
		ns.op = &ebiten.DrawImageOptions{}
		s.notes[i] = ns
	}
	// range stamps를 outer loop로 두고 짜도 큰 차이 없을 듯
	var stampIdx int
	stamp := s.stamps[0]
	for i, n := range s.chart.Notes {
		for si := range s.stamps[stampIdx:] {
			if n.Time < s.stamps[stampIdx+si].nextTime {
				if si != 0 {
					stamp = s.stamps[stampIdx+si]
					stampIdx += si
				}
				break
			}
		}
		s.notes[i].position = float64(n.Time-stamp.time)*stamp.factor + stamp.position
	}
	s.lnotes = make([]LNSprite, 0, s.chart.LNCount())
	lastLNHeads := make([]int, s.chart.Keys)
	for i, n := range s.chart.Notes {
		switch n.Type {
		case TypeLNHead:
			lastLNHeads[n.Key] = i
		case TypeLNTail:
			var ls LNSprite
			ls.head = &s.notes[lastLNHeads[n.Key]]
			ls.tail = &s.notes[i]
			ls.i = s.stage.LNBodys[n.Key][0].Image()
			ls.length = ls.tail.position - ls.head.position
			ls.bodyop = &ebiten.DrawImageOptions{}
			s.lnotes = append(s.lnotes, ls)
		}
	}

	var err error
	s.bg, err = s.chart.Background()
	if err != nil {
		panic(err)
	}
	s.bgop = game.BackgroundOp(game.ScreenSize(), image.Pt(s.bg.Size()))
	var dimness uint8
	switch {
	default:
		dimness = game.GeneralDimness()
	}
	s.bgop.ColorM.ChangeHSV(0, 1, float64(dimness)/100)

	switch {
	default:
		s.speed = Settings.GeneralSpeed
	}
	s.hitPosition = Settings.HitPosition
	s.displayScale = game.ScaleY()

	s.audioPlayer = game.NewAudioPlayer(s.chart.AbsPath(s.chart.AudioFilename))
	// s.layout = Settings.KeyLayout[s.chart.Keys]
	s.endTime = s.chart.EndTime()

	s.karma = 100
	s.hp = 100
	s.lastPressed = make([]bool, s.chart.Keys)

	s.staged = make([]int, s.chart.Keys)
	for k := range s.staged {
		s.staged[k] = -1
	}
	for k := range s.staged {
		for i, n := range s.chart.Notes {
			if n.Key == k {
				s.staged[k] = i
				break
			}
		}
	}
	s.done = false
	return s
}

func (s *Scene) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		_ = s.audioPlayer.Close()
		s.done = true
	}

	now := s.Time()
	var stamp timeStamp
	for si := range s.stamps[s.stampIdx:] {
		if now < s.stamps[s.stampIdx+si].nextTime {
			stamp = s.stamps[s.stampIdx+si]
			s.stampIdx += si
			break
		}
	}

	measure := float64(now-stamp.time)*stamp.factor + stamp.position
	for i, n := range s.notes {
		pos := (n.position-measure)*s.speed - s.hitPosition
		s.notes[i].y = -pos*s.displayScale + float64(n.h)/2
	}
	for i, n := range s.lnotes {
		s.lnotes[i].height = n.length * s.speed * s.displayScale
	}
	s.Tick++

	// 매 업데이트마다 score 업데이트; staged도 따라서 업데이트
	var keyEvents []keyEvent
	for _, e := range keyEvents {
		s.judge(e)
	}

	// 따로 처리: lost, scored되고 시간 다 된 LNTail
	lost := func(timeDiff int64) bool { return timeDiff < -bad.Window } // never hit
	flushable := func(n Note, timeDiff int64) bool { return n.scored && timeDiff < miss.Window }
	for k, i := range s.staged {
		n := s.chart.Notes[i]
		timeDiff := n.Time - s.Time()

		if lost(timeDiff) {
			s.applyScore(i, miss)
		}

		if n.Type == TypeLNTail && flushable(n, timeDiff) {
			s.staged[k] = n.next
		}
	}
	return nil
}

func (s *Scene) Draw(screen *ebiten.Image) {
	// hp는 마스크 이미지를 씌우면 되지 않을까
	screen.DrawImage(s.bg, s.bgop)
	screen.DrawImage(s.stage.Fixed.Image(), &ebiten.DrawImageOptions{})
	for _, n := range s.lnotes {
		n.DrawLN(screen)
	}
	// op를 매번 생성하는 게 더 빠를까? 근데 그럴 것 같지는 않아
	// 화면 범위 바깥은 Draw 생략해야할까?
	for i, n := range s.notes {
		s.notes[i].op.GeoM.Reset()
		s.notes[i].op.GeoM.Translate(n.x, n.y)
		screen.DrawImage(n.i, s.notes[i].op)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		`CurrentFPS: %.2f
CurrentTPS: %.2f
Time: %.3fs

score: %f
karma: %.2f
hp: %.2f
combo: %d
`, ebiten.CurrentFPS(), ebiten.CurrentTPS(), float64(s.Time())/1000,
		s.score, s.karma, s.hp, s.combo))
}

func (s *Scene) Init() {
	ebiten.SetWindowTitle(fmt.Sprintf("gosu - %s [%s]", s.chart.MusicName, s.chart.ChartName))
	s.audioPlayer.Play()
}

func (s *Scene) Done() bool {
	return s.done
}
