package play

import (
	"fmt"
	"io/fs"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/times"
)

type play interface {
	Update(now int32, kas []game.KeyboardAction) any
	SampleBuffer() []game.Sample
	Draw(dst draws.Image)
}

// Todo: draw 4:3 screen on 16:9 screen
type Scene struct {
	game.ChartHeader

	// Timer
	startTime   time.Time
	pauseTime   time.Time
	paused      bool
	musicOffset int32
	musicPlayed bool // This really matters.

	keyboard    input.Keyboard
	musicPlayer audios.MusicPlayer
	soundPlayer audios.SoundPlayer

	play play
}

func NewScene(fsys fs.FS, name string) (s Scene, err error) {
	format, hash, err := game.LoadChartFile(fsys, name)
	if err != nil {
		err = fmt.Errorf("failed to load chart file: %w", err)
		return
	}
	s.ChartHeader = game.NewChartHeader(format, hash)

	const wait = 1100 * time.Millisecond
	s.startTime = times.Now().Add(wait)
	s.musicOffset = musicOffset

	if replay != nil {
		s.keyboard = game.NewReplay(replayFile)
	} else {
		keys := input.NamesToKeys(s.KeySettings[s.KeyCount])
		s.keyboard = input.NewKeyboard(keys)
		defer s.keyboard.Listen(s.startTime)
	}

	mp, err := audios.NewMusicPlayerFromFile(fsys, s.MusicFilename)
	if err != nil {
		err = fmt.Errorf("failed to load music file: %w", err)
		return
	}
	s.musicPlayer = mp
	s.SetMusicVolume(*s.MusicVolume)

	s.SoundMap = audios.NewSoundMap(fsys, s.DefaultHitSoundFormat, s.SoundVolume)
	// It is possible for empty string to be a key of a map.
	// https://go.dev/play/p/nn-peGAjawW
	s.SoundMap.AddSound("", s.DefaultHitSoundStreamer)

	ebiten.SetWindowTitle(s.WindowTitle())
	return
}

func (s Scene) WindowTitle() string {
	c := s.ChartHeader
	return fmt.Sprintf("gosu | %s - %s [%s] (%s) ", c.Artist, c.MusicName, c.ChartName, c.Charter)
}

func (s Scene) Now() int32 {
	var d time.Duration
	if s.paused {
		d = s.pauseTime.Sub(s.startTime)
	} else {
		d = times.Since(s.startTime)
	}
	return int32(d.Milliseconds()) // + s.musicOffset
}

// TL;DR: If you tend to hit early, set positive offset.
// It leads to delayed music / early start time.
func (s *Scene) SetMusicOffset(new int32) {
	// Once the music starts, there isn't much we can do,
	// since music is hard to seek precisely.
	// Instead, we adjust the start time.
	if s.musicPlayed {
		// Positive offset in music infers music is delayed.
		// Delayed music is same as early start time.
		// Hence, as offset increases, start time decreases.
		// This leads to a instant, slight movement of notes.

		// Changing offset might affect to indexing keyboard input,
		// but it would be negligible because a player tend to
		// hands off the keys for a while when setting offset.
		// Maybe the fact that 'osu!' doesn't allow to change offset
		// during pausing is related to this.

		// No need to consider playback rate, since it is supported naturally.
		// Times themselves are not affected, only now flows faster or slower.
		old := s.musicOffset
		diff := time.Duration(new-old) * time.Millisecond
		s.startTime = s.startTime.Add(-diff)
		s.musicOffset = new
	} else {
		// If the music has not played yet, we can adjust the offset
		// and let the music played at given delayed time.
		s.musicOffset = new
	}
}

func (s *Scene) Update() any {
	now := s.Now() // Use unified time.

	// No update t.startTime when playing music, unless
	// notes would look like they suddenly teleport at the beginning.
	if !s.musicPlayed && now >= s.musicOffset {
		s.musicPlayer.Play()
		s.musicPlayed = true
	}
	nowTime := time.Duration(now) * time.Millisecond
	kss := s.keyboard.Read(nowTime)
	kas := game.KeyboardActions(kss)
	r := s.play.Update(now, kas)
	s.PlaySounds()
	return r
}

func (s Scene) PlaySounds() {
	for _, samp := range s.play.SampleBuffer() {
		vol := samp.Volume * s.SoundVolumeScale
		s.soundPlayer.Play(samp.Filename, vol)
	}
}

func (s *Scene) Pause() {
	s.pauseTime = times.Now()
	s.musicPlayer.Pause()
	s.keyboard.Stop()
	s.paused = true

}

func (s *Scene) Resume() {
	elapsedTime := times.Since(s.pauseTime)
	s.startTime = s.startTime.Add(elapsedTime)
	s.musicPlayer.Resume()
	s.keyboard.Listen(s.startTime)
	s.paused = false
}

// Music keeps playing at result scene.
func (s *Scene) Close() {
	// s.MusicPlayer.Close()
	s.keyboard.Stop()
}

// func (t *Timer) sync() {
// 	const threshold = 30 * 1000
// 	since := int32(times.Since(t.startTime).Milliseconds())
// 	if e := since - t.Now(); e >= threshold {
// 		fmt.Printf("%dms: adjusting time error (%dms)\n", since, e)
// 		t.Tick += ToTick(e)
// 	}
// }
