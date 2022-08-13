package gosu

import (
	"fmt"
	"io"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu/parse/osr"
)

const (
	DefaultWaitBefore int64 = int64(-1.8 * 1000)
	DefaultWaitAfter  int64 = 3 * 1000
)

// Todo: tick dependent variables should not be global variable.
var (
	MaxComboCountdown    int = MsecToTick(2000)
	MaxJudgmentCountdown int = MsecToTick(600)
	// The following formula is for make score scroll speed constant regardless of TPS.
	DelayedScorePower float64 = 1 - math.Exp(-math.Log(MaxScore)/(float64(MaxTPS)*0.4))
)

// ScenePlay: struct, PlayScene: function
type ScenePlay struct {
	Chart     *Chart
	PlayNotes []*PlayNote

	MainBPM      float64
	BaseSpeed    float64
	Tick         int
	FetchPressed func() []bool
	LastPressed  []bool
	Pressed      []bool
	StagedNotes  []*PlayNote
	LowestTails  []*PlayNote // For drawing long note efficiently
	*TransPoint

	Play        bool // Whether the scene is for play or not
	MusicFile   io.ReadSeekCloser
	MusicPlayer AudioPlayer
	Skin
	Background            Sprite
	LastJudgment          Judgment
	LastJudgmentCountdown int
	ComboCountdown        int
	DelayedScore          float64
	BarLineTimes          []int64
	LowestBarLineIndex    int

	Combo int
	// NoteWeights is a sum of weight of marked notes.
	// This is also max value of each score sum can get at the time.
	NoteWeights    float64
	MaxNoteWeights float64 // Upper bound of NoteWeights
	Flow           float64
	FlowSum        float64
	AccSum         float64
	KoolSum        float64
	JudgmentCounts []int
}

// Todo: May user change speed during playing
func NewScenePlay(c *Chart, cpath string, rf *osr.Format, play bool) *ScenePlay {
	s := new(ScenePlay)
	s.Chart = c
	s.PlayNotes, s.StagedNotes, s.LowestTails, s.MaxNoteWeights = NewPlayNotes(c) // Todo: add Mods to input param
	s.MainBPM, _, _ = c.BPMs()
	s.BaseSpeed = BaseSpeed // From global variable
	waitBefore := DefaultWaitBefore
	if rf != nil && rf.BufferTime() < waitBefore {
		waitBefore = rf.BufferTime()
	}
	s.Tick = MsecToTick(waitBefore)
	if rf != nil {
		s.FetchPressed = NewReplayListener(rf, s.Chart.KeyCount, waitBefore)
	} else {
		s.FetchPressed = NewListener(KeySettings[s.Chart.KeyCount])
	}
	s.LastPressed = make([]bool, c.KeyCount)
	s.Pressed = make([]bool, c.KeyCount)
	s.TransPoint = s.Chart.TransPoints[0]
	for s.TransPoint.Time == s.TransPoint.Next.Time {
		s.TransPoint = s.TransPoint.Next
	}
	s.Flow = 1
	s.JudgmentCounts = make([]int, 5)

	s.Play = play
	if !s.Play {
		return s
	}
	s.MusicFile, s.MusicPlayer = NewAudioPlayer(c.MusicPath(cpath))
	s.MusicPlayer.SetVolume(Volume)
	s.Skin = SkinMap[c.KeyCount]
	if img := NewImage(s.Chart.BackgroundPath(cpath)); img != nil {
		s.Background = Sprite{
			I:      NewImage(s.Chart.BackgroundPath(cpath)),
			Filter: ebiten.FilterLinear,
		}
		s.Background.SetWidth(screenSizeX)
		s.Background.SetCenterY(screenSizeY / 2)
	} else {
		s.Background = s.DefaultBackground
	}
	s.BarLineTimes = s.Chart.BarLineTimes(waitBefore, DefaultWaitAfter)
	return s
}

