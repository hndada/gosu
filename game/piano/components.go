package piano

import (
	"github.com/hndada/gosu/draws"
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

func NewComponents(res *Resources, opts *Options, c *Chart) (cmps Components) {
	cmps.field = NewFieldComponent(res, opts, c.keyCount)
	cmps.bars = NewBarsComponent(res, opts, c)
	cmps.hint = NewHintComponent(res, opts, c.keyCount)
	cmps.notes = NewNotesComponent(res, opts, c)
	cmps.keyButtons = NewKeyButtonsComponent(res, opts, c.keyCount)
	cmps.backlights = NewBacklightsComponent(res, opts, c.keyCount)
	cmps.hitLights = NewHitLightsComponent(res, opts, c.keyCount)
	cmps.holdLights = NewHoldLightsComponent(res, opts, c)
	cmps.judgment = NewJudgmentComponent(res, opts)
	cmps.combo = game.NewComboComponent(res.ComboImages, &opts.Combo)
	cmps.score = game.NewScoreComponent(res.ScoreImages, &opts.Score)
	return
}

func (cmps *Components) Update(ka game.KeyboardAction, dys game.Dynamics, s Scorer) any {
	cursor := dys.Position(ka.Time)
	cmps.field.Update()
	cmps.bars.Update(cursor)
	cmps.hint.Update()
	cmps.notes.Update(ka, cursor)
	cmps.keyButtons.Update(ka)
	// cmps.backlights.Update(ka)
	// cmps.hitLights.Update(s.keysJudgmentKind)
	// fmt.Println("AAAAAAAAAAAAAAAAAAA")
	// cmps.holdLights.Update(ka)
	// fmt.Println("BBBBBBBBBBBBBBBBBBB")
	// cmps.judgment.Update(s.keysJudgmentKind)
	// cmps.combo.Update(s.Combo)
	// cmps.score.Update(s.Score)
	return nil
}

func (cmps Components) Draw(dst draws.Image) {
	cmps.field.Draw(dst)
	cmps.bars.Draw(dst)
	cmps.hint.Draw(dst)
	cmps.notes.Draw(dst)
	cmps.keyButtons.Draw(dst)
	// cmps.backlights.Draw(dst)
	// cmps.hitLights.Draw(dst)
	// // cmps.holdLights.Draw(dst)
	// // cmps.judgment.Draw(dst)
	// cmps.combo.Draw(dst)
	// cmps.score.Draw(dst)
}
