package gosu

import (
	"fmt"
	"io"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu/parse/osr"
)

// ScenePlay: struct, PlayScene: function
type ScenePlay struct {
	Chart     *Chart
	PlayNotes []*PlayNote

	MusicFile   io.ReadSeekCloser
	MusicPlayer *audio.Player

	Skin
	Background        Sprite
	Judgment          Judgment
	JudgmentCountdown int

	Speed        float64
	Tick         int
	FetchPressed func() []bool
	LastPressed  []bool
	Pressed      []bool
	TransPoint
	StagedNotes []*PlayNote

	Combo          int
	Karma          float64
	KarmaSum       float64
	JudgmentCounts []int
}

func NewScenePlay(c *Chart, cpath string, rf *osr.Format) *ScenePlay {
	s := new(ScenePlay)
	s.Chart = c
	s.PlayNotes, s.StagedNotes = NewPlayNotes(c) // Todo: add Mods to input param

	var err error
	s.MusicFile, err = os.Open(c.MusicPath(cpath))
	if err != nil {
		panic(err)
	}
	s.MusicPlayer, err = Context.NewPlayer(s.MusicFile)
	if err != nil {
		panic(err)
	}

	s.Skin = SkinMap[c.KeyCount]
	s.Background = Sprite{
		I: NewImage(c.BgPath(cpath)),
		W: screenSizeX,
		H: screenSizeY,
	}

	s.Speed = Speed      // From global variable.
	s.Tick = -2 * MaxTPS // Put 2 seconds of waiting
	if rf != nil {
		s.FetchPressed = NewReplayListener(rf, s.Chart.KeyCount)
	} else {
		s.FetchPressed = NewListener(KeySettings[s.Chart.KeyCount])
	}
	s.LastPressed = make([]bool, c.KeyCount)
	s.Pressed = make([]bool, c.KeyCount)
	s.TransPoint = TransPoint{
		s.Chart.SpeedFactors[0],
		s.Chart.Tempos[0],
		s.Chart.Volumes[0],
		s.Chart.Effects[0],
	}

	s.Karma = 1
	s.JudgmentCounts = make([]int, 5)
	return s
}

// TPS affects only on Update(), not on Draw().
func (s *ScenePlay) Update(g *Game) {
	if s.IsFinished() {
		if s.MusicPlayer != nil {
			s.MusicFile.Close()
			// Music still plays even the chart is finished.
		}
		g.Scene = selectScene
		return
	}
	if s.Tick == 0 {
		s.MusicPlayer.Play()
	}
	for s.Time() < s.SpeedFactor.Next.Time {
		s.SpeedFactor = s.SpeedFactor.Next
	}
	for s.Time() < s.Tempo.Next.Time {
		s.Tempo = s.Tempo.Next
	}
	for s.Time() < s.Volume.Next.Time {
		s.Volume = s.Volume.Next
	}
	for s.Time() < s.Effect.Next.Time {
		s.Effect = s.Effect.Next
	}
	s.Pressed = s.FetchPressed()
	for k, n := range s.StagedNotes {
		if n == nil {
			continue
		}
		if n.Type != Tail && s.KeyAction(k) == Hit {
			n.PlaySE()
		}
		td := n.Time - s.Time() // Time difference; negative values means late hit
		if n.Scored {
			if n.Type != Tail {
				panic("non-tail note has not flushed")
			}
			if td < Miss.Window { // Keep Tail being staged until nearly ends
				s.StagedNotes[n.Key] = n.Next
			}
			continue
		}
		var worst Judgment
		if j := Verdict(n.Type, s.KeyAction(n.Key), td); j.Window != 0 {
			s.Score(n, j)
			if worst.Window < j.Window {
				worst = j
			}
		}
		if s.Judgment.Window < worst.Window {
			s.Judgment = worst
			s.JudgmentCountdown = MsecToTick(1000)
		} else {
			s.JudgmentCountdown--
		}
	}
	s.Tick++
}
func (s *ScenePlay) Draw(screen *ebiten.Image) {
	bgop := s.Background.Op()
	bgop.ColorM.ChangeHSV(0, 1, BgDimness)
	screen.DrawImage(s.Background.I, bgop)

	s.FieldSprite.Draw(screen)
	if s.IsFinished() {
		s.ClearSprite.Draw(screen)
	} else {
		s.DrawNotes(screen)
		if s.Combo > 0 {
			s.DrawCombo(screen)
		}
		if s.JudgmentCountdown > 0 { // Draw the same judgment for a while.
			for i, j := range Judgments {
				if j.Window == s.Judgment.Window {
					s.JudgmentSprites[i].Draw(screen)
					break
				}
			}
		}
		s.DrawScore(screen)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"CurrentFPS: %.2f\nCurrentTPS: %.2f\nTime: %.3fs\n"+
			"Score: %.0f\nKarma: %.2f\nCombo: %d\n"+
			"judge: %v", ebiten.CurrentFPS(), ebiten.CurrentTPS(), float64(s.Time())/1000,
		s.CurrentScore(), s.Karma, s.Combo,
		s.JudgmentCounts))
}

// DrawCombo supposes each number image has different size.
// Wait, we loaded number image with adjusting size.
func (s *ScenePlay) DrawCombo(screen *ebiten.Image) {
	var wsum int
	vs := make([]int, 0)
	for v := s.Combo; v > 0; v /= 10 {
		vs = append(vs, v%10) // Little endian
		wsum += int(s.ComboSprites[v%10].W - ComboGap)
	}
	wsum += int(ComboGap)
	x := (screenSizeX + float64(wsum)) / 2
	for _, v := range vs {
		x -= s.ComboSprites[v].W + ComboGap
		sprite := s.ComboSprites[v]
		sprite.X = x
		sprite.Y = ComboPosition - sprite.H/2
		sprite.Draw(screen)
	}
}

func (s *ScenePlay) DrawScore(screen *ebiten.Image) {
	var wsum int
	vs := make([]int, 0)
	for v := int(s.CurrentScore()); v > 0; v /= 10 {
		vs = append(vs, v%10) // Little endian
		wsum += int(s.ComboSprites[v%10].W)
	}
	x := float64(screenSizeX)
	for _, v := range vs {
		x -= s.ScoreSprites[v].W
		sprite := s.ScoreSprites[v]
		sprite.X = x
		sprite.Draw(screen)
	}
}

func (s ScenePlay) Time() int64 { // In milliseconds.
	return int64(float64(s.Tick) / float64(MaxTPS) * 1000)
}
func (s ScenePlay) IsFinished() bool {
	const buffer = 3000
	return s.Time() > buffer+s.PlayNotes[len(s.PlayNotes)-1].Time
}
func TickToMsec(tick int) int64 { return int64(1000 * float64(tick) / float64(MaxTPS)) }
func MsecToTick(msec int64) int { return int(msec) * MaxTPS }
