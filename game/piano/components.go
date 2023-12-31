package piano

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
)

type Resources struct {
	field      FieldResources
	bars       BarsResources
	hint       HintResources
	notes      NotesResources
	keyButtons KeyButtonsResources
	backlights BacklightsResources
	hitLights  HitLightsResources
	holdLights HoldLightsResources
	judgment   JudgmentResources
	combo      game.ComboResources
	score      game.ScoreResources
}

func NewResources(fsys fs.FS) (res Resources) {
	res.field.Load(fsys)
	res.bars.Load(fsys)
	res.hint.Load(fsys)
	res.notes.Load(fsys)
	res.keyButtons.Load(fsys)
	res.backlights.Load(fsys)
	res.hitLights.Load(fsys)
	res.holdLights.Load(fsys)
	res.judgment.Load(fsys)
	res.combo.Load(fsys)
	res.score.Load(fsys)
	return
}

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

func NewOptions(keyCount int) (opts Options) {
	// const defaultKeyCount = 4
	stage := NewStageOptions(keyCount)
	keys := NewKeysOptions(stage)
	return Options{
		KeyCount: keyCount,
		Stage:    stage,
		Key:      keys,

		Field:      NewFieldOptions(stage),
		Bars:       NewBarsOptions(stage),
		Hint:       NewHintOptions(stage),
		Notes:      NewNotesOptions(stage, keys),
		KeyButtons: NewKeyButtonsOptions(keys),
		Backlights: NewBacklightsOptions(keys),
		HitLights:  NewHitLightsOptions(keys),
		HoldLights: NewHoldLightsOptions(keys),
		Judgment:   NewJudgmentOptions(stage),
		Combo:      NewComboOptions(stage),
		Score:      game.NewScoreOptions(),
	}
}

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

func NewComponents(res Resources, opts Options, dys game.Dynamics, ns Notes) (cmps Components) {
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
