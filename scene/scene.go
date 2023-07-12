package scene

import (
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/draws"
)

type Scene interface {
	Update() any
	Draw(screen draws.Image)
}

// Todo: struct or pointer to struct?
var TheBaseScene *BaseScene

type BaseScene struct {
	MusicVolumeKeyHandler          ctrl.KeyHandler
	SoundVolumeKeyHandler          ctrl.KeyHandler
	BackgroundBrightnessKeyHandler ctrl.KeyHandler
	OffsetKeyHandler               ctrl.KeyHandler
	DebugPrintKeyHandler           ctrl.KeyHandler
	SpeedScaleKeyHandlers          []ctrl.KeyHandler

	BackgroundSprite draws.Sprite
}

func NewBaseScene(cfg *Config, asset *Asset) *BaseScene {
	s := &BaseScene{}
	s.setMusicVolumeKeyHandler(cfg, asset)
	s.setSoundVolumeKeyHandler(cfg, asset)
	s.setBackgroundBrightnessKeyHandler(cfg, asset)
	s.setOffsetKeyHandler(cfg, asset)
	s.setDebugPrintKeyHandler(cfg, asset)
	s.setSpeedScaleKeyHandlers(cfg, asset)
	return s
}

// type FromChooseToPlay struct {
// 	cfg     *Config
// 	asset   *Asset
// 	fsys    fs.FS
// 	name    string
// 	_mode   int
// 	subMode int
// 	mods    any
// 	rf      *osr.Format
// }
