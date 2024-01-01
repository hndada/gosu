package scene

import (
	"fmt"
	"image/color"
	"io/fs"

	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/draws"
)

type Asset struct {
	DefaultBackgroundSprite draws.Sprite
	BoxMaskSprite           draws.Sprite
	LoadingSprite           draws.Sprite
	SearchBoxSprite         draws.Sprite

	EnterSound       audios.SoundPlayer
	SwipeSoundPod    audios.SoundPlayer
	TapSoundPod      audios.SoundPlayer
	ToggleSounds     [2]audios.SoundPlayer
	TransitionSounds [2]audios.SoundPlayer
}

func NewAsset(cfg *Config, fsys fs.FS) *Asset {
	asset := &Asset{}
	asset.setDefaultBackgroundSprite(cfg, fsys)
	asset.setBoxMaskSprite(cfg, fsys)
	asset.setLoadingSprite(cfg, fsys)
	asset.setSearchBoxSprite(cfg, fsys)

	asset.setEnterSound(cfg, fsys)
	asset.setSwipeSoundPod(cfg, fsys)
	asset.setTapSoundPod(cfg, fsys)
	asset.setToggleSounds(cfg, fsys)
	asset.setTransitionSounds(cfg, fsys)
	return asset
}

func (asset *Asset) setDefaultBackgroundSprite(cfg *Config, fsys fs.FS) {
	s := draws.NewSpriteFromFile(fsys, "interface/default-bg.jpg")
	s.MultiplyScale(cfg.ScreenSize.X / s.Width())
	s.Locate(cfg.ScreenSize.X/2, cfg.ScreenSize.Y/2, draws.CenterMiddle)
	asset.DefaultBackgroundSprite = s
}

func (asset *Asset) setBoxMaskSprite(cfg *Config, fsys fs.FS) {
	s := draws.NewSpriteFromFile(fsys, "interface/box-mask.png")
	s.Locate(cfg.ScreenSize.X, cfg.ScreenSize.Y/2, draws.RightMiddle)
	s.SetSize(cfg.ChartTreeNodeWidth, cfg.ChartTreeNodeHeight)
	asset.BoxMaskSprite = s
}

func (asset *Asset) setLoadingSprite(cfg *Config, fsys fs.FS) {
	s := draws.NewSpriteFromFile(fsys, "interface/loading.png")
	s.Locate(cfg.ScreenSize.X/2, cfg.ScreenSize.Y/2, draws.CenterMiddle)
	asset.LoadingSprite = s
}

func (asset *Asset) setSearchBoxSprite(cfg *Config, fsys fs.FS) {
	img := draws.NewImage(cfg.SearchBoxWidth, cfg.SearchBoxHeight)
	// color.NRGBA{153, 217, 234, 192} // blue
	img.Fill(color.RGBA{128, 128, 128, 128}) // semi-transparent gray

	var (
		x float64 = cfg.ScreenSize.X - cfg.ChartTreeNodeWidth
		y float64 = 25
	)
	s := draws.NewSprite(img)
	s.Locate(x, y, draws.RightMiddle)
	asset.SearchBoxSprite = s
}

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
