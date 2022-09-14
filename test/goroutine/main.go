package main

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu"
	"github.com/hndada/gosu/input"
)

// input.Event
type InputEvent struct {
	Time    int64
	Pressed []bool
}
type KeyLogger struct {
	TPS          int
	Tick         int
	FetchPressed func() []bool
	Events       []InputEvent
	Cursor       int
}

func NewKeyLogger(keySettings []input.Key, time int64) (kl KeyLogger) {
	kl.TPS = 1000
	kl.SetTime(time)
	kl.FetchPressed = input.NewListener(keySettings)
	return
}
func (kl *KeyLogger) SetTime(time int64) {
	kl.Tick = int(float64(time) / 1000 * float64(kl.TPS))
}
func (kl KeyLogger) Time() int64 {
	return int64(float64(kl.Tick) / float64(kl.TPS) * 1000)
}
func (kl *KeyLogger) Update() {
	pressed := kl.FetchPressed()
	event := InputEvent{
		Time:    kl.Time(),
		Pressed: pressed,
	}
	if len(kl.Events) == 0 {
		kl.Events = append(kl.Events, event)
	} else {
		lastPressed := kl.Events[len(kl.Events)-1].Pressed
		for i, p := range pressed {
			lp := lastPressed[i]
			if p != lp {
				kl.Events = append(kl.Events, event)
				break
			}
		}
	}
	kl.Tick++
}

const (
	screenSizeX = 320
	screenSizeY = 240
)

type Game struct {
	TPS  int
	Tick int
	KeyLogger
	Events []InputEvent // For printing.
}

func (g *Game) Update() (err error) {
	// g.KeyLogger.Update()
	es := g.KeyLogger.Events[g.KeyLogger.Cursor:]
	g.Events = append(g.Events, es...)
	g.KeyLogger.Cursor += len(es)
	g.Tick++
	return
}
func (g *Game) Draw(screen *ebiten.Image) {
	time := gosu.TickToTime(g.Tick)
	var s string
	for _, e := range g.Events {
		s += fmt.Sprintf("%v\n", e)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"FPS: %.2f\nTPS: %.2f\nTime: %.3fs\n\n"+
			"Key inputs:\n%s",
		ebiten.ActualFPS(), ebiten.ActualTPS(), float64(time)/float64(g.TPS), //time%int64(g.TPS),
		s))
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenSizeX, screenSizeY
}
func (g Game) TimeToTick(time int64) int {
	return int(float64(time) / 1000 * float64(g.TPS))
}
func (g Game) TickToTime() int64 {
	return int64(float64(g.Tick) / float64(g.TPS) * 1000)
}
func main() {
	g := new(Game)
	g.TPS = 5
	ebiten.SetTPS(g.TPS)
	g.Tick = g.TimeToTick(-1800)
	// fmt.Println(g.Tick, g.TickToTime())
	keySettings := []input.Key{input.KeyD, input.KeyF, input.KeyJ, input.KeyK}
	g.KeyLogger = NewKeyLogger(keySettings, g.TickToTime())
	go func() {
		for {
			start := time.Now()
			g.KeyLogger.Update()
			// fmt.Println(start, time.Since(start))
			// fmt.Println(50*time.Millisecond - time.Since(start))
			// fmt.Println(g.KeyLogger.TPS, g.KeyLogger.Tick)
			time.Sleep(1000*time.Microsecond - time.Since(start))
		}
	}()
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
