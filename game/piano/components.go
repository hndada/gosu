package piano

import (
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
)

type Components struct {
	field         FieldComponent
	bars          BarsComponent
	hint          HintComponent
	keysNotes     KeysNotesComponent
	keysButton    KeysButtonComponent
	keysBacklight KeysBacklightComponent
	keysHitLight  KeysHitLightComponent
	keysHoldLight KeysHoldLightComponent
	judgment      JudgmentComponent
	combo         game.ComboComponent
	score         game.ScoreComponent
}

func (cmps Components) Draw(dst draws.Image) {
	cmps.field.Draw(dst)
	cmps.bars.Draw(dst)
	cmps.hint.Draw(dst)
	cmps.keysNotes.Draw(dst)
	cmps.keysButton.Draw(dst)
	cmps.keysBacklight.Draw(dst)
	cmps.keysHitLight.Draw(dst)
	cmps.keysHoldLight.Draw(dst)
	cmps.judgment.Draw(dst)
	cmps.combo.Draw(dst)
	cmps.score.Draw(dst)
}
