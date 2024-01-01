package ctrl

import "github.com/hndada/gosu/input"

// Required ticks to trigger pressing key:
// 0, long, short, short, short, ...
type UIKeyListener struct {
	keys []input.Key
	// Modifiers    []input.Key
	holdingKey   input.Key
	cooltimeTick int
	active       bool
}

func NewUIKeyListener(keys []input.Key) UIKeyListener {
	return UIKeyListener{
		keys:       keys,
		holdingKey: input.KeyNone,
	}
}

func (kl *UIKeyListener) Listen() input.Key {
	tk := input.KeyNone // triggered key
	if !input.IsKeyPressed(kl.holdingKey) {
		kl.holdingKey = input.KeyNone
		kl.cooltimeTick = 0
		kl.active = false
	}

	for _, key := range kl.keys {
		if input.IsKeyPressed(key) {
			kl.holdingKey = key
			if kl.cooltimeTick > 0 {
				kl.cooltimeTick--
			} else {
				tk = key
				if kl.active {
					kl.cooltimeTick = shortTicks
				} else {
					kl.cooltimeTick = longTicks
				}
				kl.active = true
			}
			break
		}
	}
	return tk
}
