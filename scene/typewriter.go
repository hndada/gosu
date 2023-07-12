package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hndada/gosu/input"
)

type TypeWriter struct {
	Text  string
	runes []rune
}

// Original source code is written by Hajimehoshi, the Author of Ebitengine.
// https://ebitengine.org/en/examples/typewriter.html
func IsRepeatKeyPressed(k input.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(k)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

func (w *TypeWriter) Update() error {
	w.runes = ebiten.AppendInputChars(w.runes[:0])
	w.Text += string(w.runes)
	if IsRepeatKeyPressed(input.KeyBackspace) {
		if len(w.Text) >= 1 {
			w.Text = w.Text[:len(w.Text)-1]
		}
	}
	if ebiten.IsKeyPressed(input.KeyEscape) {
		w.Text = ""
	}
	return nil
}

func (w *TypeWriter) Reset() {
	w.Text = ""
	w.runes = []rune{}
}
