package piano

import (
	"io/fs"
	"time"

	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/draws"
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
	audios.MusicPlayer
	audios.SoundMap
	stagedNotes []*Note // Scorer has same slice. This is for playing samples.

	Scorer
	Dynamic *mode.Dynamic

	// draw
	speedScale   float64
	cursor       float64
	highestBar   *Bar
	highestNotes []*Note

	isKeyHolds    []bool // long note body, hold lightings
	isKeyPresseds []bool // keys, key lightings, and hold lightings
	isNoteHits    []bool // 'hit' lighting
	worstJudgment mode.Judgment

	// draw: animation or transition
	keyTimers          []draws.Timer
	noteTimers         []draws.Timer
	keyLightingTimers  []draws.Timer
	hitLightingTimers  []draws.Timer
	holdLightingTimers []draws.Timer
	judgmentTimer      draws.Timer
	comboTimer         draws.Timer

	drawScore func(draws.Image)
	drawCombo func(draws.Image)
}

// Todo: pass key count beforehand so that s.Asset can be initialized before s.Chart.
func NewScenePlay(cfg *Config, assets map[int]*Asset, fsys fs.FS, name string, mods Mods, replay input.KeyboardReader) (s *ScenePlay, err error) {
	s = &ScenePlay{Config: cfg}
	s.Mods = mods
	s.Chart, err = NewChart(s.Config, fsys, name, mods)
	if err != nil {
		return
	}
	s.Asset = assets[s.KeyCount]

	const wait = 1800 * time.Millisecond
	s.Timer = mode.NewTimer(*s.MusicOffset, wait)
	s.now = s.Now()

	if replay.IsEmpty() {
		keys := input.NamesToKeys(s.KeySettings[s.KeyCount])
		s.Keyboard = input.NewKeyboard(keys, s.StartTime())
		defer s.Keyboard.Listen()
	} else {
		s.KeyboardReader = replay
	}

	const ratio = 1
	s.MusicPlayer, err = audios.NewMusicPlayerFromFile(fsys, s.MusicFilename, ratio)
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
	s.highestNotes = make([]*Note, len(s.Chart.Notes))
	copy(s.highestNotes, s.stagedNotes)

	// Since timers are now updated in Draw(), their ticks would be dependent on FPS.
	// However, so far TPS and FPS goes synced by SyncWithFPS().
	s.keyTimers = s.newTimers(mode.ToTick(30), 0)
	s.noteTimers = s.newTimers(0, mode.ToTick(400))
	s.keyLightingTimers = s.newTimers(mode.ToTick(30), 0)
	s.hitLightingTimers = s.newTimers(mode.ToTick(150), mode.ToTick(150))
	s.holdLightingTimers = s.newTimers(0, mode.ToTick(300))
	s.judgmentTimer = draws.NewTimer(mode.ToTick(250), mode.ToTick(40))
	s.comboTimer = draws.NewTimer(mode.ToTick(2000), 0)

	const comboBounce = 0.85
	s.drawScore = mode.NewScoreDrawer(s.ScoreSprites, &s.Score, s.ScoreSpriteScale)
	s.drawCombo = mode.NewComboDrawer(s.ComboSprites, &s.Combo, &s.comboTimer, s.ComboDigitGap, comboBounce)
	return
}

func (s ScenePlay) newStagedNotes() []*Note {
	staged := make([]*Note, s.KeyCount)
	for k := range staged {
		for _, n := range s.Notes {
			if k == n.Key {
				staged[n.Key] = n
				break
			}
		}
	}
	return staged
}

// Todo: Timer -> DrawTimer
func (s ScenePlay) newTimers(maxTick, period int) []draws.Timer {
	timers := make([]draws.Timer, s.Chart.KeyCount)
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
func (s ScenePlay) SetMusicVolume(vol float64)    { s.MusicPlayer.SetVolume(vol) }

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

	// draw
	s.isKeyHolds = make([]bool, s.KeyCount)
	s.isKeyPresseds = make([]bool, s.KeyCount)
	s.isNoteHits = make([]bool, s.KeyCount)
	s.worstJudgment = s.kool()
	for _, ka := range s.readInput() {
		s.Scorer.flushStagedNotes(ka.Time)
		s.Dynamic = mode.NextDynamics(s.Dynamic, ka.Time) // for Volume
		s.playSounds(ka)
		js := s.Scorer.tryJudge(ka)

		// draw
		for k, a := range ka.KeyActions {
			switch a {
			case input.Hit:
				s.isKeyPresseds[k] = true
				s.keyTimers[k].Reset()
				s.keyLightingTimers[k].Reset()
				s.hitLightingTimers[k].Reset()
				s.holdLightingTimers[k].Reset()
			case input.Hold:
				s.isKeyPresseds[k] = true
				s.isKeyHolds[k] = true
			}
		}

		for k, j := range js {
			// Tail also makes hit lighting on.
			if !j.Is(s.miss()) {
				s.isNoteHits[k] = true
			}
			if s.worstJudgment.Window < j.Window {
				s.worstJudgment = j
			}
		}

		// Todo: Add time error meter mark
		// Todo: Use different color for error meter of Tail
	}

	// draw
	s.updateCursor()
	s.updateHighestBar()
	s.updateHighestNotes()
	s.ticker()
	return nil
}

// readInput guarantees that it length is at least one.
func (s ScenePlay) readInput() []input.KeyboardAction {
	if s.Keyboard != nil {
		return s.Keyboard.Read(s.now)
	}
	return s.KeyboardReader.Read(s.now)
}

func (s *ScenePlay) tryPlayMusic() {
	if s.MusicPlayer.IsPlayed() {
		return
	}
	if s.now >= *s.MusicOffset && s.now < 300 {
		s.MusicPlayer.Play()
		s.Timer.SetMusicPlayed(time.Now())
	}
}

// No need to check whether staged note is Tail or not,
// since Tail has no sample in advance.

// Todo: set all sample volumes in advance?
func (s ScenePlay) playSounds(ka input.KeyboardAction) {
	for k, n := range s.stagedNotes {
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

func (s *ScenePlay) ticker() {
	for k := 0; k < s.KeyCount; k++ {
		s.keyTimers[k].Ticker()
		s.noteTimers[k].Ticker()
		s.keyLightingTimers[k].Ticker()
		s.hitLightingTimers[k].Ticker()
		s.holdLightingTimers[k].Ticker()
	}
	s.judgmentTimer.Ticker()
	s.comboTimer.Ticker()
}

func (s *ScenePlay) Pause() {
	s.Timer.Pause()
	s.MusicPlayer.Pause()
	s.Keyboard.Pause()
}

func (s *ScenePlay) Resume() {
	s.Timer.Resume()
	s.MusicPlayer.Resume()
	s.Keyboard.Resume()
}

func (s ScenePlay) Finish() any {
	s.MusicPlayer.Close()
	s.Keyboard.Close()
	return s.Scorer
}
