package mania

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu/common"
	"github.com/hndada/gosu/engine/audio"
	"github.com/hndada/gosu/engine/kb"
	"github.com/hndada/gosu/engine/scene"
	"github.com/hndada/gosu/engine/ui"

	_ "image/jpeg"
)

const auto = false

type Scene struct {
	ready      bool // whether scene has been loaded
	close      bool
	startTime  time.Time
	initUpdate bool

	mods        Mods
	chart       *Chart
	keyLayout   []kb.Code
	audioPlayer *audio.Player
	speed       float64

	sceneUI
	bg            ui.FixedSprite
	jm            *common.JudgmentMeter // temp
	timingSprites []ui.Animation        // temp

	score       float64
	karma       float64
	hp          float64
	combo       int32
	judgeCounts [len(Judgments)]int
	staged      []int

	timeStamp func(time int64) common.TimeStamp
	auto      func(int64) []keyEvent
	playSE    func()

	lastPressed []bool
}

func NewScene(c *Chart, mods Mods, cwd string) *Scene {
	s := new(Scene)
	const instability = 0 // 0~100; 0 makes the play 'Perfect'

	s.speed = Settings.GeneralSpeed
	s.mods = mods
	s.chart = c.ApplyMods(s.mods)
	s.chart.ScratchMode = Settings.ScratchMode[c.KeyCount] // only for replay
	s.auto = s.chart.GenAutoKeyEvents(instability)
	{
		dir := filepath.Join(cwd, "skin")
		name := "soft-slidertick.wav"
		sePath := filepath.Join(dir, name)
		s.playSE = audio.NewSEPlayer(sePath, 25) // temp: volume
	}
	s.timeStamp = c.TimeStampFinder()
	s.keyLayout = Settings.KeyLayout[WithScratch(c.KeyCount)]
	{
		path := s.chart.AudioPath()
		if path != "" { // There could be keysound-only music
			s.audioPlayer = audio.NewPlayer(path)
		}
	}
	var img *ebiten.Image
	keyKinds := keyKindsMap[WithScratch(c.KeyCount)]
	for i, n := range c.Notes { // temp: All Note, LNHead, LNTail use skin's 'Note' image
		img = Skin.Note[keyKinds[n.key]]
		s.chart.Notes[i].Sprite.SetImage(img)
	}
	s.lastPressed = make([]bool, c.KeyCount)
	{
		s.staged = make([]int, c.KeyCount)
		for k := range s.staged {
			s.staged[k] = -1
		}
		for k := range s.staged {
			for i, n := range c.Notes {
				if n.key == k {
					s.staged[k] = i
					break
				}
			}
		}
	}

	s.karma = 100
	s.hp = 100

	s.sceneUI = newSceneUI(c.KeyCount)
	s.setNoteSprites()
	s.bg = c.BG(common.Settings.BackgroundDimness)
	// s.jm = common.NewJudgmentMeter(Judgments[:]) // TODO: severely lagged

	s.hpScreen = ebiten.NewImage(common.Settings.ScreenSize.X, common.Settings.ScreenSize.Y)
	s.timingSprites = make([]ui.Animation, 0, len(s.chart.Notes))
	if !auto {
		go kb.Listen()
	}
	s.ready = true
	return s
}

func (s *Scene) Ready() bool { return s.ready }

