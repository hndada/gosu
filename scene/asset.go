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
	CursorSprites           [3]draws.Sprite
	DefaultBackgroundSprite draws.Sprite
	BoxMaskSprite           draws.Sprite
	ClearSprite             draws.Sprite
	// Intro   draws.Sprite
	// Loading draws.Sprite

	EnterSound       audios.SoundPlayer
	SwipeSoundPod    audios.SoundPlayer
	TapSoundPod      audios.SoundPlayer
	ToggleSounds     [2]audios.SoundPlayer
	TransitionSounds [2]audios.SoundPlayer

	// Each key count has different asset in piano mode.
	PianoAssets map[int]*piano.Asset
}

func NewAsset(cfg *Config, fsys fs.FS) *Asset {
	asset := &Asset{}
	asset.setCursorSprites(cfg, fsys)
	asset.setDefaultBackgroundSprite(cfg, fsys)
	asset.setBoxMaskSprite(cfg, fsys)
	asset.setClearSprite(cfg, fsys)
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
func (asset *Asset) setCursorSprites(cfg *Config, fsys fs.FS) {
	for i, name := range []string{"base", "additive", "trail"} {
		s := draws.NewSpriteFromFile(fsys, fmt.Sprintf("cursor/%s.png", name))
		s.MultiplyScale(cfg.CursorSpriteScale)
		s.Locate(cfg.ScreenSize.X/2, cfg.ScreenSize.Y/2, draws.LeftTop)
		asset.CursorSprites[i] = s
	}
}

func (asset *Asset) setDefaultBackgroundSprite(cfg *Config, fsys fs.FS) {
	s := NewBackgroundSprite(fsys, "interface/default-bg.jpg", cfg.ScreenSize)
	asset.DefaultBackgroundSprite = s
}

func NewBackgroundSprite(fsys fs.FS, name string,
	screenSize draws.Vector2) draws.Sprite {

	s := draws.NewSpriteFromFile(fsys, name)
	s.MultiplyScale(screenSize.X / s.W())
	s.Locate(screenSize.X/2, screenSize.Y/2, draws.CenterMiddle)
	return s
}

// Todo: MultiplyScale by cfg.ChooseEntryBoxCount
func (asset *Asset) setBoxMaskSprite(cfg *Config, fsys fs.FS) {
	s := draws.NewSpriteFromFile(fsys, "interface/box-mask.png")
	s.Locate(cfg.ScreenSize.X, cfg.ScreenSize.Y/2, draws.RightMiddle)
	// s.MultiplyScale(cfg.ChooseEntryBox)
	asset.BoxMaskSprite = s
}

func (asset *Asset) setClearSprite(cfg *Config, fsys fs.FS) {
	s := draws.NewSpriteFromFile(fsys, "interface/clear.png")
	s.Locate(cfg.ScreenSize.X/2, cfg.ScreenSize.Y/2, draws.CenterMiddle)
	s.MultiplyScale(cfg.ClearSpriteScale)
	asset.ClearSprite = s
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
	format, err := audios.FormatFromFS(subFS)
	if err != nil {
		return
	}
	asset.SwipeSoundPod = audios.NewSoundPod(subFS, format, &cfg.SoundVolume)
}

func (asset *Asset) setTapSoundPod(cfg *Config, fsys fs.FS) {
	subFS, err := fs.Sub(fsys, "sound/tap")
	if err != nil {
		return
	}
	format, err := audios.FormatFromFS(subFS)
	if err != nil {
		return
	}
	asset.TapSoundPod = audios.NewSoundPod(subFS, format, &cfg.SoundVolume)
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
