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

type ChartHeader struct {
	// ChartPath     string
	ChartID       int64 // 6byte: setID, 2byte: subID
	MusicName     string
	MusicUnicode  string
	Artist        string
	ArtistUnicode string
	MusicSource   string
	ChartName     string // diff name
	Producer      string
	HolderID      int64 // 0: gosu Chart Management

	AudioFilename   string
	PreviewTime     int64
	ImageFilename   string
	VideoFilename   string
	VideoTimeOffset int64
	Parameter       struct {
		CircleSize float64
		KeyCount   int
	}
	Level float64
}

// Chart should avoid redundant data as much as possible
type Chart struct {
	ChartHeader
	KeyCount    int
	ScratchMode int

	TransPoints
	Notes []Note
}

func NewChartHeaderFromOsu(o *osu.Format) ChartHeader {
	c := ChartHeader{
		// ChartPath:     path,
		MusicName:     o.Title,
		MusicUnicode:  o.TitleUnicode,
		Artist:        o.Artist,
		ArtistUnicode: o.ArtistUnicode,
		MusicSource:   o.Source,
		ChartName:     o.Version,
		Producer:      o.Creator, // field name may change

		AudioFilename: o.AudioFilename,
		PreviewTime:   int64(o.PreviewTime),
	}
	switch o.General.Mode {
	case ModeStandard, ModeCatch:
		c.Parameter.CircleSize = o.CircleSize
	case ModeMania:
		c.Parameter.KeyCount = int(o.CircleSize)
	}
	return c
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
		c.KeyCount = c.Parameter.KeyCount
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
