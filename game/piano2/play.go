package piano

import (
	"fmt"
	"io/fs"
	"strings"

	draws "github.com/hndada/gosu/draws6"
	"github.com/hndada/gosu/game"
)

// Struct Play goes a part of ScenePlay.
type Play struct {
	*Resources
	*Options
	*Chart
	Scorer
	Components
}

func NewPlay(res *Resources, opts *Options, fsys fs.FS, name string, mods Mods) (*Play, error) {
	// dys, err := game.NewDynamics(chart, opts.Stage.H)
	// if err != nil {
	// 	return Play{}, fmt.Errorf("failed to create dynamics: %w", err)
	// }
	// ns := NewNotes(chart, dys)

	c, err := NewChart(fsys, name, mods)
	if err != nil {
		return nil, fmt.Errorf("failed to create chart: %w", err)
	}

	return &Play{
		Chart: c,
		// Mods may affect judgment range.
		Scorer:     NewScorer(c.Notes, mods),
		Components: NewComponents(res, opts, c),
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
	ns := p.notes.notes.data
	for i := range ns {
		ns[i].position *= scale
	}
	bs := p.bars.bars.data
	for i := range bs {
		bs[i].position *= scale
	}
}

func (p Play) Draw(dst draws.Image) {
	p.Components.Draw(dst)
}

func (p Play) NoteExposureDuration() int32 {
	return p.Chart.NoteExposureDuration(p.KeyPositionY)
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
