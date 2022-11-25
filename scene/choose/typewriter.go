package choose

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/scene"
)

// Original source code is written by Hajimehoshi, the Author of Ebitengine.
// https://ebitengine.org/en/examples/typewriter.html
type TypeWriter struct {
	Text  string
	runes []rune
}

func (w *TypeWriter) Update() error {
	// Add runes that are input by the user by AppendInputChars.
	// Note that AppendInputChars result changes every frame, so you need to call this
	// every frame.
	w.runes = ebiten.AppendInputChars(w.runes[:0])
	w.Text += string(w.runes)

	// If the backspace key is pressed, remove one character.
	if w.isKeyPressed(input.KeyBackspace) {
		if len(w.Text) >= 1 {
			w.Text = w.Text[:len(w.Text)-1]
		}
	}
	if ebiten.IsKeyPressed(input.KeyEscape) {
		w.Text = ""
	}
	return nil
}
func (w TypeWriter) isKeyPressed(k input.Key) bool {
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
func (w TypeWriter) Draw(dst draws.Image) {
	text.Draw(dst.Image, w.Text, scene.Face16, 1200, 100, color.Black)
}
func (w *TypeWriter) Reset() {
	w.Text = ""
	w.runes = []rune{}
}
