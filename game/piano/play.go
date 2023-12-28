package piano

import (
	"fmt"
	"strings"
	"time"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/input"
)

// Todo: game.ErrorMeterComp
// Todo: FlowPoint (kind of HP)
type Play struct {
	Scorer

	dynamics   game.Dynamics
	cursor     float64
	field      FieldComp
	bar        BarComp
	hint       HintComp
	notes      NotesComp
	keyButtons KeyButtonsComp
	backlights BacklightsComp
	hitLights  HitLightsComp
	holdLights HoldLightsComp
	judgment   JudgmentComp
	combo      game.ComboComp
	score      game.ScoreComp

	isKeyPresseds []bool // for keys, key lightings, and hold lightings
	isKeyHolds    []bool // for long note body, hold lightings
	// isJudgeOKs         []bool // for 'hit' lighting
	isLongNoteHoldings []bool // for long note body
}

// Need to re-calculate positions when Speed has changed.
func (s *Play) SetSpeedScale(new float64) {
	old := s.dynamics.SpeedScale
	s.dynamics.SpeedScale = new

	scale := new / old
	s.cursor *= scale
	ds := s.dynamics.Dynamics
	for i := range ds {
		ds[i].Position *= scale
	}
	ns := s.notes.notes
	for i := range ns {
		ns[i].Position *= scale
	}
	bs := s.bar.bars
	for i := range bs {
		bs[i].Position *= scale
	}
}

func (s Play) Draw(dst draws.Image) {
	s.field.Draw(dst)
	s.bar.Draw(dst)
	s.hint.Draw(dst)
	s.notes.Draw(dst)
	s.keyButtons.Draw(dst)
	s.backlights.Draw(dst)
	s.hitLights.Draw(dst)
	s.holdLights.Draw(dst)
	s.judgment.Draw(dst)
	s.combo.Draw(dst)
	s.score.Draw(dst)
}

// Just assigning slice will shallow copy.
// NewXxx returns struct, while LoadXxx doesn't.
func NewPlay(res Resources, opts Options, format any) (s Play, err error) {
	s.dynamics, err = game.NewDynamics(format)
	if err != nil {
		err = fmt.Errorf("failed to create dynamics: %w", err)
		return
	}
	s.cursor = s.Speed() * (2 * game.ScreenH)

	s.isKeyPresseds = make([]bool, s.KeyCount)
	s.isKeyHolds = make([]bool, s.KeyCount)
	// s.isJudgeOKs = make([]bool, s.KeyCount)
	s.isLongNoteHoldings = make([]bool, s.KeyCount)
	// s.kool() is just for placeholder.
	s.worstJudgment = s.kool()
	return
}

// All components in Play use unified time.
func (s *Play) Update(now time.Time, kas []game.KeyboardAction) any {
	for _, ka := range kas {
		missed := s.flushStagedNotes(ka.Time)
		if missed {
			worstJudgment = s.miss()
		}

		s.playSounds(ka)
		js := s.tryJudge(ka)

		// draw
		for k, a := range ka.KeyActions {
			switch a {
			case game.Idle, game.Release:
				s.isKeyPresseds[k] = false
				s.isKeyHolds[k] = false
				s.isLongNoteHoldings[k] = false
			case game.Hit:
				s.isKeyPresseds[k] = true
				s.drawKeyTimers[k].Reset()
				s.drawKeyLightingTimers[k].Reset()
				s.drawHitLightingTimers[k].Reset()
				s.drawHoldLightingTimers[k].Reset()
			case game.Hold:
				s.isKeyPresseds[k] = true
				s.isKeyHolds[k] = true
				isLN := s.stagedNotes[k] != nil && s.stagedNotes[k].Type == Tail
				s.isLongNoteHoldings[k] = isLN
			}
		}

		for k, j := range js {
			// Tail also makes hit lighting on.
			if !j.Is(s.miss()) {
				// s.isJudgeOKs[k] = true
				s.drawHitLightingTimers[k].Reset()
			}
			if worstJudgment.Window < j.Window { // j is worse
				worstJudgment = j
			}
		}

		if !worstJudgment.IsBlank() {
			s.worstJudgment = worstJudgment
			s.drawJudgmentTimer.Reset()
		}
	}
	s.cursor = s.dynamics.Cursor(now)
	return nil
}

// No need to check whether staged note is Tail or not,
// since Tail has no sample in advance.

// Todo: set all sample volumes in advance?
func (s Play) playSounds(ka input.KeyboardAction) {
	for k, n := range s.stagedNotes {
		if n == nil {
			continue
		}
		a := ka.KeyActions[k]
		if a != input.Hit {
			continue
		}

		sample := game.DefaultSample
		if n != nil {
			sample = n.Sample
		}

		vol := sample.Volume
		if vol == 0 {
			vol = s.Dynamic.Volume
		}
		scale := *s.SoundVolume
		s.SoundMap.Play(sample.Filename, vol*scale)
	}
}

func (s Play) DebugString() string {
	var b strings.Builder
	f := fmt.Fprintf

	f(&b, "Time: %.3fs/%.0fs\n", s.Now(), game.ToSecond(s.Duration()))
	f(&b, "\n")
	f(&b, "Score: %.0f \n", s.Score)
	f(&b, "Combo: %d\n", s.Combo)
	f(&b, "Flow: %.0f/%2d\n", s.flow, maxFlow)
	f(&b, " Acc: %.0f/%2d\n", s.acc, maxAcc)
	f(&b, "Judgment counts: %v\n", s.JudgmentCounts)
	f(&b, "\n")
	f(&b, "Speed scale (PageUp/Down): x%.2f (x%.2f)\n", s.SpeedScale, s.Speed())
	f(&b, "(Exposure time: %dms)\n", s.NoteExposureDuration(s.Speed()))
	return b.String()
}
