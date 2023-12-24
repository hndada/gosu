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

const (
	maxFlow = 50
	maxAcc  = 20
)

func (asset *Asset) setDefaultHitSound(cfg *Config, fsys fs.FS) {
	streamer, format, _ := audios.DecodeFromFile(fsys, "piano/sound/hit.wav")
	asset.DefaultHitSoundStreamer = streamer
	asset.DefaultHitSoundFormat = format
}

// Alternative names of Mods:
// Modifiers, Parameters
// Occupied: Options, Settings, Configs
// If Mods is gonna be used, it might be good to change "Mode".
type ScenePlay struct {
	// Todo: FlowPoint
	flow       float64
	acc        float64
	unitScores [3]float64

	isKeyPresseds []bool // for keys, key lightings, and hold lightings
	isKeyHolds    []bool // for long note body, hold lightings
	// isJudgeOKs         []bool // for 'hit' lighting
	isLongNoteHoldings []bool // for long note body

	drawKeyTimers          []draws.Timer
	drawKeyLightingTimers  []draws.Timer
	drawHitLightingTimers  []draws.Timer
	drawHoldLightingTimers []draws.Timer
}

// Just assigning slice will shallow copy.
// NewXxx returns struct, while LoadXxx doesn't.
func NewScenePlay(res Resources, opts Options) (s ScenePlay, err error) {
	c.Notes = NewNotes(format, c.KeyCount())
	s.Bars = s.Dynamics.NewBars(c.Duration())

	c.setDynamicPositions()
	c.setNotePositions()
	c.setBarPositions()
	s.Chart.updateTailPosition(cfg.TailExtraDuration)

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
	s.SoundMap.AppendSound("", s.DefaultHitSoundStreamer)

	s.Dynamic = s.Chart.Dynamics[0]

	s.SetSpeedScale()
	s.cursor = float64(s.now) * s.SpeedScale

	s.isKeyPresseds = make([]bool, s.KeyCount)
	s.isKeyHolds = make([]bool, s.KeyCount)
	// s.isJudgeOKs = make([]bool, s.KeyCount)
	s.isLongNoteHoldings = make([]bool, s.KeyCount)
	// s.kool() is just for placeholder.
	s.worstJudgment = s.kool()

	s.drawKeyTimers = s.newDrawTimers(mode.ToTick(30), 0)
	s.drawKeyLightingTimers = s.newDrawTimers(mode.ToTick(30), 0)
	s.drawHitLightingTimers = s.newDrawTimers(mode.ToTick(150), mode.ToTick(150))
	s.drawHoldLightingTimers = s.newDrawTimers(0, mode.ToTick(300))
	return
}

// Need to re-calculate positions when Speed has changed.
func (s *ScenePlay) SetSpeedScale(new, old float64) {
	s.cursor *= new / old
	for _, d := range s.Dynamics {
		d.Position *= new / old
	}
	for _, n := range s.Notes {
		n.Position *= new / old
	}
	for _, b := range s.Bars {
		b.Position *= new / old
	}
}

func (s *ScenePlay) SetMusicOffset(offset int32) { s.Timer.SetMusicOffset(offset) }

func (s *ScenePlay) Update(kas []input.KeyboardAction) any {
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
func (s ScenePlay) playSounds(ka input.KeyboardAction) {
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

func (s ScenePlay) DebugString() string {
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

// func (s *ScenePlay) tickerDrawTimers() {
// 	for k := 0; k < s.KeyCount; k++ {
// 		s.drawKeyTimers[k].Ticker()
// 		s.drawNoteTimers[k].Ticker()
// 		s.drawKeyLightingTimers[k].Ticker()
// 		s.drawHitLightingTimers[k].Ticker()
// 		s.drawHoldLightingTimers[k].Ticker()
// 	}
// 	s.drawJudgmentTimer.Ticker()
// 	s.drawComboTimer.Ticker()
// }