// TPS affects only on Update(), not on Draw().
// Todo: Apply other values of TransPoint
// Todo: keep playing music when making SceneResult
func (s *ScenePlay) Update(g *Game) {
	for s.TransPoint.Next != nil && s.TransPoint.Next.Time <= s.Time() {
		s.TransPoint = s.TransPoint.Next
	}
	s.LastPressed = s.Pressed
	s.Pressed = s.FetchPressed()
	var worst Judgment
	for k, n := range s.StagedNotes {
		if n == nil {
			continue
		}
		if s.Play && n.Type != Tail && s.KeyAction(k) == Hit {
			n.PlaySE()
		}
		td := n.Time - s.Time() // Time difference; negative values means late hit
		if n.Marked {
			if n.Type != Tail {
				panic("non-tail note has not flushed")
			}
			if td < Miss.Window { // Keep Tail being staged until nearly ends
				s.StagedNotes[n.Key] = n.Next
			}
			continue
		}

		if j := Verdict(n.Type, s.KeyAction(n.Key), td); j.Window != 0 {
			s.MarkNote(n, j)
			if worst.Window < j.Window {
				worst = j
			}
		}
	}
	if !s.Play {
		s.Tick++
		return
	}
	if worst.Window != 0 {
		s.LastJudgment = worst
		s.LastJudgmentCountdown = MaxJudgmentCountdown
	}
	if s.LastJudgmentCountdown == 0 {
		s.LastJudgment = Judgment{}
	} else {
		s.LastJudgmentCountdown--
	}
	if s.ComboCountdown > 0 {
		s.ComboCountdown--
	}
	s.DelayedScore += DelayedScorePower * (s.Score() - s.DelayedScore)

	if ebiten.IsKeyPressed(ebiten.KeyEscape) || s.IsFinished() {
		s.MusicPlayer.Close()
		s.MusicFile.Close()
		g.Scene = selectScene
		return
	}
	if s.Tick == 0 {
		s.MusicPlayer.Play()
	}
	t := s.BarLineTimes[s.LowestBarLineIndex]
	// Bar line and Hint are anchored at the bottom.
	for s.LowestBarLineIndex < len(s.BarLineTimes)-1 &&
		int(s.Position(t)+NoteHeigth/2) >= screenSizeY {
		s.LowestBarLineIndex++
		t = s.BarLineTimes[s.LowestBarLineIndex]
	}
	s.Tick++
}
func (s ScenePlay) Draw(screen *ebiten.Image) {
	bgop := s.Background.Op()
	bgop.ColorM.ChangeHSV(0, 1, BgDimness)
	screen.DrawImage(s.Background.I, bgop)
	s.FieldSprite.Draw(screen)
	s.HintSprite.Draw(screen)
	s.DrawBarLine(screen)
	s.DrawLongNotes(screen)
	s.DrawNotes(screen)
	s.DrawCombo(screen)
	s.DrawJudgment(screen)
	s.DrawScore(screen)
	var fr, ar, rr float64 = 1, 1, 1
	if s.NoteWeights > 0 {
		fr = s.FlowSum / s.NoteWeights
		ar = s.AccSum / s.NoteWeights
		rr = s.KoolSum / s.NoteWeights
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"CurrentFPS: %.2f\nCurrentTPS: %.2f\nTime: %.3fs/%.0fs\n\n"+
			"Score: %.0f | %.0f \nFlow: %.0f/100\nCombo: %d\n\n"+
			"Flow rate: %.2f%%\nAccuracy: %.2f%%\n(Kool: %.2f%%)\nJudgment counts: %v\n\n"+
			"Speed: %.0f | %.0f\n(Exposure time: %.fms)\n\n",
		ebiten.CurrentFPS(), ebiten.CurrentTPS(), float64(s.Time())/1000, float64(s.Chart.EndTime())/1000,
		s.Score(), s.ScoreBound(), s.Flow*100, s.Combo,
		fr*100, ar*100, rr*100, s.JudgmentCounts,
		s.Speed()*100, s.BaseSpeed*100, ExposureTime(s.Speed())))
}

func (s ScenePlay) DrawBarLine(screen *ebiten.Image) {
	for _, t := range s.BarLineTimes[s.LowestBarLineIndex:] {
		sprite := s.BarLineSprite
		sprite.Y = s.Position(t) + NoteHeigth/2
		if sprite.Y < 0 {
			break
		}
		sprite.Draw(screen)
	}
}

