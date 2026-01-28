package piano

import (
	"fmt"
	"strings"

	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
)

// Struct Play goes a part of ScenePlay.
// TODO: ErrorMeter, Timer
type Play struct {
	*Resources
	*Options
	*Chart
	Scorer // This seems okay to be non-pointer as this is created in NewPlay.
	// soundPlayer *audios.SoundPlayer

	// Drawers
	field          Field
	barDrawer      BarDrawer
	hint           Hint
	noteDrawer     NoteDrawer
	keyButtons     KeyButtons
	backlights     Backlights
	hitLights      HitLights
	holdLights     HoldLights
	judgmentDrawer JudgmentDrawer
	comboDrawer    game.ComboDrawer
	scoreDrawer    game.ScoreDrawer
}

func NewPlay(res *Resources, opts *Options, c *Chart, mods Mods, sp *audios.SoundPlayer) (*Play, error) {
	scn := &Play{
		Resources: res,
		Options:   opts,
		Chart:     c,
		// Mods may affect judgment range.
		// Scorer plays a corresponding sample when a key is hit.
		Scorer: NewScorer(c, sp),
		// soundPlayer: sp,

		field:          NewField(res, opts, c.keyCount),
		barDrawer:      NewBarDrawer(res, opts, c),
		hint:           NewHint(res, opts, c.keyCount),
		noteDrawer:     NewNoteDrawer(res, opts, c),
		keyButtons:     NewKeyButtons(res, opts, c.keyCount),
		backlights:     NewBacklights(res, opts, c.keyCount),
		hitLights:      NewHitLights(res, opts, c.keyCount),
		holdLights:     NewHoldLights(res, opts, c),
		judgmentDrawer: NewJudgmentDrawer(res, opts),
		comboDrawer:    game.NewComboDrawer(res.ComboImages, &opts.Combo),
		scoreDrawer:    game.NewScoreDrawer(res.ScoreImages, &opts.Score),
	}
	// TODO: refactor
	scn.SetSpeedScale(1, scn.SpeedScale)
	return scn, nil
}

func (p *Play) Update(now int32, kas []game.KeyboardAction) any {
	for _, ka := range kas {
		// fmt.Printf("ka: %v\n", ka)
		p.Scorer.update(ka)

		cursor := p.Dynamics.Position(ka.Time)
		p.field.Update()
		p.barDrawer.Update(cursor)
		p.hint.Update()
		p.noteDrawer.Update(ka, cursor)
		p.keyButtons.Update(ka)
		p.backlights.Update(ka)
		p.hitLights.Update(p.Scorer.keysJudgmentKind)
		p.holdLights.Update(ka)
		p.judgmentDrawer.Update(p.Scorer.keysJudgmentKind)
		p.comboDrawer.Update(p.Scorer.Combo)
		p.scoreDrawer.Update(p.Scorer.Score)
	}
	return nil
}

// Need to re-calculate positions when Speed has changed.
func (p *Play) SetSpeedScale(oldScale, newScale float64) {
	// oldScale := p.SpeedScale
	ratio := newScale / oldScale
	p.SpeedScale = newScale

	ds := p.Dynamics.Dynamics()
	for i := range ds {
		ds[i].Position *= ratio
	}
	p.Dynamics.SpeedScale = newScale

	ns := p.Chart.notes
	for i := range ns {
		ns[i].position *= ratio
	}
	// for lowermost and uppermost
	p.noteDrawer.scaledScreenSize = game.ScreenSizeY / ratio

	bs := p.Chart.bars
	for i := range bs {
		bs[i].position *= ratio
	}
}

func (p Play) Draw(dst draws.Image) {
	p.field.Draw(dst)
	p.barDrawer.Draw(dst)
	p.hint.Draw(dst)
	p.noteDrawer.Draw(dst)
	p.keyButtons.Draw(dst)
	p.backlights.Draw(dst)
	p.hitLights.Draw(dst)
	p.holdLights.Draw(dst)
	// p.judgmentDrawer.Draw(dst)
	p.comboDrawer.Draw(dst)
	p.scoreDrawer.Draw(dst)
}

func (p Play) NoteExposureDuration() int32 {
	return p.Chart.NoteExposureDuration(p.KeyPositionY)
}

func (p Play) DebugString() string {
	var b strings.Builder
	f := fmt.Fprintf

	// f(&b, "Time: %ds/%ds\n", s.now/1000, s.Span()/1000)
	// f(&b, "\n")
	f(&b, p.Scorer.DebugString())
	f(&b, "Speed scale (PageUp/Down): x%.2f (x%.2f)\n", p.SpeedScale, p.Speed())
	f(&b, "(Exposure time: %dms)\n", p.NoteExposureDuration())
	return b.String()
}
