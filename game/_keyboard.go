package game

import (
	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/types"
	"time"
)

// GetAsyncKeyState: most significant bit 가 set이면 해당 키 눌린 거

type KeyState int8

// const (
// 	KeyStateUnknown KeyState = iota - 1
// 	KeyStateUp
// 	KeyStateDown
// )

// func state(message types.Message) KeyState {
// 	switch message {
// 	case types.WM_KEYDOWN:
// 		return KeyStateDown
// 	case types.WM_KEYUP:
// 		return KeyStateUp
// 	default:
// 		return KeyStateUnknown
// 	}
// }

//
// type KeyboardEvent struct {
// 	Time    int64
// 	KeyCode types.VKCode
// 	State   KeyState
// }

type KeyboardEventChannel struct {
	start time.Time
	// q     chan chan types.KeyboardEvent
	Chan chan types.KeyboardEvent
	// events []KeyboardEvent
	// finish chan bool
}

func NewKeyboardEventChannel() (*KeyboardEventChannel, error) {
	// var c KeyboardEventChannel
	// c.q = make(chan chan types.KeyboardEvent, 5)
	// c.addChannel()
	c := &KeyboardEventChannel{
		Chan: make(chan types.KeyboardEvent, 100),
	}
	if err := keyboard.Install(nil, c.Chan); err != nil {
		return c, err
	}
	// c.events = make([]KeyboardEvent, 0, 100)
	// c.finish = make(chan bool, 1)
	return c, nil
}

// func (c *KeyboardEventChannel) addChannel() {
// 	newChan := make(chan types.KeyboardEvent, 100)
// 	if err := keyboard.Install(nil, newChan); err != nil {
// 		panic(err)
// 	}
// 	c.q <- newChan
// }
func (c *KeyboardEventChannel) Close() error {
	close(c.Chan)
	return keyboard.Uninstall()
}

func (c *KeyboardEventChannel) SetStartTime(t time.Time) { c.start = t }

func (c *KeyboardEventChannel) Time() time.Duration { return time.Since(c.start) }

// 별도 goroutine에서 돌게 해야할듯
// func (c *KeyboardEventChannel) Dequeue() []KeyboardEvent {
// 	ch := <-c.q
// 	_ = keyboard.Uninstall()
// 	c.addChannel()
// 	close(ch)
// 	events := make([]KeyboardEvent, 0, len(ch))
// 	for k := range ch {
// 		var e KeyboardEvent
// 		e.Time = time.Since(c.start).Milliseconds()
// 		e.KeyCode = k.VKCode
// 		e.State = state(k.Message)
// 		events = append(events, e)
// 	}
// 	return events
// }
