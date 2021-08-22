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

type TimeBool struct {
	Time  int64
	Value bool
}
type Scene struct {
	game.Scene

	speed float64
	bg    *ebiten.Image
	bgop  *ebiten.DrawImageOptions

	mods  Mods
	chart *Chart

	lastPressed []TimeBool // todo: time도 있어야 여러 프레임동안 pressed할 수 있음
	staged      []int

	score float64
	karma float64
	hp    float64
	combo int32

	done bool

	ready       bool      // whether scene has been loaded
	startTime   time.Time // int64
	initUpdate  bool
	auto        func(int64) []keyEvent
	playSE      func()
	judgeCounts [len(Judgments)]int

	timeStamp func(time int64) game.TimeStamp

	sceneUI
	lns []game.LongSprite

	timeDiffs []int64
	jm        *game.JudgmentMeter // temp
	lastJudge game.Judgment
}

func NewScene(c *Chart, mods Mods, p image.Point, cwd string) *Scene {
	const instability = 0 // 0~100; 0 is Auto
	s := new(Scene)
	// s.CWD = cwd
	s.ScreenSize = p
	s.speed = Settings.GeneralSpeed
	s.mods = mods
	s.chart = c.ApplyMods(s.mods)

	const dimness = 30
	bg, err := c.Background()
	if err != nil {
		panic(err)
	}
	s.bg = bg
	s.bgop = game.BackgroundOp(p, image.Pt(s.bg.Size()))
	s.bgop.ColorM.ChangeHSV(0, 1, float64(dimness)/100)

	var img *ebiten.Image
	kinds := keyKindsMap[c.KeyCount]
	for i, n := range c.Notes { // temp: Note, LNHead, LNTail 전부 Note 이미지 사용
		img = Skin.Note[kinds[n.Key]]
		s.chart.Notes[i].Sprite.SetImage(img)
	}
	s.lastPressed = make([]TimeBool, c.KeyCount)
	s.initStaged(c)

	s.karma = 100
	s.hp = 100

	s.AudioPlayer = game.NewAudioPlayer(s.chart.AbsPath(s.chart.AudioFilename))
	s.AudioPlayer.Play()
	s.AudioPlayer.Pause()
	s.auto = s.chart.GenAutoKeyEvents(instability)
	s.playSE = SEPlayer(cwd)
	s.timeStamp = c.TimeStampFinder()
	s.sceneUI = newSceneUI(p, s.chart.KeyCount)
	s.setNoteSprites()
	s.ready = true
	s.jm = game.NewJudgmentMeter(Judgments[:])

	return s
}

func (s *Scene) Ready() bool { return s.ready }

// 음악이 그림보다 느리다
// 여러번 로딩되면 음악이 그림 속도를 따라가다
func (s *Scene) Update() error {
	var now int64
	if !s.initUpdate {
		ebiten.SetWindowTitle(fmt.Sprintf("gosu - %s [%s]", s.chart.MusicName, s.chart.ChartName))
		s.AudioPlayer.Play()
		s.startTime = time.Now()
		s.initUpdate = true
		return nil
	}
	if !game.AudioContext.IsReady() {
		return nil
	}
	now = time.Since(s.startTime).Milliseconds()
	// if now < 3000 { // unsafe: 꼬로록 소리 남
	//		s.AudioPlayer.Seek(time.Now().Sub(s.startTime))
	//		}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) || now > s.chart.EndTime()+2000 { // temp: 2초 여유 두기
		_ = s.AudioPlayer.Close()
		s.done = true
	}
	ts := s.timeStamp(now)
	cursor := float64(now-ts.Time)*ts.Factor + ts.Position
	for i, n := range s.chart.Notes {
		rp := (n.position-cursor)*s.speed - Settings.HitPosition // relative position
		s.chart.Notes[i].Sprite.Y = int(-rp*(float64(s.ScreenSize.Y)/100) + float64(n.Sprite.H)/2)
		if n.Type == TypeLNTail {
			s.chart.Notes[i].LongSprite.Y = n.Sprite.Y + n.Sprite.H // todo: why?
		}
	}
	// judge: score과 staged도 따라서 업데이트
	s.lastJudge = empty
	for _, e := range s.auto(now) {
		s.judge(e)
	}

	// 따로 처리: lost, scored되고 시간 다 된 LNTail
	lost := func(timeDiff int64) bool { return timeDiff < -Bad.Window } // never hit
	flushable := func(n Note, timeDiff int64) bool { return n.scored && timeDiff < Miss.Window }
	for k, i := range s.staged {
		if i < 0 {
			continue
		}
		n := s.chart.Notes[i]
		timeDiff := n.Time - now

		if lost(timeDiff) {
			s.applyScore(i, Miss)
		}

		if n.Type == TypeLNTail && flushable(n, timeDiff) {
			s.staged[k] = n.next
		}
	}
	return nil
}

func (s *Scene) Draw(screen *ebiten.Image) {
	now := time.Since(s.startTime).Milliseconds()
	screen.DrawImage(s.bg, s.bgop)
	s.playfield.Draw(screen)
	// for i, j := range Judgments {
	// 	if s.lastJudge == j {
	// 		s.judgeSprite[i].Draw(screen)
	// 		break
	// 	}
	// }
	s.judgeSprite[0].Draw(screen)
	for _, n := range s.chart.Notes {
		if n.Type == TypeLNTail {
			n.LongSprite.Draw(screen)
		}
	}
	// to make sure LNs go most behind
	for _, n := range s.chart.Notes {
		n.Sprite.Draw(screen)
	}

	// s.jm.DrawTiming(screen, s.timeDiffs)
	// if len(s.timeDiffs) > 20 {
	// 	s.timeDiffs = s.timeDiffs[:20]
	// }
	// scoreStr := fmt.Sprintf("%.0f", s.score)
	// text.Draw(screen, scoreStr, arcadeFont, s.ScreenSize.X-4*fontSize, fontSize, color.White)
	// comboStr := fmt.Sprintf("%d", s.combo)
	// text.Draw(screen, comboStr, titleArcadeFont, s.ScreenSize.X/2-1.5*fontSize/2, s.ScreenSize.Y/2-fontSize, color.White)
	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		`CurrentFPS: %.2f
CurrentTPS: %.2f
Time: %.3fs

score: %.0f
karma: %.2f
hp: %.2f
combo: %d
judge: %v
`, ebiten.CurrentFPS(), ebiten.CurrentTPS(), float64(now)/1000,
		s.score, s.karma, s.hp, s.combo, s.judgeCounts))
	s.drawCombo(screen)
	s.drawScore(screen)
	for i, tb := range s.lastPressed {
		if tb.Value || now-tb.Time < 90 { // temp
			s.stageKeysPressed[i].Draw(screen)
		} else {
			s.stageKeys[i].Draw(screen)
		}
	}
	s.jm.Sprite.Draw(screen)
}

func (s *Scene) initStaged(c *Chart) {
	s.staged = make([]int, c.KeyCount)
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
}

func (s *Scene) Done(args *game.TransSceneArgs) bool {
	if s.done && args.Next == "" {
		args.Next = "sceneSelect"
		args.Args = nil
	}
	return s.done
}
