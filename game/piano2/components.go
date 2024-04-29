package piano

import (
	draws "github.com/hndada/gosu/draws5"
	"github.com/hndada/gosu/game"
)

// Too dedicated structs harms readability.
// Resources, Options, and other arguments, explicity.

// Todo: game.ErrorMeterComponent
// Todo: game.TimerComponent
type Components struct {
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

func NewComponents(res *Resources, opts *Options, c *Chart) (cmps *Components) {
	cmps.field = NewFieldComponent(res.field, opts.Field)
	cmps.bars = NewBarsComponent(res.bars, opts.Bars, dys)
	cmps.hint = NewHintComponent(res.hint, opts.Hint)
	cmps.notes = NewNotesComponent(res.notes, opts.Notes, ns, dys)
	cmps.keyButtons = NewKeyButtonsComponent(res.keyButtons, opts.KeyButtons)
	cmps.backlights = NewBacklightsComponent(res.backlights, opts.Backlights)
	cmps.hitLights = NewHitLightsComponent(res.hitLights, opts.HitLights)
	cmps.holdLights = NewHoldLightsComponent(res.holdLights, opts.HoldLights)
	cmps.judgment = NewJudgmentComponent(res.judgment, opts.Judgment)
	cmps.combo = game.NewComboComponent(res.combo, opts.Combo)
	cmps.score = game.NewScoreComponent(res.score, opts.Score)
	return
}

func (cmps *Components) Update(ka game.KeyboardAction, dys game.Dynamics, s Scorer) any {
	cursor := dys.Position(ka.Time)
	cmps.field.Update()
	cmps.bars.Update(cursor)
	cmps.hint.Update()
	cmps.notes.Update(ka, cursor)
	cmps.keyButtons.Update(ka)
	cmps.backlights.Update(ka)
	cmps.hitLights.Update(s.keysJudgmentKind)
	cmps.holdLights.Update(ka, s.keysFocusNote())
	cmps.judgment.Update(s.keysJudgmentKind)
	cmps.combo.Update(s.Combo)
	cmps.score.Update(s.Score)
	return nil
}

func (cmps Components) Draw(dst draws.Image) {
	cmps.field.Draw(dst)
	cmps.bars.Draw(dst)
	cmps.hint.Draw(dst)
	cmps.notes.Draw(dst)
	cmps.keyButtons.Draw(dst)
	cmps.backlights.Draw(dst)
	cmps.hitLights.Draw(dst)
	cmps.holdLights.Draw(dst)
	cmps.judgment.Draw(dst)
	cmps.combo.Draw(dst)
	cmps.score.Draw(dst)
}
