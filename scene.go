package main

type KeyAction int

const (
	Idle KeyAction = iota
	Hit
	Release
	Hold
)

// Todo: ebiten의 IsKeyPressed()를 사용할 것이므로, 얘는 keyboard input 쪽으로 보내야 할듯
// Todo: Move to gosu/input/kb
func (s *ScenePlay) KeyAction(k int) KeyAction {
	return CurrentKeyAction(s.LastPressed[k], s.Pressed[k])
}

func CurrentKeyAction(last, now bool) KeyAction {
	switch {
	case !last && !now:
		return Idle
	case !last && now:
		return Hit
	case last && !now:
		return Release
	case last && now:
		return Hold
	default:
		panic("not reach")
	}
}

// ScenePlay: struct, PlayScene: function
type ScenePlay struct {
	Chart     *Chart
	PlayNotes []*PlayNote

	Pressed     []bool
	LastPressed []bool
	Time        int64
	LastTime    int64

	StagedNotes    []*PlayNote
	Combo          int
	Karma          float64
	KarmaSum       float64 // ScoreA; 700k
	JudgmentCounts []int   // ScoreB, ScoreC
}

func NewScenePlay(c *Chart) *ScenePlay {
	s := new(ScenePlay)
	s.PlayNotes = NewPlayNotes(c) // Todo: add Mods to input param
	return s
}

// PlayNote is for in-game. Handled by pointers to modify its fields easily.
type PlayNote struct {
	Note
	Prev   *PlayNote
	Next   *PlayNote
	Scored bool
}

func NewPlayNotes(c *Chart) []*PlayNote {
	pns := make([]*PlayNote, 0, len(c.Notes))
	prevs := make([]*PlayNote, c.Parameter.KeyCount)
	for _, n := range c.Notes {
		prev := prevs[n.Key]
		next := &PlayNote{
			Note: n,
			Prev: prev,
		}
		if prev != nil { // Set Next value later
			prev.Next = next
		}
		prevs[n.Key] = next
	}
	return pns
}

func (s *ScenePlay) Update() {
	s.FetchInput()
	s.CalcNoteXY()
	s.CheckScore()
	s.PlaySEs()
}

func (s *ScenePlay) Draw() {
	DrawBG()
	DrawField()
	DrawNote()
	DrawCombo()
	DrawJudgment()
	DrawOthers() // Score, judgment counts and other states
}

func (s ScenePlay) PlaySEs() {
	for k, n := range s.StagedNotes {
		if s.KeyAction(k) == Hit {
			n.PlaySE()
		}
	}
}

func (n PlayNote) PlaySE() {}

// TimeStamp도 Scene에서 관리해야 할 것 같다
// type TimeStamp struct {
// 	Time     int64
// 	NextTime int64
// 	Position float64
// 	Factor   float64
// }