func (s *Scene) Update() error {
	var now int64
	if !s.initUpdate {
		ebiten.SetWindowTitle(fmt.Sprintf("gosu - %s [%s]", s.chart.MusicName, s.chart.ChartName))
		if s.audioPlayer != nil { // There could be keysound-only chart, aka virtual
			s.audioPlayer.Play()
		}
		startTime := time.Now()
		kb.SetTime(startTime)
		s.startTime = startTime
		s.initUpdate = true
		return nil
	}

	now = time.Since(s.startTime).Milliseconds()
	// if now < 3000 { // unsafe: It sounds sinking noise
	//		s.audioPlayer.Seek(time.Now().Sub(s.startTime))
	//		}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) || now > s.chart.EndTime()+2000 { // temp: 2초 여유 두기
		if s.audioPlayer != nil {
			_ = s.audioPlayer.Close()
		}
		s.close = true
	}
	ts := s.timeStamp(now)
	cursor := float64(now-ts.Time)*ts.Factor + ts.Position
	for i, n := range s.chart.Notes {
		rp := (n.position-cursor)*s.speed - Settings.HitPosition // relative position
		s.chart.Notes[i].Sprite.Y = int(-rp*(float64(common.Settings.ScreenSize.Y)/100) - float64(n.Sprite.H)/2)
		if n.Type == TypeLNTail {
			s.chart.Notes[i].LongSprite.Y = n.Sprite.Y + n.Sprite.H // why?: center of tail sprite ~ center of head sprite
			if s.chart.Notes[i].scored {
				s.chart.Notes[i].LongSprite.Saturation = 0.5
				s.chart.Notes[i].LongSprite.Dimness = 0.3
			} else {
				s.chart.Notes[n.prev].Sprite.Saturation = 1
				s.chart.Notes[n.prev].Sprite.Dimness = 1
			}
		}
	}
	// judge: score과 staged도 따라서 업데이트
	if auto {
		for _, e := range s.auto(now) {
			s.judge(e)
			s.lastPressed[e.Key] = e.Pressed // scored되지 않는 누름에도 업데이트 되어야함
		}
	} else {
		events := kb.Fetch()
		for _, e := range events {
			for k, v := range s.keyLayout {
				if v == e.KeyCode {
					e2 := keyEvent{
						Time:    e.Time,
						KeyCode: e.KeyCode,
						Pressed: e.Pressed,
						Key:     k,
					}
					s.judge(e2)
					s.lastPressed[k] = e.Pressed // scored되지 않는 누름에도 업데이트 되어야함
					continue
				}
			}
		}
	}

	// Handle timed-out LNTail after handling lost and scored
	// LNTail should be kept staged even when a player held off LN in middle, got missed.
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
	s.HPBarMask.H = int(float64(s.HPBarColor.H) * (100 - s.hp) / 100)
	return nil
}

func (s *Scene) Draw(screen *ebiten.Image) {
	now := time.Since(s.startTime).Milliseconds()
	s.bg.Draw(screen)
	s.playfield.Draw(screen)

	for i, pressed := range s.lastPressed {
		if pressed {
			s.Spotlights[i].Draw(screen)
		}
	}
	for _, n := range s.chart.Notes {
		if n.Type == TypeLNTail {
			n.LongSprite.Draw(screen)
		}
	}
	// to make sure LNs go most behind
	for _, n := range s.chart.Notes {
		n.Sprite.Draw(screen)
	}
	for i, pressed := range s.lastPressed {
		if pressed {
			s.stageKeysPressed[i].Draw(screen)
		} else {
			s.stageKeys[i].Draw(screen)
		}
	}
	for _, l := range s.Lighting {
		l.Draw(screen)
	}
	for _, l := range s.LightingLN {
		l.Draw(screen)
	}
	{
		var latest int
		for i, js := range s.judgeSprite {
			if js.BornTime.After(s.judgeSprite[latest].BornTime) {
				latest = i
			}
		}
		s.judgeSprite[latest].Draw(screen)
	}
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
	if s.combo > 0 {
		s.drawCombo(screen)
	}
	s.drawScore(screen)

	// s.HPBar.Draw(screen) // temp: hard to keep consistent HPBar image and HP color
	s.hpScreen.Clear()
	s.HPBarColor.Draw(s.hpScreen)
	s.HPBarMask.Draw(s.hpScreen)
	screen.DrawImage(s.hpScreen, &ebiten.DrawImageOptions{})
	// s.jm.Sprite.Draw(screen)
	// for _, sprite := range s.timingSprites {
	// 	sprite.Draw(screen)
	// }
}

func (s *Scene) Close(args *scene.Args) bool {
	if s.close && args.Next == "" {
		args.Next = "SceneSelect"
		args.Args = nil
	}
	return s.close
}
