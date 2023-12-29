package piano

import (
	"fmt"
	"strings"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
)

// Todo: game.ErrorMeterComp
// Todo: FlowPoint (kind of HP)
type Play struct {
	Scorer
	now      int32
	dynamics game.Dynamics
	cursor   float64
	stage    KeysOpts // Todo: KeysOpts -> StageOpts

	field      FieldComp
	bar        BarComp
	hint       HintComp
	notes      NotesComp
	keyButtons KeyButtonsComp
	backlights BacklightsComp
	hitLights  HitLightsComp
	holdLights HoldLightsComp
	judgment   JudgmentComp
	combo      game.ComboComp
	score      game.ScoreComp
}

// Just assigning slice will shallow copy.
// NewXxx returns struct, while LoadXxx doesn't.
func NewPlay(res Resources, opts Options, format any) (s Play, err error) {
	s.dynamics, err = game.NewDynamics(format)
	if err != nil {
		err = fmt.Errorf("failed to create dynamics: %w", err)
		return
	}
	return
}

// All components in Play use unified time.
func (s *Play) Update(now int32, kas []game.KeyboardAction) any {
	for _, ka := range kas {
		kji := make([]int, s.keyCount)
		s.updateKeysFocusNoteIndex(ka.Time, kji)
		s.addKeysSampleToBuffer(ka)
		s.tryJudge(ka, kji)

		kn, kok := s.keysFocusNote()
		klnh := s.keysLongNoteHolding(kn, kok, ka)
		s.updateTime(ka.Time)
	}
	s.updateTime(now)
	return nil
}

func (s Play) keysLongNoteHolding(kn []Note, kok []bool, ka game.KeyboardAction) []bool {
	klnh := make([]bool, s.keyCount)
	for k, hold := range ka.KeysHolding() {
		if hold && kok[k] && kn[k].Type == Tail {
			klnh[k] = true
		}
	}
	return klnh
}

func (s *Play) updateTime(now int32) {
	s.now = now
	s.cursor = s.dynamics.Cursor(now)
}

// Need to re-calculate positions when Speed has changed.
func (s *Play) SetSpeedScale(new float64) {
	old := s.dynamics.SpeedScale
	s.dynamics.SpeedScale = new
	scale := new / old

	s.cursor *= scale
	ds := s.dynamics.Dynamics
	for i := range ds {
		ds[i].Position *= scale
	}
	ns := s.notes.notes
	for i := range ns {
		ns[i].position *= scale
	}
	bs := s.bar.bars
	for i := range bs {
		bs[i].position *= scale
	}
}

func (s Play) Draw(dst draws.Image) {
	s.field.Draw(dst)
	s.bar.Draw(dst)
	s.hint.Draw(dst)
	s.notes.Draw(dst)
	s.keyButtons.Draw(dst)
	s.backlights.Draw(dst)
	s.hitLights.Draw(dst)
	s.holdLights.Draw(dst)
	s.judgment.Draw(dst)
	s.combo.Draw(dst)
	s.score.Draw(dst)
}

func (s Play) DebugString() string {
	var b strings.Builder
	f := fmt.Fprintf

	f(&b, "Time: %.3fs/%.0fs\n", s.now, s.notes.Span()/1000)
	f(&b, "\n")
	f(&b, "Score: %.0f \n", s.Score)
	f(&b, "Combo: %d\n", s.Combo)
	f(&b, "Flow: %.0f/%2d\n", s.Scorer.factors[flow], s.Scorer.maxFactors[flow])
	f(&b, " Acc: %.0f/%2d\n", s.Scorer.factors[acc], s.Scorer.maxFactors[acc])
	f(&b, "Judgment counts: %v\n", s.JudgmentCounts)
	f(&b, "\n")
	f(&b, "Speed scale (PageUp/Down): x%.2f (x%.2f)\n", s.dynamics.SpeedScale, s.dynamics.Speed())
	f(&b, "(Exposure time: %dms)\n", s.dynamics.NoteExposureDuration(s.stage.BaselineY))
	return b.String()
}

// isKeyPresseds []bool // for keys, key lightings, and hold lightings
// isKeyHolds    []bool // for long note body, hold lightings
