package ctrl

import (
	"fmt"

	draws "github.com/hndada/gosu/draws2"
)

// type BaseBox[T any] struct {
// 	Parent  *BaseBox[T]
// 	Befores []*BaseBox[T]
// 	Afters  []*BaseBox[T]
// }

// type Box BaseBox[struct{}]

type Widget struct {
	Parent *Widget
	Before []*Widget
	After  []*Widget

	draws.Box
}

func NewButton() Widget {
	return Widget{
		Box: draws.Box{},
	}
}

type EventHandler func() any

func (b *Button) AddEventHandler(e Event, eh EventHandler) {

}

func (b *Button) Update() any {
	switch r := c.Button.Update().(type) {
	case *draws.Box:
		fmt.Println(r.Afters)
	}
	return nil
}

type ChartItem struct {
	ChartHeader struct{}
	Button
}

func NewChartItem() *ChartItem {
	return &ChartItem{
		Button: Button{},
	}
}

func (c *ChartItem) Update() {

}
