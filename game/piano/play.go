package piano

import (
	"fmt"
	"strings"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
)

// Options is for passing from game to mode.
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
	dynamics game.Dynamics
	scorer   Scorer

	// keyCount   int
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
		s.scorer.update(ka)
		// cursor := s.dynamics.Cursor(ka.Time)
	}
	return nil
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
	ns := s.notes.notes.notes
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
	f(&b, "Score: %.0f \n", s.scorer.Score)
	f(&b, "Combo: %d\n", s.scorer.Combo)
	f(&b, "Flow: %.0f/%.0f\n", s.scorer.factors[flow], s.scorer.maxFactors[flow])
	f(&b, " Acc: %.0f/%.0f\n", s.scorer.factors[acc], s.scorer.maxFactors[acc])
	f(&b, "Judgment counts: %v\n", s.scorer.JudgmentCounts)
	f(&b, "\n")
	f(&b, "Speed scale (PageUp/Down): x%.2f (x%.2f)\n", s.dynamics.SpeedScale, s.dynamics.Speed())
	// f(&b, "(Exposure time: %dms)\n", s.dynamics.NoteExposureDuration(s.stage.H))
	return b.String()
}

// isKeyPresseds []bool // for keys, key lightings, and hold lightings
// isKeyHolds    []bool // for long note body, hold lightings
