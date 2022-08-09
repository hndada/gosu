//go:build !windows

package main

func NewListener(keySettings []ebiten.Key) func(int64) KeysState {
	return func(now int64) KeysState {
		state := KeysState{now, make([]bool, len(keySettings))}
		for k, ek := range keySettings {
			state.Pressed[k] = ebiten.IsKeyPressed(ek)
		}
		return state
	}
}
