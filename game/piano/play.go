package piano

import (
	"fmt"
	"strings"

	"github.com/hndada/gosu/game"
)

type Play struct {
	game.Dynamics
	Scorer
	Components
}

func NewPlay(res Resources, opts Options, chart any, mods Mods) (Play, error) {
	dys, err := game.NewDynamics(chart, opts.Stage.H)
	if err != nil {
		return Play{}, fmt.Errorf("failed to create dynamics: %w", err)
	}
	ns := NewNotes(chart, dys)

	return Play{
		Dynamics:   dys,
		Scorer:     NewScorer(ns, mods),
		Components: NewComponents(res, opts, dys, ns),
	}, nil
}

func (p *Play) Update(now int32, kas []game.KeyboardAction) any {
	for _, ka := range kas {
		p.Scorer.update(ka)
		p.Components.Update(ka, p.Dynamics, p.Scorer)
	}
	return nil
}

// Need to re-calculate positions when Speed has changed.
func (p *Play) SetSpeedScale(newScale float64) {
	oldScale := p.SpeedScale
	scale := newScale / oldScale
	p.SpeedScale = newScale

	ds := p.Dynamics.Dynamics()
	for i := range ds {
		ds[i].Position *= scale
	}
	ns := p.notes.notes
	for i := range ns {
		ns[i].position *= scale
	}
	bs := p.bars.bars
	for i := range bs {
		bs[i].position *= scale
	}
}

func (p Play) DebugString() string {
	var b strings.Builder
	f := fmt.Fprintf

	// f(&b, "Time: %ds/%ds\n", s.now/1000, s.Span()/1000)
	f(&b, "\n")
	f(&b, p.Scorer.DebugString())
	f(&b, "Speed scale (PageUp/Down): x%.2f (x%.2f)\n", p.SpeedScale, p.Speed())
	f(&b, "(Exposure time: %dms)\n", p.NoteExposureDuration())
	return b.String()
}
