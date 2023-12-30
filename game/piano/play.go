package piano

import (
	"fmt"
	"strings"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
)

// Objective: manage UI cmponents with each own struct.
// Options is for passing from game to mode.
// A field is plural form if is drawn per key.

// No XxxArgs. It just makes the code too verbose.
// Introducing interface as a field would make the code too verbose.
type Options struct {
	KeyCount int
	Stage    StageOptions
	Key      KeysOptions

	Field      FieldOptions
	Bars       BarsOptions
	Hint       HintOptions
	Notes      NotesOptions
	KeyButtons KeyButtonsOptions
	Backlights BacklightsOptions
	HitLights  HitLightsOptions
	HoldLights HoldLightsOptions
	Judgment   JudgmentOptions
	Combo      game.ComboOptions
	Score      game.ScoreOptions
}

// Todo: game.ErrorMeterComponent
// Todo: FlowPoint (kind of HP)
type Play struct {
	sampleBuffer []game.Sample
	dynamics     game.Dynamics
	Scorer       Scorer

	keyCount   int
	field      FieldComponent
	bars       BarsComponent
	hint       HintComponent
	notes      NotesComponent
	keyButtons KeyButtonsComponent
	backlights BacklightsComponent
	hitLights  HitLightsComponent
	holdLights HoldLightsComponent
	judgment   JudgmentComponent
	combo      game.ComboComponent
	score      game.ScoreComponent
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
	s.sampleBuffer = s.sampleBuffer[:0]
	for _, ka := range kas {
		kji := make([]int, s.keyCount)
		s.Scorer.markKeysUntouchedNote(ka.Time, kji)
		kn := s.Scorer.keysFocusNote()

		s.addKeysSampleToBuffer(kn, ka)
		s.Scorer.tryJudge(ka, kji)
		cursor := s.dynamics.Cursor(now)
	}
	return nil
}

func (s Play) addKeysSampleToBuffer(kn []Note, ka game.KeyboardAction) {
	for k, n := range kn {
		if n.valid && ka.KeysAction[k] == game.Hit {
			s.sampleBuffer = append(s.sampleBuffer, n.Sample)
		}
	}
}

// Need to re-calculate positions when Speed has changed.
func (s *Play) SetSpeedScale(new float64) {
	old := s.dynamics.SpeedScale
	s.dynamics.SpeedScale = new
	scale := new / old

	// s.cursor *= scale
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

func (p Play) Draw(dst draws.Image) {
	p.field.Draw(dst)
	p.bars.Draw(dst)
	p.hint.Draw(dst)
	p.notes.Draw(dst)
	p.keyButtons.Draw(dst)
	p.backlights.Draw(dst)
	p.hitLights.Draw(dst)
	p.holdLights.Draw(dst)
	p.judgment.Draw(dst)
	p.combo.Draw(dst)
	p.score.Draw(dst)
}

func (s Play) DebugString() string {
	var b strings.Builder
	f := fmt.Fprintf

	// f(&b, "Time: %.3fs/%.0fs\n", s.now, s.notes.Span()/1000)
	f(&b, "\n")
	f(&b, "Score: %.0f \n", s.Scorer.Score)
	f(&b, "Combo: %d\n", s.Scorer.Combo)
	f(&b, "Flow: %.0f/%2d\n", s.Scorer.factors[flow], s.Scorer.maxFactors[flow])
	f(&b, " Acc: %.0f/%2d\n", s.Scorer.factors[acc], s.Scorer.maxFactors[acc])
	f(&b, "Judgment counts: %v\n", s.Scorer.JudgmentCounts)
	f(&b, "\n")
	f(&b, "Speed scale (PageUp/Down): x%.2f (x%.2f)\n", s.dynamics.SpeedScale, s.dynamics.Speed())
	// f(&b, "(Exposure time: %dms)\n", s.dynamics.NoteExposureDuration(s.stage.H))
	return b.String()
}

// isKeyPresseds []bool // for keys, key lightings, and hold lightings
// isKeyHolds    []bool // for long note body, hold lightings