// DrawJudgment draws the same judgment for a while.
func (s ScenePlay) DrawJudgment(screen *ebiten.Image) {
	if s.LastJudgmentCountdown <= 0 {
		return
	}
	var sprite Sprite
	for i, j := range Judgments {
		if j.Window == s.LastJudgment.Window {
			sprite = s.JudgmentSprites[i]
			break
		}
	}
	t := MaxJudgmentCountdown - s.LastJudgmentCountdown
	age := float64(t) / float64(MaxJudgmentCountdown)
	switch {
	case age < 0.1:
		sprite.ApplyScale(sprite.ScaleW() * 1.15 * (1 + age))
	case age >= 0.1 && age < 0.2:
		sprite.ApplyScale(sprite.ScaleW() * 1.15 * (1.2 - age))
	case age > 0.9:
		sprite.ApplyScale(sprite.ScaleW() * (1 - 1.15*(age-0.9)))
	}
	sprite.SetCenterX(screenSizeX / 2)
	sprite.SetCenterY(JudgmentPosition)
	sprite.Draw(screen)
}

// DrawCombo draws each number at constant x regardless of their widths.
// Each number image has different size; The standard width is number 0's.
func (s ScenePlay) DrawCombo(screen *ebiten.Image) {
	var wsum int
	if s.ComboCountdown == 0 {
		return
	}
	vs := make([]int, 0)
	for v := s.Combo; v > 0; v /= 10 {
		vs = append(vs, v%10) // Little endian
		// wsum += int(s.ComboSprites[v%10].W + ComboGap)
		wsum += int(s.ComboSprites[0].W) + int(ComboGap)
	}
	wsum -= int(ComboGap)

	t := MaxJudgmentCountdown - s.LastJudgmentCountdown
	age := float64(t) / float64(MaxJudgmentCountdown)
	x := screenSizeX/2 + float64(wsum)/2 - s.ComboSprites[0].W/2
	for _, v := range vs {
		// x -= s.ComboSprites[v].W + ComboGap
		x -= s.ComboSprites[0].W + ComboGap
		sprite := s.ComboSprites[v]
		// sprite.X = x
		sprite.X = x + (s.ComboSprites[0].W - sprite.W/2)
		sprite.SetCenterY(ComboPosition)
		switch {
		case age < 0.1:
			sprite.Y += 0.85 * age * sprite.H
		case age >= 0.1 && age < 0.2:
			sprite.Y += 0.85 * (0.2 - age) * sprite.H
		}
		sprite.Draw(screen)
	}
}

// DrawScore draws each number at constant x regardless of their widths,.
// same as DrawCombo.
func (s ScenePlay) DrawScore(screen *ebiten.Image) {
	var wsum int
	vs := make([]int, 0)
	for v := int(math.Ceil(s.DelayedScore)); v > 0; v /= 10 {
		vs = append(vs, v%10) // Little endian
		// wsum += int(s.ComboSprites[v%10].W)
		wsum += int(s.ComboSprites[0].W)
	}
	if len(vs) == 0 {
		vs = append(vs, 0) // Little endian
		wsum += int(s.ComboSprites[0].W)
	}
	x := float64(screenSizeX) - s.ScoreSprites[0].W/2
	for _, v := range vs {
		// x -= s.ScoreSprites[v].W
		x -= s.ScoreSprites[0].W
		sprite := s.ScoreSprites[v]
		sprite.X = x + (s.ScoreSprites[0].W - sprite.W/2)
		sprite.Draw(screen)
	}
}

func (s ScenePlay) Time() int64 { // In milliseconds.
	return int64(float64(s.Tick) / float64(MaxTPS) * 1000)
}
func (s ScenePlay) IsFinished() bool {
	return s.Time() > s.Chart.EndTime()+DefaultWaitAfter
}
func TickToMsec(tick int) int64        { return int64(1000 * float64(tick) / float64(MaxTPS)) }
func MsecToTick(msec int64) int        { return int(float64(msec) * float64(MaxTPS) / 1000) }
func (s ScenePlay) BeatRatio() float64 { return s.TransPoint.BPM / s.MainBPM }
func (s ScenePlay) Speed() float64     { return s.BaseSpeed * s.BeatRatio() * s.BeatScale }
