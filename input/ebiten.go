package input

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type EbitenInput struct {
	start       time.Time
	keyEvents   []KeyEvent
	lastPressed [256]bool
	closed      bool
}

// Todo: rename either one of 'start's
func (ei *EbitenInput) Listen(start time.Time) {
	ei.start = start
	ei.keyEvents = make([]KeyEvent, 0, 10)
	go ei.scan()
}
func (ei *EbitenInput) Flush() []KeyEvent {
	es := ei.keyEvents
	ei.keyEvents = make([]KeyEvent, 0, 10)
	return es
}
func (ei *EbitenInput) Close() { ei.closed = true }

// Todo: has not tested yet.
func (ei *EbitenInput) scan() {
	const d = 1 * time.Millisecond // Todo: Test from 9 to 1 gradually
	for {
		enter := time.Now()
		t := time.Since(ei.start).Milliseconds()
		for i := 0; i < int(ebiten.KeyMax); i++ {
			switch {
			case !ei.lastPressed[i] && ei.isKeyPressed(i):
				// fmt.Printf("%s pressed at %v ms\n", code, t)
				e := KeyEvent{
					Time:    t,
					KeyCode: ebitenKeyToCode(i),
					Pressed: true,
				}
				ei.keyEvents = append(ei.keyEvents, e)
				ei.lastPressed[i] = true
			case ei.lastPressed[i] && !ei.isKeyPressed(i):
				// fmt.Printf("%s released at %v ms\n", code, t)
				e := KeyEvent{
					Time:    t,
					KeyCode: ebitenKeyToCode(i),
					Pressed: false,
				}
				ei.keyEvents = append(ei.keyEvents, e)
				ei.lastPressed[i] = false
			}
		}
		remained := d - time.Since(enter) // Todo: should subtract -1?
		time.Sleep(remained)              // prevents 100% CPU usage
		if ei.closed {
			return
		}
	}
}
func (ei EbitenInput) isKeyPressed(i int) bool {
	return ebiten.IsKeyPressed(ebiten.Key(i))
}
