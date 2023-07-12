package piano

import (
	"io/fs"
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
	*Chart
	Scorer
	mode.BaseScenePlay

	lastSpeedScale float64
	cursor         float64
	highestBar     *Bar
	highestNotes   []*Note
	lastKeyActions []input.KeyActionType

	// For animation or transition.
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

// Todo: replay listener
func NewScenePlay(cfg *Config, assets map[int]*Asset, fsys fs.FS, name string,
	mods Mods, rf *osr.Format) (s *ScenePlay, err error) {
	s = new(ScenePlay)
	s.Config = cfg
	s.Chart, err = NewChart(cfg, fsys, name, mods)
	if err != nil {
		return
	}
	s.Asset = assets[s.KeyCount]
	s.Scorer = NewScorer(s.Chart)

	const ratio = 1
	s.MusicPlayer, err = audios.NewMusicPlayerFromFile(fsys, s.MusicFilename, ratio)
	if err != nil {
		return
	}
	s.SoundMap = audios.NewSoundMap(fsys, s.SoundVolume)

	const wait = 1800 * time.Millisecond
	if rf != nil {
		s.Keyboard = NewReplayListener(rf, s.KeyCount, wait)
	} else {
		keys := input.NamesToKeys(s.KeySettings[s.KeyCount])
		s.Keyboard = input.NewKeyboardListener(keys, wait)
	}

	s.Dynamic = s.Chart.Dynamics[0]
	s.lastSpeedScale = s.SpeedScale
	s.cursor = float64(s.Now()) * s.SpeedScale
	s.highestBar = s.Chart.Bars[0]
	s.highestNotes = make([]*Note, s.Chart.KeyCount)

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
	s.drawScore = mode.NewDrawScoreFunc(s.ScoreSprites, &s.Scorer.Score,
		s.ScoreSpriteScale)
	s.drawCombo = mode.NewDrawComboFunc(s.ComboSprites, &s.Scorer.Combo, &s.comboTimer,
		s.ComboDigitGap, comboBounce)
	return
}

func (s ScenePlay) newTimers(maxTick, period int) []draws.Timer {
	timers := make([]draws.Timer, s.Chart.KeyCount)
	for k := range timers {
		timers[k] = draws.NewTimer(maxTick, period)
	}
	return timers
}

func (s ScenePlay) Speed() float64 { return s.Dynamic.Speed * s.SpeedScale }

// Need to re-calculate positions when Speed has changed.
func (s *ScenePlay) SetSpeedScale() {
	c := s.Chart
	old := s.lastSpeedScale
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
	s.lastSpeedScale = s.SpeedScale
}

func (s *ScenePlay) Update() any {
	s.Ticker()

	s.Scorer.worstJudgment = blank
	var kas []input.KeyboardAction
	for _, ka := range kas {
		s.Scorer.Check(ka)
		for k, n := range s.Scorer.Staged {
			a := ka.Action[k]
			if n.Type != Tail && a == input.Hit {
				s.PlaySound(n.Sample, *s.SoundVolume)
			}
		}
	}
	if len(kas) > 0 {
		s.lastKeyActions = kas[len(kas)-1].Action
	}

	// Changed speed might not be applied after positions are calculated.
	// But this is not tested.
	s.updateHighestBar()
	s.updateHighestNotes()
	s.UpdateDynamic()
	s.updateCursor()
	return nil
}

func (s *ScenePlay) updateCursor() {
	duration := float64(s.Now() - s.Dynamic.Time)
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
	for b := s.highestBar; b.Position < upperBound; b = b.Next {
		s.highestBar = b
		if b.Next == nil {
			break
		}
	}
}
func (s *ScenePlay) updateHighestNotes() {
	upperBound := s.cursor + s.ScreenSize.Y + 100
	for k, n := range s.highestNotes {
		for ; n.Position < upperBound; n = n.Next {
			s.highestNotes[k] = n
			if n.Next == nil {
				break
			}
		}
	}
}

func (s ScenePlay) WindowTitle() string {
	return s.Chart.WindowTitle()
}

func (s ScenePlay) BackgroundFilename() string {
	return s.Chart.ImageFilename
}

// ExposureTime is the time that cursor takes to move 1 logical pixel.
// The unit is millisecond.
func (s ScenePlay) ExposureTime() int32 {
	return int32(s.HitPosition / s.Speed())
}

func (s ScenePlay) Finish() any {
	s.BaseScenePlay.Finish()
	return s.Scorer
}

func (s ScenePlay) isKeyHit(k int) bool {
	return s.lastKeyActions[k] == input.Hit
}

func (s ScenePlay) isKeyPressed(k int) bool {
	return s.lastKeyActions[k] == input.Hit ||
		s.lastKeyActions[k] == input.Hold
}
