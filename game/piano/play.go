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
	sampleBuffer []game.Sample
	Scorer       Scorer
	now          int32
	dynamics     game.Dynamics
	cursor       float64

	keyCount int
	stage    StageOpts
	key      KeysOpts

	field      FieldComp
	bars       BarsComp
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
		s.Scorer.markKeysUntouchedNote(ka.Time, kji)
		kn := s.Scorer.keysFocusNote()
		s.addKeysSampleToBuffer(kn, ka)
		s.Scorer.tryJudge(ka, kji)

		klnh := s.keysLongNoteHolding(kn, ka)
		s.updateTime(ka.Time)
	}
	s.updateTime(now)
	return nil
}

func (s Play) addKeysSampleToBuffer(kn []Note, ka game.KeyboardAction) {
	s.sampleBuffer = s.sampleBuffer[:0]
	for k, n := range kn {
		if n.valid && ka.KeysAction[k] == game.Hit {
			s.sampleBuffer = append(s.sampleBuffer, n.Sample)
		}
	}
}

func (s Play) keysLongNoteHolding(kn []Note, ka game.KeyboardAction) []bool {
	klnh := make([]bool, s.keyCount)
	kh := ka.KeysHolding()
	for k, n := range kn {
		if n.valid && n.Type == Tail && kh[k] {
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
	bs := s.bars.bars
	for i := range bs {
		bs[i].position *= scale
	}
}

func (s Play) Draw(dst draws.Image) {
	s.field.Draw(dst)
	s.bars.Draw(dst)
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
	f(&b, "Score: %.0f \n", s.Scorer.Score)
	f(&b, "Combo: %d\n", s.Scorer.Combo)
	f(&b, "Flow: %.0f/%2d\n", s.Scorer.factors[flow], s.Scorer.maxFactors[flow])
	f(&b, " Acc: %.0f/%2d\n", s.Scorer.factors[acc], s.Scorer.maxFactors[acc])
	f(&b, "Judgment counts: %v\n", s.Scorer.JudgmentCounts)
	f(&b, "\n")
	f(&b, "Speed scale (PageUp/Down): x%.2f (x%.2f)\n", s.dynamics.SpeedScale, s.dynamics.Speed())
	f(&b, "(Exposure time: %dms)\n", s.dynamics.NoteExposureDuration(s.stage.H))
	return b.String()
}

// isKeyPresseds []bool // for keys, key lightings, and hold lightings
// isKeyHolds    []bool // for long note body, hold lightings
