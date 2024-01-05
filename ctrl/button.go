package ctrl

import (
	"fmt"

	draws "github.com/hndada/gosu/draws2"
)

type Box draws.Box

// func NewButton() Box {}

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
