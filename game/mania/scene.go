package mania

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hndada/gosu/game"

	"image"
	_ "image/jpeg"
)

type Scene struct {
	game.PlayScene
	sceneSettings // constant
	sceneChart    // constant
	sceneImage    // constant
	sceneNotes    // partially constant
	sceneState
	sceneTally

	ready      bool      // whether scene has been loaded
	startTime  time.Time // int64
	initUpdate bool
	auto       func(int64) []keyEvent
	playSE     func()
}
type sceneSettings struct {
	speed        float64
	hitPosition  float64
	displayScale float64
	// layout      []types.VKCode
}
type sceneChart struct {
	mods         Mods
	chart        *Chart
	endTime      int64
	speedFactors []game.SpeedFactorPoint
	timeStamps   []timeStamp
}
type sceneImage struct {
	bg               *ebiten.Image
	bgop             *ebiten.DrawImageOptions
	fixedStageSprite *ebiten.Image
}

type sceneNotes struct {
	notes  []NoteSprite
	lnotes []LNSprite // 롱노트 특성상, 2개로 나누는 게 불가피해보임
}
type sceneState struct {
	done         bool
	timeStampIdx int
	lastPressed  []bool
	staged       []int
}
type sceneTally struct { // set of variouis scores in the scene
	score float64
	karma float64
	hp    float64
	combo int32
}

type timeStamp struct {
	time     int64
	nextTime int64
	position float64
	factor   float64
}

func newSceneSettings() sceneSettings {
	s := new(sceneSettings)
	switch {
	default:
		s.speed = Settings.GeneralSpeed
	}
	s.hitPosition = Settings.HitPosition
	s.displayScale = game.ScaleY()
	// s.layout = Settings.KeyLayout[s.chart.Keys]
	return *s
}

func newSceneChart(c *Chart, mods Mods) sceneChart {
	s := new(sceneChart)
	s.mods = mods
	s.chart = c.ApplyMods(s.mods)
	s.endTime = s.chart.EndTime()

	initSpeedFactor := game.SpeedFactorPoint{Time: 0, Factor: 1}
	s.speedFactors = append([]game.SpeedFactorPoint{initSpeedFactor}, s.chart.TimingPoints.SpeedFactors...)
	s.timeStamps = make([]timeStamp, len(s.speedFactors))
	var position float64
	for i, sf := range s.speedFactors {
		timeStamp := timeStamp{
			time:     sf.Time,
			position: position,
			factor:   sf.Factor,
		}
		if i < len(s.speedFactors)-1 {
			nextTime := s.speedFactors[i+1].Time
			timeStamp.nextTime = nextTime
			position += float64(nextTime-sf.Time) * sf.Factor
		} else {
			timeStamp.nextTime = 9223372036854775807 // max int64
		}
		s.timeStamps[i] = timeStamp
	}
	return *s
}
func newSceneImage(c *Chart) sceneImage {
	s := new(sceneImage)
	bg, err := c.Background()
	if err != nil {
		panic(err)
	}
	s.bg = bg
	s.bgop = game.BackgroundOp(game.ScreenSize(), image.Pt(s.bg.Size()))
	var dimness uint8
	switch {
	default:
		dimness = game.GeneralDimness()
	}
	s.bgop.ColorM.ChangeHSV(0, 1, float64(dimness)/100)
	s.fixedStageSprite = SpriteMap.Stages[c.Keys].Fixed.Image()
	return *s
}

func newSceneNotes(c *Chart, timeStamps []timeStamp) sceneNotes {
	s := new(sceneNotes)
	stageSprite := SpriteMap.Stages[c.Keys]
	s.notes = make([]NoteSprite, len(c.Notes))
	for i, n := range c.Notes {
		var ns NoteSprite
		var sprite game.Sprite
		switch n.Type {
		case TypeNote:
			sprite = stageSprite.Notes[n.Key]
		case TypeLNHead:
			sprite = stageSprite.LNHeads[n.Key]
		case TypeLNTail:
			sprite = stageSprite.LNTails[n.Key]
		}
		ns.x = float64(sprite.Position().X)
		ns.i = sprite.Image()
		_, ns.h = ns.i.Size()
		ns.op = &ebiten.DrawImageOptions{}
		s.notes[i] = ns
	}

	// set position of notes
	// should be precede lnotes loading
	// performance: range timeStamps를 outer loop로 두고 짜도 큰 차이 없을 듯
	var timeStampIdx int
	ts := timeStamps[0]
	for i, n := range c.Notes {
		for si := range timeStamps[timeStampIdx:] {
			if n.Time < timeStamps[timeStampIdx+si].nextTime {
				if si != 0 {
					ts = timeStamps[timeStampIdx+si]
					timeStampIdx += si
				}
				break
			}
		}
		s.notes[i].position = float64(n.Time-ts.time)*ts.factor + ts.position
	}
	s.lnotes = make([]LNSprite, 0, c.LNCount())
	lastLNHeads := make([]int, c.Keys)
	for i, n := range c.Notes {
		switch n.Type {
		case TypeLNHead:
			lastLNHeads[n.Key] = i
		case TypeLNTail:
			ls := LNSprite{
				head: &s.notes[lastLNHeads[n.Key]],
				tail: &s.notes[i],

				i:      stageSprite.LNBodys[n.Key][0].Image(),
				bodyop: &ebiten.DrawImageOptions{},
			}
			ls.length = ls.tail.position - ls.head.position
			s.lnotes = append(s.lnotes, ls)
		}
	}

	return *s
}

