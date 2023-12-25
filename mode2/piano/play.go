package piano

import (
	"fmt"
	"io/fs"
	"strings"
	"time"

	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/input"
	mode "github.com/hndada/gosu/mode2"
)

// Todo: mode.ErrorMeterComp
type Scene struct {
	mode.Scene
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
	combo      mode.ComboComp
	score      mode.ScoreComp

	// Todo: Mods, FlowPoint (kind of HP)
	judge         Judge
	isKeyPresseds []bool // for keys, key lightings, and hold lightings
	isKeyHolds    []bool // for long note body, hold lightings
	// isJudgeOKs         []bool // for 'hit' lighting
	isLongNoteHoldings []bool // for long note body
}

func (s Scene) Draw(dst draws.Image) {
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
func NewScene(res Resources, opts Options) (s Scene, err error) {
	const wait = 1100 * time.Millisecond
	s.Timer = mode.NewTimer(*s.MusicOffset, wait)

	if replay != nil {
		s.KeyboardReader = replay.KeyboardReader(s.KeyCount)
	} else {
		keys := input.NamesToKeys(s.KeySettings[s.KeyCount])
		s.Keyboard = input.NewKeyboard(keys, s.StartTime())
		defer s.Keyboard.Listen()
	}

	const ratio = 1
	s.musicPlayer, _ = audios.NewMusicPlayerFromFile(fsys, s.MusicFilename, ratio)
	s.SetMusicVolume(*s.MusicVolume)

	s.SoundMap = audios.NewSoundMap(fsys, s.DefaultHitSoundFormat, s.SoundVolume)
	// It is possible for empty string to be a key of a map.
	// https://go.dev/play/p/nn-peGAjawW
	s.SoundMap.AddSound("", s.DefaultHitSoundStreamer)

	s.cursor = float64(s.now) * s.SpeedScale

	s.isKeyPresseds = make([]bool, s.KeyCount)
	s.isKeyHolds = make([]bool, s.KeyCount)
	// s.isJudgeOKs = make([]bool, s.KeyCount)
	s.isLongNoteHoldings = make([]bool, s.KeyCount)
	// s.kool() is just for placeholder.
	s.worstJudgment = s.kool()

	return
}

// Need to re-calculate positions when Speed has changed.
func (s *Scene) SetSpeedScale(new float64) {
	old := s.SpeedScale

	s.cursor *= new / old

	ds := s.Dynamics.Dynamics
	for i := range ds {
		ds[i].Position *= new / old
	}
	ns := s.notes.notes
	for i := range ns {
		ns[i].Position *= new / old
	}
	bs := s.bar.bars
	for i := range bs {
		bs[i].Position *= new / old
	}

	s.SpeedScale = new
}

func (s *Scene) SetMusicOffset(offset int32) { s.Timer.SetMusicOffset(offset) }

func (s *Scene) Update(kas []input.KeyboardAction) any {
	for _, ka := range kas {
		// Todo: solve this
		// if len(ka.KeyActions) != s.KeyCount {
		// 	fmt.Println("len(ka.KeyActions) != s.KeyCount")
		// 	continue
		// }

		missed := s.Scorer.flushStagedNotes(ka.Time)
		if missed {
			worstJudgment = s.miss()
		}

		s.Dynamic = mode.NextDynamics(s.Dynamic, ka.Time) // for Volume
		s.playSounds(ka)
		js := s.Scorer.tryJudge(ka)

		// draw
		for k, a := range ka.KeyActions {
			switch a {
			case input.Idle, input.Release:
				s.isKeyPresseds[k] = false
				s.isKeyHolds[k] = false
				s.isLongNoteHoldings[k] = false
			case input.Hit:
				s.isKeyPresseds[k] = true
				s.drawKeyTimers[k].Reset()
				s.drawKeyLightingTimers[k].Reset()
				s.drawHitLightingTimers[k].Reset()
				s.drawHoldLightingTimers[k].Reset()
			case input.Hold:
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
		// Todo: Add time error meter mark
		// Todo: Use different color for error meter of Tail
	}

	// update cursor
	duration := float64(s.now - s.Dynamic.Time)
	s.cursor = s.Dynamic.Position + duration*s.Speed()
	return nil
}

// No need to check whether staged note is Tail or not,
// since Tail has no sample in advance.

// Todo: set all sample volumes in advance?
func (s Scene) playSounds(ka input.KeyboardAction) {
	for k, n := range s.stagedNotes {
		if n == nil {
			continue
		}
		a := ka.KeyActions[k]
		if a != input.Hit {
			continue
		}

		sample := mode.DefaultSample
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

func (s Scene) DebugString() string {
	var b strings.Builder
	f := fmt.Fprintf

	f(&b, "Time: %.3fs/%.0fs\n", mode.ToSecond(s.now), mode.ToSecond(s.Duration()))
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

func (asset *Asset) setDefaultHitSound(cfg *Config, fsys fs.FS) {
	streamer, format, _ := audios.DecodeFromFile(fsys, "piano/sound/hit.wav")
	asset.DefaultHitSoundStreamer = streamer
	asset.DefaultHitSoundFormat = format
}

// Alternative names of Mods:
// Modifiers, Parameters
// Occupied: Options, Settings, Configs
// If Mods is gonna be used, it might be good to change "Mode".
