package piano

import (
	"io/fs"

	"github.com/hndada/gosu/game"
)

type Resourcer interface {
	Load(fsys fs.FS)
}
type Resources2 []Resourcer

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
