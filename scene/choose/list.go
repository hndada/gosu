package choose

import (
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/scene"
)

// Group1, Group2, Sort, Filter int

const (
	listDepthSearch = -1
	listDepthMusic  = iota
	listDepthChart
	// FocusKeySettings
)

// var NewSliceIndexKeyHandler func(index *int, len int) ctrl.KeyHandler
func (s Scene) NewListMoveKeyHandler(handler ctrl.Handler) ctrl.KeyHandler {
	return ctrl.KeyHandler{
		Handler:  handler,
		Modifier: input.KeyNone,
		Keys:     scene.UpDownKeys,
		Sounds:   [2]audios.SoundPlayer{s.SwipeSoundPod, s.SwipeSoundPod},
		Volume:   &s.SoundVolume,
	}
}

//	if len(children) != 0 {
//		handler := NewSliceIndexKeyHandler(&list.Cursor, len(list.Children))
//		list.keyHandleCursor = handler.Handle
//	}
