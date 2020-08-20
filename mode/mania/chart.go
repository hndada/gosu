package mania

import (
	"github.com/hndada/gosu/mode"
	"github.com/hndada/rg-parser/osugame/osu"
	"log"
)

type Chart struct {
	mode.BaseChart
	Keys  int
	Notes []Note
}

// raw 차트에는 모드가 들어가면 안됨
// 모드마다 TransPoint(TimingPoint), Note건듦
func NewChartFromOsu(o *osu.Format) (*Chart, error) {
	var c Chart
	c.BaseChart = mode.NewBaseChartFromOsu(o)
	c.Keys = int(c.Parameter["Scale"])
	c.Notes = make([]Note, 0, len(o.HitObjects)*2)
	for _, ho := range o.HitObjects {
		ns, err := newNote(ho, c.Keys)
		if err != nil {
			log.Fatal("invalid hit object") // todo: normal log
		}
		c.Notes = append(c.Notes, ns...)
	}
	return &c, nil
}

func (c *Chart) ApplyMods(mods Mods) *Chart {
	var c2 Chart
	c2.BaseChart = c.BaseChart // todo: pointer?
	c2.Keys = c.Keys
	c2.Notes = make([]Note, len(c.Notes))
	for i, n := range c.Notes {
		n.Time = int64(float64(n.Time) / mods.TimeRate)
		n.Time2 = int64(float64(n.Time2) / mods.TimeRate)
		if mods.Mirror { // todo: scartch는 따로 분리? -> 까다로워질지도, 아니면 미러로 그냥 쇼부 봐
			n.Key = c.Keys - 1 - n.Key
		}
		c2.Notes[i] = n
	}
	return &c2
}
