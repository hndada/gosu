package choose

import (
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/input"
)

// Todo: OneAtATime option
func (s *Scene) setSubModeKeyHandler(cfg *Config, asset *Asset) {
	// s.SubMode = ctrl.KeyHandler{
	// 	Handler: ctrl.IntHandler{
	// 		Value: &s.subMode,
	// 		Min:   4,
	// 		Max:   9,
	// 		Loop:  true,
	// 	},
	// 	Modifiers: []input.Key{},
	// 	Keys:      [2]input.Key{input.KeyF2, input.KeyF3},
	// 	Sounds:    [2]audios.Sounder{scene.UserSkin.Swipe, scene.UserSkin.Swipe},
	// 	Volume:    &mode.S.SoundVolume,
	// }
}

// Todo: ctrl.UpDownKeys -> input.UpDownKeys?
func (s *Scene) setCursorKeyHandler(slice []any) {
	s.cursorKeyHandlers = make([]ctrl.KeyHandler, len(slice))
	for i := range s.cursorKeyHandlers {
		s.cursorKeyHandlers[i] = ctrl.KeyHandler{
			Handler: &ctrl.IntHandler{
				Value: &cfg.cursor[i],
				Min:   0,
				Max:   len(ts) - 1,
				Loop:  false,
			},
			Modifier: input.KeyNone,
			Keys:     ctrl.UpDownKeys,
			Sounds:   asset.SwipeSounds,
			Volume:   &cfg.SoundVolume,
		}
	}
}
