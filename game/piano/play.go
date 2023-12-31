package piano

import (
	"fmt"
	"strings"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
)

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
	now int32 // current time in millisecond
	Scorer
	game.Dynamics

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

func NewPlay(res Resources, opts Options, format any) (s Play, err error) {
	s.Dynamics, err = game.NewDynamics(format, opts.Stage.H)
	if err != nil {
		err = fmt.Errorf("failed to create dynamics: %w", err)
		return
	}
	return
}

func (s *Play) Update(now int32, kas []game.KeyboardAction) any {
	for _, ka := range kas {
		s.now = ka.Time
		s.Scorer.update(ka)

		cursor := s.Cursor(s.now)
		s.field.Update()
		s.bars.Update(cursor)
		s.hint.Update()
		s.notes.Update(ka, cursor)
		s.keyButtons.Update(ka)
		s.backlights.Update(ka)
		s.hitLights.Update(s.keysJudgmentKind)
		s.holdLights.Update(ka, s.notes.keysFocusNote())
		s.judgment.Update(s.keysJudgmentKind)
		s.combo.Update(s.Combo)
		s.score.Update(s.Score)
	}
	s.now = now
	return nil
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

// Need to re-calculate positions when Speed has changed.
func (s *Play) SetSpeedScale(new float64) {
	old := s.SpeedScale
	scale := new / old
	s.SpeedScale = new

	ds := s.Dynamics.Dynamics()
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

func (s Play) DebugString() string {
	var b strings.Builder
	f := fmt.Fprintf

	f(&b, "Time: %ds/%ds\n", s.now/1000, s.Span()/1000)
	f(&b, "\n")
	f(&b, s.Scorer.DebugString())
	f(&b, "Speed scale (PageUp/Down): x%.2f (x%.2f)\n", s.SpeedScale, s.Speed())
	f(&b, "(Exposure time: %dms)\n", s.NoteExposureDuration())
	return b.String()
}
