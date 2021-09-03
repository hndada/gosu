package mania

import (
	"path/filepath"
	"strings"

	"github.com/hndada/gosu/game"
	"github.com/hndada/rg-parser/osugame/osu"
)

type Chart struct {
	*game.ChartHeader
	game.TimingPoints
	KeyCount    int
	ScratchMode int
	Notes       []Note
	TimeStamps  []game.TimeStamp
}

// raw 차트에는 Mods가 들어가면 안됨
// Mods마다 TransPoint(TimingPoint), Note건드림
func NewChart(path string) (*Chart, error) {
	var c Chart
	switch strings.ToLower(filepath.Ext(path)) {
	case ".osu":
		o, err := osu.Parse(path)
		if err != nil {
			panic(err)
		}
		c.ChartHeader = game.NewChartHeaderFromOsu(o, path)
		c.TimingPoints = game.NewTimingPointsFromOsu(o)
		c.KeyCount = int(c.Parameter["KeyCount"])
		err = c.loadNotesFromOsu(o)
		if err != nil {
			panic(err)
		}
		c.TimingPoints.SpeedFactors = append(
			[]game.SpeedFactorPoint{game.DefaultSpeedFactor},
			c.TimingPoints.SpeedFactors...)

		c.setStamps()
		c.setNotePosition()
		c.CalcDifficulty()
		return &c, nil
	default:
		panic("not reach")
	}
}

func (c *Chart) ApplyMods(mods Mods) *Chart {
	var c2 Chart
	c2.ChartHeader = c.ChartHeader // todo: pointer?
	c2.KeyCount = c.KeyCount
	c2.ScratchMode = mods.ScratchMode // temp
	c2.Notes = make([]Note, len(c.Notes))
	for i, n := range c.Notes {
		n.Time = int64(float64(n.Time) / mods.TimeRate)
		n.Time2 = int64(float64(n.Time2) / mods.TimeRate)
		if mods.Mirror { // todo: scartch는 따로 분리? -> 까다로워질지도, 아니면 미러로 그냥 쇼부 봐
			n.Key = c.KeyCount - 1 - n.Key
		}
		c2.Notes[i] = n
	}
	c2.CalcDifficulty()
	return &c2
}

func (c *Chart) setStamps() {
	const maxInt64 = 9223372036854775807
	var position float64
	fs := c.TimingPoints.SpeedFactors
	stamps := make([]game.TimeStamp, len(fs))
	for i, f := range fs {
		s := game.TimeStamp{
			Time:     f.Time,
			Position: position,
			Factor:   f.Factor,
		}
		if i < len(fs)-1 {
			nextTime := fs[i+1].Time
			s.NextTime = nextTime
			position += float64(nextTime-f.Time) * f.Factor
		} else {
			s.NextTime = maxInt64
		}
		stamps[i] = s
	}
	c.TimeStamps = stamps
}
func (c Chart) EndTime() int64 {
	return c.Notes[len(c.Notes)-1].Time
}

// return value: timeStamp()
func (c Chart) TimeStampFinder() func(time int64) game.TimeStamp {
	var cursor int
	return func(time int64) game.TimeStamp {
		for si := range c.TimeStamps[cursor:] {
			if time < c.TimeStamps[cursor+si].NextTime {
				ts := c.TimeStamps[cursor+si]
				cursor += si
				return ts
			}
		}
		panic("not reach")
	}
}

func (c Chart) LNCount() int {
	var num int
	for _, n := range c.Notes {
		if n.Type == TypeLNHead {
			num++
		}
	}
	return num
}
