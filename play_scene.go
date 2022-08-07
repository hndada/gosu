package main

import (
	"fmt"
	"io"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

// ScenePlay: struct, PlayScene: function
type ScenePlay struct {
	Tick int

	Chart       *Chart
	PlayNotes   []*PlayNote
	KeySettings []ebiten.Key // Todo: separate ebiten

	Pressed     []bool
	LastPressed []bool

	StagedNotes    []*PlayNote
	Combo          int
	Karma          float64
	KarmaSum       float64
	JudgmentCounts []int

	// In dev
	ReplayMode   bool
	ReplayStates []ReplayState
	ReplayCursor int

	Speed float64
	TransPoint
	NoteSprites []Sprite
	BodySprites []Sprite
	HeadSprites []Sprite
	TailSprites []Sprite

	MusicFile   io.ReadSeekCloser
	MusicPlayer *audio.Player
}

func TickToMsec(tick int) int64 { return int64(1000 * float64(tick) / float64(MaxTPS)) }
func MsecToTick(msec int64) int { return int(msec) * MaxTPS }
func NewScenePlay(c *Chart, cpath string) *ScenePlay {
	s := new(ScenePlay)
	s.Tick = -2 * MaxTPS // Put 2 seconds of waiting
	s.Chart = c
	s.PlayNotes, s.StagedNotes = NewPlayNotes(c) // Todo: add Mods to input param
	// s.KeySettings = KeySettings[s.Chart.Parameter.KeyCount]
	s.JudgmentCounts = make([]int, 5)
	s.LastPressed = make([]bool, c.KeyCount)
	s.Pressed = make([]bool, c.KeyCount)
	s.Karma = 1
	s.TransPoint = TransPoint{
		s.Chart.SpeedFactors[0],
		s.Chart.Tempos[0],
		s.Chart.Volumes[0],
		s.Chart.Effects[0],
	}
	s.NoteSprites = make([]Sprite, c.KeyCount)
	s.BodySprites = make([]Sprite, c.KeyCount)
	var wsum int
	for k, kind := range NoteKindsMap[c.KeyCount] {
		w := int(NoteWidths[c.KeyCount][int(kind)] * Scale()) // w should be integer, since it is a play note's width.
		var fpath string
		fpath = "skin/note/" + fmt.Sprintf("n%d.png", []int{1, 2, 3, 3}) // Todo: 4th note image
		s.NoteSprites[k] = Sprite{
			I: NewImage(fpath),
			W: float64(w),
			H: NoteHeigth * Scale(),
		}
		fpath = "skin/note/" + fmt.Sprintf("l%d.png", []int{1, 2, 3, 3}) // Todo: 4th note image
		s.BodySprites[k] = Sprite{
			I: NewImage(fpath),
			W: float64(w),
			H: NoteHeigth * Scale(), // Long note body does not have to be scaled though.
		}
		wsum += w
	}
	s.HeadSprites = make([]Sprite, c.KeyCount)
	s.TailSprites = make([]Sprite, c.KeyCount)
	copy(s.HeadSprites, s.NoteSprites)
	copy(s.TailSprites, s.NoteSprites)

	// Todo: Scratch should be excluded to width sum.
	x := (ScreenSizeX - wsum) / 2 // x should be integer as well as w
	for k, kind := range NoteKindsMap[c.KeyCount] {
		s.NoteSprites[k].X = float64(x)
		x += int(NoteWidths[c.KeyCount][kind] * Scale())
	}
	f, err := os.Open(c.MusicPath(cpath))
	if err != nil {
		panic(err)
	}
	s.MusicPlayer, err = Context.NewPlayer(f)
	if err != nil {
		panic(err)
	}
	return s
}

func (s *ScenePlay) Update() {
	s.Tick++
	if s.IsFinished() {
		s.MusicFile.Close()
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

	for k, p := range s.Pressed {
		s.LastPressed[k] = p
		if s.ReplayMode {
			for s.ReplayCursor < len(s.ReplayStates)-1 && s.Time() > s.ReplayStates[s.ReplayCursor].Time {
				s.ReplayCursor++
			}
			s.ReplayCursor--
			if s.ReplayCursor < 0 {
				s.ReplayCursor = 0
			}
			s.Pressed = s.ReplayStates[s.ReplayCursor].Pressed
		} else {
			s.Pressed[k] = ebiten.IsKeyPressed(s.KeySettings[k])
		}
	}
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
		if j := Verdict(n.Type, s.KeyAction(n.Key), td); j.Window != 0 {
			s.Score(n, j)
		}
	}
}
func (s ScenePlay) Time() int64 {
	return int64(float64(s.Tick) / float64(MaxTPS) * 1000)
}

func (s ScenePlay) IsFinished() bool {
	return s.Time() > 3000+s.PlayNotes[len(s.PlayNotes)-1].Time
}

func (s *ScenePlay) Draw(screen *ebiten.Image) {
	if s.IsFinished() {
		s.DrawClear()
		return
	}
	s.DrawBG()
	s.DrawField()
	s.DrawNotes(screen)
	s.DrawCombo()
	s.DrawJudgment()
	s.DrawOthers() // Score, judgment counts and other states
}

func (s ScenePlay) DrawBG()       {}
func (s ScenePlay) DrawField()    {}
func (s ScenePlay) DrawCombo()    {}
func (s ScenePlay) DrawJudgment() {}
func (s ScenePlay) DrawScore()    {}
func (s ScenePlay) DrawClear()    {}
func (s ScenePlay) DrawOthers()   {} // judgment counts and scene's state
