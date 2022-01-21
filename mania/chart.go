package mania

import (
	"path/filepath"
	"strings"

	"github.com/hndada/gosu/common"
	"github.com/hndada/rg-parser/osugame/osu"
)

type Chart struct {
	*common.ChartHeader
	common.TimingPoints
	KeyCount    int
	ScratchMode int
	Notes       []Note
	TimeStamps  []common.TimeStamp
}

// A raw chart data should be not 'Mods-affected': Mods modify TransPoint(TimingPoint) and Note.
func NewChart(path string) (*Chart, error) {
	var c Chart
	switch strings.ToLower(filepath.Ext(path)) {
	case ".osu":
		o, err := osu.Parse(path)
		if err != nil {
			panic(err)
		}
		c.ChartHeader = common.NewChartHeaderFromOsu(o, path)
		c.TimingPoints = common.NewTimingPointsFromOsu(o)
		c.KeyCount = int(c.Parameter["KeyCount"])
		err = c.loadNotesFromOsu(o)
		if err != nil {
			panic(err)
		}
		c.TimingPoints.SpeedFactors = append(
			[]common.SpeedFactorPoint{common.DefaultSpeedFactor},
			c.TimingPoints.SpeedFactors...)

		c.setStamps()
		c.setNotePosition()
		c.CalcDifficulty()
		// fmt.Println(c.MusicName, c.ChartName)
		return &c, nil
	default:
		panic("not reach")
	}
}

func (c *Chart) ApplyMods(mods Mods) *Chart {
	var c2 Chart
	c2.ChartHeader = c.ChartHeader // TODO: value -> pointer?
	c2.KeyCount = c.KeyCount
	c2.Notes = make([]Note, len(c.Notes))
	for i, n := range c.Notes {
		n.Time = int64(float64(n.Time) / mods.TimeRate)
		n.Time2 = int64(float64(n.Time2) / mods.TimeRate)
		if mods.Mirror { // TODO: separate scartch lane?
			n.key = c.KeyCount - 1 - n.key
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
	stamps := make([]common.TimeStamp, len(fs))
	for i, f := range fs {
		s := common.TimeStamp{
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
func (c Chart) TimeStampFinder() func(time int64) common.TimeStamp {
	var cursor int
	return func(time int64) common.TimeStamp {
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
