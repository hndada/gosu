package mania

import "github.com/hndada/gosu/engine/kb"

type keyEvent struct {
	kb.KeyEvent
	Key int
}
