package play

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/game/piano"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/scene"
	"github.com/hndada/gosu/times"
)

type play interface {
	Update(now int32, kas []game.KeyboardAction) any
	// PopSamples() []game.Sample
	Draw(dst draws.Image)
	DebugString() string
}

// Todo: draw 4:3 screen on 16:9 screen
type Scene struct {
	*scene.Resources
	*scene.Options

	*game.ChartHeader
	play              play
	musicPlayer       audios.MusicPlayer
	keyboard          input.KeyboardReader
	lastKeyboardState input.KeyboardState

	// Timer
	firstUpdated bool
	startTime    time.Time
	pauseTime    time.Time
	paused       bool
	musicOffset  int32
	musicPlayed  bool // This really matters.
}

// (*Scene, error) is typically used for regular functions that operate on struct pointers.
// (s *Scene, err error) is typically used for methods attached to structs.
// chartFS fs.FS, cname string, replayFS fs.FS, rname string, mods game.Mods) (*Scene, error) {
func NewScene(res *scene.Resources, opts *scene.Options, args scene.PlayArgs) (*Scene, error) {
	s := &Scene{
		Resources: res,
		Options:   opts,
	}

	switch opts.Mode {
	case game.ModePiano:
		mods := args.Mods.(piano.Mods)
		c, err := piano.NewChart(args.ChartFS, args.ChartFilename, mods)
		if err != nil {
			err = fmt.Errorf("failed to create chart: %w", err)
			return nil, err
		}

		s.ChartHeader = c.ChartHeader
		// Todo: add default sound
		// soft-hitnormal.wav
		sp := s.newSamplePlayer(args.ChartFS, s.MusicFilename)

		play, err := piano.NewPlay(res.Piano, opts.Piano, c, mods, &sp)
		if err != nil {
			err = fmt.Errorf("failed to create play scene: %w", err)
			return nil, err
		}
		s.play = play
	}

	mp, err := audios.NewMusicPlayerFromFile(args.ChartFS, s.MusicFilename)
	if err != nil {
		err = fmt.Errorf("failed to load music file: %w", err)
		return nil, err
	}
	s.musicPlayer = mp
	mp.SetVolume(s.MusicVolume)
	s.musicOffset = s.MusicOffset

	var keyCount int
	var keyNames []string
	switch opts.Mode {
	case game.ModePiano:
		keyCount = opts.SubMode
		keyNames = opts.Piano.KeyMappings[keyCount]
	case game.ModeDrum:
		keyCount = 4
		// keyNames = opts.Drum.Key.Mappings
	}

	if args.ReplayFS != nil {
		kb, _, err := game.NewReplay(args.ReplayFS, args.ReplayFilename, keyCount)
		if err != nil {
			err = fmt.Errorf("failed to load replay file: %w", err)
			return nil, err
		}
		s.keyboard = kb
	} else {
		keys := input.NamesToKeys(keyNames)
		s.keyboard = input.NewKeyboard(keys)
	}

	return s, nil
}

func (s Scene) newSamplePlayer(fsys fs.FS, musicFilename string) audios.SoundPlayer {
	sp := audios.NewSoundPlayer(&s.SoundVolumeScale)
	fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		f, err := fsys.Open(path)
		if err != nil {
			return err
		}
		info, err := f.Stat()
		if err != nil {
			return err
		}

		// Skip music file.
		if info.Name() == musicFilename {
			return nil
		}

		switch ext := filepath.Ext(path); ext {
		case ".wav", ".ogg", ".mp3", ".OGG", ".WAV", ".MP3":
			sp.AddFile(fsys, path)
		}
		return nil
	})
	times.Now()
	return sp
}

// Now returns current time in millisecond.
func (s Scene) Now() time.Duration {
	var d time.Duration
	if s.paused {
		d = s.pauseTime.Sub(s.startTime)
	} else {
		d = times.Since(s.startTime)
	}
	return d
	// return int32(d.Milliseconds()) // + s.musicOffset
}

