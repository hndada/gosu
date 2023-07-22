package piano

import (
	"fmt"
	"io/fs"
	"strings"
	"time"

	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/mode"
)

type ScenePlay struct {
	*Config
	*Asset
	Mods
	*Chart

	mode.Timer
	// Store a certain time point to now.
	// Each Now() may yield different time point.
	now int32
	*input.Keyboard
	input.KeyboardReader // for replay
	musicPlayer          audios.MusicPlayer
	audios.SoundMap
	// Scorer also has stagedNotes.
	// ScenePlay.stagedNotes is for playing samples.
	stagedNotes []*Note

	Scorer
	Dynamic *mode.Dynamic // Todo: Dynamic -> dynamic

	// draw
	speedScale   float64
	cursor       float64
	highestBar   *Bar
	highestNotes []*Note

	isKeyPresseds []bool // for keys, key lightings, and hold lightings
	isKeyHolds    []bool // for long note body, hold lightings
	// isJudgeOKs         []bool // for 'hit' lighting
	isLongNoteHoldings []bool // for long note body
	worstJudgment      mode.Judgment

	// draw: animation or transition
	drawKeyTimers          []draws.Timer
	drawNoteTimers         []draws.Timer
	drawKeyLightingTimers  []draws.Timer
	drawHitLightingTimers  []draws.Timer
	drawHoldLightingTimers []draws.Timer
	drawJudgmentTimer      draws.Timer
	drawComboTimer         draws.Timer

	drawScore func(draws.Image)
	drawCombo func(draws.Image)
}

// Todo: initialize s.Asset with s.KeyCount, then set s.Chart.
func NewScenePlay(cfg *Config, assets map[int]*Asset, fsys fs.FS, name string, replay *osr.Format) (s *ScenePlay, err error) {
	s = &ScenePlay{Config: cfg}
	s.Chart, err = NewChart(s.Config, fsys, name)
	if err != nil {
		return
	}
	if s.Chart.KeyCount == 0 {
		return s, fmt.Errorf("key count is zero")
	}

	if _, ok := assets[s.KeyCount]; !ok {
		assets[s.KeyCount] = NewAsset(s.Config, fsys, s.KeyCount, NoScratch)
	}
	s.Asset = assets[s.KeyCount]

	const wait = 1800 * time.Millisecond
	s.Timer = mode.NewTimer(*s.MusicOffset, wait)
	s.now = s.Now()

	if replay != nil {
		s.KeyboardReader = replay.KeyboardReader(s.KeyCount)
	} else {
		keys := input.NamesToKeys(s.KeySettings[s.KeyCount])
		s.Keyboard = input.NewKeyboard(keys, s.StartTime())
		defer s.Keyboard.Listen()
	}

	const ratio = 1
	s.musicPlayer, err = audios.NewMusicPlayerFromFile(fsys, s.MusicFilename, ratio)
	if err != nil {
		return
	}
	s.SoundMap = audios.NewSoundMap(fsys, s.DefaultHitSoundFormat, s.SoundVolume)
	// It is possible for empty string to be a key of a map.
	// https://go.dev/play/p/nn-peGAjawW
	s.SoundMap.AppendSound("", s.DefaultHitSoundStreamer)
	s.stagedNotes = s.newStagedNotes()

	s.Scorer = s.newScorer()
	s.Dynamic = s.Chart.Dynamics[0]

	// draw
	s.speedScale = s.SpeedScale
	s.cursor = float64(s.now) * s.SpeedScale
	s.highestBar = s.Chart.Bars[0]
	// Just assigning slice will shallow copy.
	s.highestNotes = make([]*Note, s.KeyCount)
	copy(s.highestNotes, s.stagedNotes)

	s.isKeyPresseds = make([]bool, s.KeyCount)
	s.isKeyHolds = make([]bool, s.KeyCount)
	// s.isJudgeOKs = make([]bool, s.KeyCount)
	s.isLongNoteHoldings = make([]bool, s.KeyCount)
	// s.kool() is just for placeholder.
	s.worstJudgment = s.kool()

	s.drawKeyTimers = s.newDrawTimers(mode.ToTick(30), 0)
	s.drawNoteTimers = s.newDrawTimers(0, mode.ToTick(400))
	s.drawKeyLightingTimers = s.newDrawTimers(mode.ToTick(30), 0)
	s.drawHitLightingTimers = s.newDrawTimers(mode.ToTick(150), mode.ToTick(150))
	s.drawHoldLightingTimers = s.newDrawTimers(0, mode.ToTick(300))
	s.drawJudgmentTimer = draws.NewTimer(mode.ToTick(250), mode.ToTick(40))
	s.drawComboTimer = draws.NewTimer(mode.ToTick(2000), 0)

	const comboBounce = 0.85
	s.drawScore = mode.NewScoreDrawer(s.ScoreSprites, &s.Score, s.ScoreSpriteScale)
	s.drawCombo = mode.NewComboDrawer(s.ComboSprites, &s.Combo, &s.drawComboTimer, s.ComboDigitGap, comboBounce)
	return
}

