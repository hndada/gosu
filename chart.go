package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hndada/gosu/parse/osu"
)

const (
	ModeStandard = iota
	ModeTaiko
	ModeCatch
	ModeMania
)

const DefaultMode = ModeStandard

// ChartHeader contains non-play information.
// Chaning ChartHeader's data will not affect integrity of the chart.
// Mode-specific fields are moved to each Chart struct
type ChartHeader struct {
	ChartID       int64 // 6byte: setID, 2byte: subID
	MusicName     string
	MusicUnicode  string
	Artist        string
	ArtistUnicode string
	MusicSource   string
	ChartName     string // diff name
	Producer      string // Name of field may change
	HolderID      int64  // 0: gosu Chart Management

	AudioFilename   string
	PreviewTime     int64
	ImageFilename   string
	VideoFilename   string
	VideoTimeOffset int64
}

// Chart should avoid redundant data as much as possible
type Chart struct {
	ChartHeader
	TransPoints
	KeyCount int
	Notes    []Note
	Level    float64
	// ScratchMode int
}

func NewChartHeaderFromOsu(o *osu.Format) ChartHeader {
	return ChartHeader{
		MusicName:     o.Title,
		MusicUnicode:  o.TitleUnicode,
		Artist:        o.Artist,
		ArtistUnicode: o.ArtistUnicode,
		MusicSource:   o.Source,
		ChartName:     o.Version,
		Producer:      o.Creator,

		AudioFilename: o.AudioFilename,
		PreviewTime:   int64(o.PreviewTime),
	}
}
func NewChart(path string) (*Chart, error) {
	var c Chart
	switch strings.ToLower(filepath.Ext(path)) {
	case ".osu":
		dat, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		o, err := osu.Parse(dat)
		if err != nil {
			return nil, err
		}
		c.ChartHeader = NewChartHeaderFromOsu(o)
		c.TransPoints = NewTransPointsFromOsu(o)
		c.KeyCount = int(o.CircleSize)
		c.Notes = make([]Note, 0, len(o.HitObjects)*2)
		for _, ho := range o.HitObjects {
			c.Notes = append(c.Notes, NewNoteFromOsu(ho, c.KeyCount)...)
		}
		sort.Slice(c.Notes, func(i, j int) bool {
			if c.Notes[i].Time == c.Notes[j].Time {
				return c.Notes[i].Key < c.Notes[j].Key
			}
			return c.Notes[i].Time < c.Notes[j].Time
		})
		return &c, nil
	}
	return nil, errors.New("not supported")
}