// TL;DR: If you tend to hit early, set positive offset.
// It leads to delayed music / early start time.
func (s *Scene) SetMusicOffset(newOffset int32) {
	// Once the music starts, there isn't much we can do,
	// since music is hard to seek precisely.
	// Instead, we adjust the start time.
	if s.musicPlayed {
		// Positive offset in music infers music is delayed.
		// Delayed music is same as early start time.
		// Hence, as offset increases, start time decreases.
		// This leads to a instant, slight movement of notes.

		// Changing offset is fine even in pausing. KeyboardStateBuffer
		// won't return any redundant states except the current index,
		// which is not contained in KeyboardAction.
		// c.f. 'osu!' doesn't allow players to change offset during pausing.

		// No need to consider playback rate, since it is supported naturally.
		// Times themselves are not affected, only now flows faster or slower.
		oldOffset := s.musicOffset
		diff := time.Duration(newOffset-oldOffset) * time.Millisecond
		s.startTime = s.startTime.Add(-diff)
		s.musicOffset = newOffset
	} else {
		// If the music has not played yet, we can adjust the offset
		// and let the music played at given delayed time.
		s.musicOffset = newOffset
	}
}

func (s *Scene) firstUpdate() {
	const wait = 1800 * time.Millisecond
	s.startTime = times.Now().Add(wait)
	if kb, ok := s.keyboard.(*input.Keyboard); ok {
		kb.Listen(s.startTime)
	}
	// s.startTime = times.Now() // TODO: OK to comment out?
}

func (s *Scene) Update() any {
	if !s.firstUpdated {
		s.firstUpdate()
		s.firstUpdated = true
	}

	// Use unified time.
	now := s.Now()
	nowMS := int32(now.Milliseconds())
	// nowDuration := time.Duration(now) * time.Millisecond

	// No update t.startTime when playing music, unless
	// notes would look like they suddenly teleport at the beginning.
	if !s.musicPlayed && nowMS >= s.musicOffset {
		s.musicPlayer.Play()
		s.musicPlayed = true
	}

	// kss's length is mostly 1.
	kss := s.keyboard.Read(now)
	kss = append([]input.KeyboardState{s.lastKeyboardState}, kss...)
	kas := game.KeyboardActions(kss)
	r := s.play.Update(nowMS, kas)
	s.lastKeyboardState = kss[len(kss)-1]

	// s.PlaySounds()
	return r
}

// func (s Scene) PlaySounds() {
// 	for _, samp := range s.play.SampleBuffer() {
// 		vol := samp.Volume * s.Audio.SoundVolumeScale
// 		s.soundPlayer.PlayWithVolume(samp.Filename, vol)
// 	}
// }

func (s *Scene) Pause() {
	s.pauseTime = times.Now()
	s.musicPlayer.Pause()
	if kb, ok := s.keyboard.(*input.Keyboard); ok {
		kb.Stop()
	}
	s.paused = true
}

func (s *Scene) Resume() {
	elapsedTime := times.Since(s.pauseTime)
	s.startTime = s.startTime.Add(elapsedTime)
	s.musicPlayer.Resume()
	if kb, ok := s.keyboard.(*input.Keyboard); ok {
		kb.Listen(s.startTime)
	}
	s.paused = false
}

// Music keeps playing at result scene.
func (s *Scene) Close() {
	// s.MusicPlayer.Close()
	if kb, ok := s.keyboard.(*input.Keyboard); ok {
		kb.Stop()
	}
}

func (s Scene) Draw(dst draws.Image) {
	s.play.Draw(dst)
}

func (s Scene) DebugString() string {
	return s.play.DebugString()
}

// func (t *Timer) sync() {
// 	const threshold = 30 * 1000
// 	since := int32(times.Since(t.startTime).Milliseconds())
// 	if e := since - t.Now(); e >= threshold {
// 		fmt.Printf("%dms: adjusting time error (%dms)\n", since, e)
// 		t.Tick += ToTick(e)
// 	}
// }