func (s ScenePlay) newStagedNotes() []*Note {
	staged := make([]*Note, s.KeyCount)
	for k := range staged {
		for _, n := range s.Chart.Notes {
			if k == n.Key {
				staged[n.Key] = n
				break
			}
		}
	}
	return staged
}

func (s ScenePlay) newDrawTimers(maxTick, period int) []draws.Timer {
	timers := make([]draws.Timer, s.KeyCount)
	for k := range timers {
		timers[k] = draws.NewTimer(maxTick, period)
	}
	return timers
}

func (s ScenePlay) ChartHeader() mode.ChartHeader { return s.Chart.ChartHeader }
func (s ScenePlay) WindowTitle() string           { return s.Chart.WindowTitle() }
func (s ScenePlay) Now() int32                    { return s.Timer.Now() }
func (s ScenePlay) Speed() float64                { return s.Dynamic.Speed * s.SpeedScale }
func (s ScenePlay) IsPaused() bool                { return s.Timer.IsPaused() }
func (s ScenePlay) SetMusicVolume(vol float64)    { s.musicPlayer.SetVolume(vol) }

// Need to re-calculate positions when Speed has changed.
func (s *ScenePlay) SetSpeedScale() {
	c := s.Chart
	old := s.speedScale
	new := s.SpeedScale
	s.cursor *= new / old
	for _, d := range c.Dynamics {
		d.Position *= new / old
	}
	for _, n := range c.Notes {
		n.Position *= new / old
	}
	for _, b := range c.Bars {
		b.Position *= new / old
	}
	s.speedScale = s.SpeedScale
}

func (s *ScenePlay) SetMusicOffset(offset int32) { s.Timer.SetMusicOffset(offset) }

func (s *ScenePlay) Update() any {
	s.now = s.Now()
	s.tryPlayMusic()

	var worstJudgment mode.Judgment
	kas := s.readInput()
	for _, ka := range kas {
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

	// draw
	s.updateCursor()
	s.updateHighestBar()
	s.updateHighestNotes()
	s.tickerDrawTimers()
	return nil
}

// readInput guarantees that length of return value is at least one.
// The receiver should be pointer for updating replay's index.
func (s *ScenePlay) readInput() []input.KeyboardAction {
	if s.Keyboard != nil {
		return s.Keyboard.Read(s.now)
	}
	return s.KeyboardReader.Read(s.now)
}

func (s *ScenePlay) tryPlayMusic() {
	if s.musicPlayer.IsPlayed() {
		return
	}
	if s.now >= *s.MusicOffset && s.now < 300 {
		s.musicPlayer.Play()
		s.Timer.SetMusicPlayed(time.Now())
	}
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

func (s *ScenePlay) updateCursor() {
	duration := float64(s.now - s.Dynamic.Time)
	s.cursor = s.Dynamic.Position + duration*s.Speed()
}

// When speed changes from fast to slow, which means there are more bars
// on the screen, updateHighestBar() will handle it optimally.
// When speed changes from slow to fast, which means there are fewer bars
// on the screen, updateHighestBar() actually does nothing, which is still
// fine because that makes some unnecessary bars are drawn.
// The same concept also applies to notes.
func (s *ScenePlay) updateHighestBar() {
	upperBound := s.cursor + s.ScreenSize.Y + 100
	b := s.highestBar
	for b != nil && b.Position < upperBound {
		if b.Next == nil {
			break
		}
		b = b.Next
		s.highestBar = b
	}
}

func (s *ScenePlay) updateHighestNotes() {
	upperBound := s.cursor + s.ScreenSize.Y + 100
	for k, n := range s.highestNotes {
		// It is possible that n is nil when there is no notes in the whole lane.
		for n != nil && n.Position < upperBound {
			if n.Next == nil {
				break
			}
			n = n.Next
			s.highestNotes[k] = n
		}
		// Update Head to Tail since drawLongNoteBody uses Tail.
		if n != nil && n.Type == Head {
			s.highestNotes[k] = n.Next
		}
	}
}

func (s *ScenePlay) tickerDrawTimers() {
	for k := 0; k < s.KeyCount; k++ {
		s.drawKeyTimers[k].Ticker()
		s.drawNoteTimers[k].Ticker()
		s.drawKeyLightingTimers[k].Ticker()
		s.drawHitLightingTimers[k].Ticker()
		s.drawHoldLightingTimers[k].Ticker()
	}
	s.drawJudgmentTimer.Ticker()
	s.drawComboTimer.Ticker()
}

func (s *ScenePlay) Pause() {
	s.Timer.Pause()
	s.musicPlayer.Pause()
	s.Keyboard.Pause()
}

func (s *ScenePlay) Resume() {
	s.Timer.Resume()
	s.musicPlayer.Resume()
	s.Keyboard.Resume()
}

func (s ScenePlay) Finish() any {
	s.musicPlayer.Close()
	s.Keyboard.Close()
	return s.Scorer
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
	f(&b, "\n")
	return b.String()
}
