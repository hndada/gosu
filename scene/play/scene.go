package play

import (
	"fmt"
	"io/fs"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/scene"
	"github.com/hndada/gosu/times"
)

// Todo: draw 4:3 screen on 16:9 screen
type Scene struct {
	chartHeader game.ChartHeader

	// Timer
	startTime   time.Time
	pauseTime   time.Time
	paused      bool
	musicOffset int32
	musicPlayed bool // This really matters.

	keyboard    input.Keyboard
	musicPlayer audios.MusicPlayer
	soundPlayer audios.SoundPlayer

	play scene.Scene
}

func NewScene(fsys fs.FS, name string, replayFile fs.File) (s Scene, err error) {
	format, hash, err := game.LoadChartFile(fsys, name)
	if err != nil {
		err = fmt.Errorf("failed to load chart file: %w", err)
		return
	}
	s.chartHeader = game.NewChartHeader(format, hash)

	const wait = 1100 * time.Millisecond
	s.startTime = times.Now().Add(wait)
	s.musicOffset = musicOffset

	if replayFile == nil {
		keys := input.NamesToKeys(s.KeySettings[s.KeyCount])
		s.keyboard = input.NewKeyboard(keys)
		defer s.keyboard.Listen(s.startTime)
	} else {
		s.KeyboardReader = replay.KeyboardReader(s.KeyCount)
	}

	const ratio = 1
	s.musicPlayer, _ = audios.NewMusicPlayerFromFile(fsys, s.MusicFilename, ratio)
	s.SetMusicVolume(*s.MusicVolume)

	s.SoundMap = audios.NewSoundMap(fsys, s.DefaultHitSoundFormat, s.SoundVolume)
	// It is possible for empty string to be a key of a map.
	// https://go.dev/play/p/nn-peGAjawW
	s.SoundMap.AddSound("", s.DefaultHitSoundStreamer)

	ebiten.SetWindowTitle(s.WindowTitle())
	return
}

func (s Scene) WindowTitle() string {
	c := s.chartHeader
	return fmt.Sprintf("gosu | %s - %s [%s] (%s) ", c.Artist, c.MusicName, c.ChartName, c.Charter)
}

func (s Scene) Now() time.Duration {
	if s.paused {
		return s.pauseTime.Sub(s.startTime)
	} else {
		return times.Since(s.startTime)
	}
}

// TL;DR: If you tend to hit early, set positive offset.
// It leads to delayed music / early start time.
func (s *Scene) SetMusicOffset(new int32) {
	// Once the music starts, there isn't much we can do,
	// since music is hard to seek precisely.
	// Instead, we adjust the start time.

	// Positive offset in music infers music is delayed.
	// Delayed music is same as early start time.
	// Hence, as offset increases, start time decreases.
	if s.musicPlayed {
		old := s.musicOffset
		diff := time.Duration(new-old) * time.Millisecond
		s.startTime = s.startTime.Add(-diff)
		s.musicOffset = new
	}
	// If the music has not played yet, we can adjust the offset
	// and let the music played at given delayed time.
	s.musicOffset = new

	// Changing offset might affect to KeyboardState indexing,
	// but it would be negligible because a player tend to hands off the keys
	// when setting offset. Maybe the fact that osu! doesn't allow to change offset
	// during pausing is related to this.
}

// No update t.startTime here, unless notes would look
// like they suddenly teleport at the beginning.
func (s *Scene) tryPlayMusic() {
	if s.MusicPlayer.IsPlayed() {
		return
	}
	now := s.Timer.Now()
	if now >= *s.MusicOffset && now < 300 {
		s.MusicPlayer.Play()
		s.musicPlayed = true
	}
}

// readInput guarantees that length of return value is at least one.
// The receiver should be pointer for updating replay's index.
func (s *Scene) readInput() []input.KeyboardAction {
	if s.Keyboard != nil {
		return s.Keyboard.Read(s.Timer.Now())
	}
	return s.Keyboard.Reader.Read(s.Timer.Now()) // for replay
}

func (s *Scene) Update() any {
	s.tryPlayMusic()
	s.readInput()
	return s.Play.Update()
}

func (s *Scene) Pause() {
	s.pauseTime = times.Now()
	s.paused = true
	s.MusicPlayer.Pause()
	s.Keyboard.Stop()
}

func (s *Scene) Resume() {
	elapsedTime := times.Since(t.pauseTime)
	t.startTime = t.startTime.Add(elapsedTime)
	t.paused = false
	s.MusicPlayer.Resume()
	s.Keyboard.Listen()
}

// func (t *Timer) sync() {
// 	const threshold = 30 * 1000
// 	since := int32(times.Since(t.startTime).Milliseconds())
// 	if e := since - t.Now(); e >= threshold {
// 		fmt.Printf("%dms: adjusting time error (%dms)\n", since, e)
// 		t.Tick += ToTick(e)
// 	}
// }

func (s *Scene) Close() {
	// Music keeps playing at result scene.
	// s.MusicPlayer.Close()
	s.Keyboard.Stop()
}
