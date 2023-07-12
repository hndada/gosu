package scene

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/mode/piano"
)

const (
	CursorBase = iota
	CursorAdditive
	CursorTrail
)

// Asset is previously known as Skin.
type Asset struct {
	Cursor            [3]draws.Sprite
	DefaultBackground draws.Sprite
	BoxMask           draws.Sprite
	Clear             draws.Sprite
	// Intro   draws.Sprite
	// Loading draws.Sprite

	EnterSound       audios.Sound
	SwipeSoundPod    audios.SoundPod
	TapSoundPod      audios.SoundPod
	ToggleSounds     [2]audios.Sound
	TransitionSounds [2]audios.Sound

	// Each key count has different asset in piano mode.
	PianoAssets map[int]*piano.Asset
}

func NewAsset(cfg *Config, fsys fs.FS) *Asset {
	asset := &Asset{}
	asset.setCursor(cfg, fsys)
	asset.setDefaultBackground(cfg, fsys)
	asset.setBoxMask(cfg, fsys)
	asset.setClear(cfg, fsys)
	// Intro   draws.Sprite
	// Loading draws.Sprite

	asset.setEnterSound(cfg, fsys)
	asset.setSwipeSoundPod(cfg, fsys)
	asset.setTapSoundPod(cfg, fsys)
	asset.setToggleSounds(cfg, fsys)
	asset.setTransitionSounds(cfg, fsys)

	asset.PianoAssets = make(map[int]*piano.Asset)
	for _, keyCount := range []int{4, 7} {
		pianoAsset := piano.NewAsset(cfg.PianoConfig, fsys, keyCount, piano.NoScratch)
		asset.PianoAssets[keyCount] = pianoAsset
	}
	return asset
}

// Cursor should be at CenterMiddle in circle mode (in far future)
func (asset *Asset) setCursor(cfg *Config, fsys fs.FS) {
	for i, name := range []string{"base", "additive", "trail"} {
		s := draws.NewSpriteFromFile(fsys, fmt.Sprintf("cursor/%s.png", name))
		s.MultiplyScale(cfg.CursorScale)
		s.Locate(cfg.ScreenSize.X/2, cfg.ScreenSize.Y/2, draws.LeftTop)
		asset.Cursor[i] = s
	}
}

func (asset *Asset) setDefaultBackground(cfg *Config, fsys fs.FS) {
	s := NewBackgroundSprite(fsys, "interface/default-bg.jpg", cfg.ScreenSize)
	asset.DefaultBackground = s
}

// Todo: MultiplyScale by cfg.ChooseEntryBoxCount
func (asset *Asset) setBoxMask(cfg *Config, fsys fs.FS) {
	s := draws.NewSpriteFromFile(fsys, "interface/box-mask.png")
	s.Locate(cfg.ScreenSize.X, cfg.ScreenSize.Y/2, draws.RightMiddle)
	// s.MultiplyScale(cfg.CursorScale)
	asset.BoxMask = s
}

func (asset *Asset) setClear(cfg *Config, fsys fs.FS) {
	s := draws.NewSpriteFromFile(fsys, "interface/clear.png")
	s.Locate(cfg.ScreenSize.X/2, cfg.ScreenSize.Y/2, draws.CenterMiddle)
	s.MultiplyScale(cfg.ClearScale)
	asset.Clear = s
}

// Intro   draws.Sprite
// Loading draws.Sprite

func (asset *Asset) setEnterSound(cfg *Config, fsys fs.FS) {
	name := "sound/ringtone2_loop.wav"
	asset.EnterSound = audios.NewSound(fsys, name, &cfg.SoundVolume)
}

func (asset *Asset) setSwipeSoundPod(cfg *Config, fsys fs.FS) {
	subFS, err := fs.Sub(fsys, "sound/swipe")
	if err != nil {
		return
	}
	asset.SwipeSoundPod = audios.NewSoundPod(subFS, &cfg.SoundVolume)
}

func (asset *Asset) setTapSoundPod(cfg *Config, fsys fs.FS) {
	subFS, err := fs.Sub(fsys, "sound/tap")
	if err != nil {
		return
	}
	asset.TapSoundPod = audios.NewSoundPod(subFS, &cfg.SoundVolume)
}

func (asset *Asset) setToggleSounds(cfg *Config, fsys fs.FS) {
	for i, name := range []string{"off", "on"} {
		name := fmt.Sprintf("sound/toggle/%s.wav", name)
		asset.ToggleSounds[i] = audios.NewSound(fsys, name, &cfg.SoundVolume)
	}
}

func (asset *Asset) setTransitionSounds(cfg *Config, fsys fs.FS) {
	for i, name := range []string{"down", "up"} {
		name := fmt.Sprintf("sound/transition/%s.wav", name)
		asset.TransitionSounds[i] = audios.NewSound(fsys, name, &cfg.SoundVolume)
	}
}