func newSceneState(c *Chart) sceneState {
	s := new(sceneState)
	s.done = false
	s.timeStampIdx = 0
	s.lastPressed = make([]bool, c.Keys)
	s.staged = make([]int, c.Keys)
	for k := range s.staged {
		s.staged[k] = -1
	}
	for k := range s.staged {
		for i, n := range c.Notes {
			if n.Key == k {
				s.staged[k] = i
				break
			}
		}
	}
	return *s
}
func newSceneTally() sceneTally {
	s := new(sceneTally)
	s.score = 0
	s.karma = 100
	s.hp = 100
	s.combo = 0
	return *s
}

func NewScene(c *Chart, mods Mods) *Scene {
	s := new(Scene)
	s.sceneSettings = newSceneSettings()
	s.sceneChart = newSceneChart(c, mods)
	s.sceneImage = newSceneImage(c)
	s.sceneNotes = newSceneNotes(c, s.timeStamps)
	s.sceneState = newSceneState(c)
	s.sceneTally = newSceneTally()
	s.AudioPlayer = game.NewAudioPlayer(s.chart.AbsPath(s.chart.AudioFilename))
	s.AudioPlayer.Play()
	s.AudioPlayer.Pause()
	s.auto = s.chart.GenAutoKeyEvents()
	s.playSE = SEPlayer()
	s.ready = true
	return s
}

func (s *Scene) Ready() bool { return s.ready }

// 음악이 그림보다 느리다
// 여러번 로딩되면 음악이 그림 속도를 따라가다
func (s *Scene) Update() error {
	s.Tick++
	if s.Tick < 20 { // 별 효용성이 없는 듯
		return nil
	}
	var now int64
	if !s.initUpdate {
		ebiten.SetWindowTitle(fmt.Sprintf("gosu - %s [%s]", s.chart.MusicName, s.chart.ChartName))
		s.AudioPlayer.Play()
		// s.Tick = 0 // 초기화
		s.startTime = time.Now()
		s.initUpdate = true
	} else {
		// now := s.Time()
		now = time.Now().Sub(s.startTime).Milliseconds() // - 180
		// unsafe: 꼬로록 소리 남
		// if now < 3000 {
		//		s.AudioPlayer.Seek(time.Now().Sub(s.startTime))
		//		}
		if ebiten.IsKeyPressed(ebiten.KeyEscape) || now > s.endTime+2000 { // temp: 2초 여유 두기
			_ = s.AudioPlayer.Close()
			s.done = true
		}
	}
	// fmt.Println(s.AudioPlayer.Time())
	var ts timeStamp
	for si := range s.timeStamps[s.timeStampIdx:] {
		if now < s.timeStamps[s.timeStampIdx+si].nextTime {
			ts = s.timeStamps[s.timeStampIdx+si]
			s.timeStampIdx += si
			break
		}
	}

	measure := float64(now-ts.time)*ts.factor + ts.position
	for i, n := range s.notes {
		pos := (n.position-measure)*s.speed - s.hitPosition
		s.notes[i].y = -pos*s.displayScale + float64(n.h)/2
	}
	for i, n := range s.lnotes {
		s.lnotes[i].height = n.length * s.speed * s.displayScale
	}

	// judge: score과 staged도 따라서 업데이트
	for _, e := range s.auto(now) {
		s.judge(e)
	}

	// 따로 처리: lost, scored되고 시간 다 된 LNTail
	lost := func(timeDiff int64) bool { return timeDiff < -bad.Window } // never hit
	flushable := func(n Note, timeDiff int64) bool { return n.scored && timeDiff < miss.Window }
	for k, i := range s.staged {
		if i < 0 {
			continue
		}
		n := s.chart.Notes[i]
		timeDiff := n.Time - now

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
	now := time.Now().Sub(s.startTime).Milliseconds()
	screen.DrawImage(s.bg, s.bgop)
	screen.DrawImage(s.fixedStageSprite, &ebiten.DrawImageOptions{})
	for _, n := range s.lnotes {
		n.DrawLN(screen)
	}
	for i, n := range s.notes {
		s.notes[i].op.GeoM.Reset()
		s.notes[i].op.GeoM.Translate(n.x, n.y)
		screen.DrawImage(n.i, s.notes[i].op)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		`CurrentFPS: %.2f
CurrentTPS: %.2f
Time: %.3fs

score: %.0f
karma: %.2f
hp: %.2f
combo: %d
`, ebiten.CurrentFPS(), ebiten.CurrentTPS(), float64(now)/1000,
		s.score, s.karma, s.hp, s.combo))
}

func (s *Scene) Done(args *game.TransSceneArgs) bool {
	if s.done && args.Next == "" {
		args.Next = "sceneSelect"
		args.Args = nil
	}
	return s.done
}
