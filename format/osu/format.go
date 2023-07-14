package osu

import (
	"bytes"
	"fmt"
	"image/color"
	"strings"
	"unicode"
)

const (
	ModeStandard = iota
	ModeTaiko
	ModeCatch
	ModeMania
)
const ModeOsu = ModeStandard

// Format is preferred name to Type
// because Type is a more general term.
type Format struct {
	FormatVersion int // delimiter:(space)
	General
	Editor
	Metadata
	Difficulty
	Events       []Event
	TimingPoints []TimingPoint
	Colours
	HitObjects []HitObject
}

type General struct { // delimiter:(space)
	AudioFilename            string
	AudioLeadIn              int
	AudioHash                string // deprecated
	PreviewTime              int
	Countdown                int
	SampleSet                string
	StackLeniency            float64
	Mode                     int
	LetterboxInBreaks        bool
	StoryFireInFront         bool // deprecated
	UseSkinSprites           bool
	AlwaysShowPlayfield      bool // deprecated
	OverlayPosition          string
	SkinPreference           string
	EpilepsyWarning          bool
	CountdownOffset          int
	SpecialStyle             bool
	WidescreenStoryboard     bool
	SamplesMatchPlaybackRate bool
}

type Editor struct { // delimiter:(space)
	Bookmarks       []int // delimiter,
	DistanceSpacing float64
	BeatDivisor     int
	GridSize        int
	TimelineZoom    float64
}

type Metadata struct { // delimiter:
	Title         string
	TitleUnicode  string
	Artist        string
	ArtistUnicode string
	Creator       string
	Version       string
	Source        string
	Tags          []string // delimiter(space)
	BeatmapID     int
	BeatmapSetID  int
}

type Difficulty struct { // delimiter:
	HPDrainRate       float64
	CircleSize        float64
	OverallDifficulty float64
	ApproachRate      float64
	SliderMultiplier  float64
	SliderTickRate    float64
}

type Colours struct {
	Combos              [8]color.RGBA
	SliderTrackOverride color.RGBA
	SliderBorder        color.RGBA
}

func NewFormat(dat []byte) (f *Format, err error) {
	dat = bytes.ReplaceAll(dat, []byte("\r\n"), []byte("\n"))

	f = &Format{
		General: General{
			PreviewTime:      -1,
			Countdown:        1,
			SampleSet:        "Normal",
			StackLeniency:    0.7,
			StoryFireInFront: true,
			OverlayPosition:  "NoChange",
		},
	}

	var section string
	for _, l := range bytes.Split(dat, []byte("\n")) {
		// TrimLeftFunc: prevent trimming delimiter
		l = bytes.TrimLeftFunc(l, unicode.IsSpace)
		line := string(l)

		if isPass(line) {
			continue
		}
		if isSection(line) {
			section = strings.Trim(line, "[]")
			continue
		}

		switch section {
		case "General":
			if err = f.setGeneralContent(line); err != nil {
				return
			}
		case "Editor":
			if err = f.setEditorContent(line); err != nil {
				return
			}
		case "Metadata":
			if err = f.setMetadataContent(line); err != nil {
				return
			}
		case "Difficulty":
			if err = f.setDifficultyContent(line); err != nil {
				return
			}
		case "Events":
			ev, err := newEvent(line)
			if err != nil {
				return f, fmt.Errorf("error at %s: %s", line, err)
			}
			f.Events = append(f.Events, ev)
		case "TimingPoints":
			tp, err := newTimingPoint(line)
			if err != nil {
				return f, fmt.Errorf("error at %s: %s", line, err)
			}
			f.TimingPoints = append(f.TimingPoints, tp)
		case "Colours":
			if err = f.setColoursContent(line); err != nil {
				return
			}
		case "HitObjects":
			ho, err := newHitObject(line)
			if err != nil {
				return f, fmt.Errorf("error at %s: %s", line, err)
			}
			f.HitObjects = append(f.HitObjects, ho)
		}
	}
	return f, nil
}

func isPass(line string) bool { return len(line) < 2 || line[:2] == "//" }
func isSection(line string) bool {
	return !isPass(line) && line[0] == '[' && line[len(line)-1] == ']'
}

func (f *Format) setGeneralContent(line string) error {
	k, v, err := keyValue(line, `: `)
	if err != nil {
		return err
	}

	switch k {
	case "AudioFilename":
		f.AudioFilename = v
	case "AudioLeadIn":
		if f.AudioLeadIn, err = parseInt(v); err != nil {
			return fmt.Errorf("error at %s: %s", line, err)
		}
	case "AudioHash":
		f.AudioHash = v
	case "PreviewTime":
		if f.PreviewTime, err = parseInt(v); err != nil {
			return fmt.Errorf("error at %s: %s", line, err)
		}
	case "Countdown":
		if f.Countdown, err = parseInt(v); err != nil {
			return fmt.Errorf("error at %s: %s", line, err)
		}
	case "SampleSet":
		f.SampleSet = v
	case "StackLeniency":
		if f.StackLeniency, err = parseFloat(v); err != nil {
			return fmt.Errorf("error at %s: %s", line, err)
		}
	case "Mode":
		if f.Mode, err = parseInt(v); err != nil {
			return fmt.Errorf("error at %s: %s", line, err)
		}
	case "LetterboxInBreaks":
		if f.LetterboxInBreaks, err = parseBool(v); err != nil {
			return fmt.Errorf("error at %s: %s", line, err)
		}
	case "StoryFireInFront":
		if f.StoryFireInFront, err = parseBool(v); err != nil {
			return fmt.Errorf("error at %s: %s", line, err)
		}
	case "UseSkinSprites":
		if f.UseSkinSprites, err = parseBool(v); err != nil {
			return fmt.Errorf("error at %s: %s", line, err)
		}
	case "AlwaysShowPlayfield":
		if f.AlwaysShowPlayfield, err = parseBool(v); err != nil {
			return fmt.Errorf("error at %s: %s", line, err)
		}
	case "OverlayPosition":
		f.OverlayPosition = v
	case "SkinPreference":
		f.SkinPreference = v
	case "EpilepsyWarning":
		if f.EpilepsyWarning, err = parseBool(v); err != nil {
			return fmt.Errorf("error at %s: %s", line, err)
		}
	case "CountdownOffset":
		if f.CountdownOffset, err = parseInt(v); err != nil {
			return fmt.Errorf("error at %s: %s", line, err)
		}
	case "SpecialStyle":
		if f.SpecialStyle, err = parseBool(v); err != nil {
			return fmt.Errorf("error at %s: %s", line, err)
		}
	case "WidescreenStoryboard":
		if f.WidescreenStoryboard, err = parseBool(v); err != nil {
			return fmt.Errorf("error at %s: %s", line, err)
		}
	case "SamplesMatchPlaybackRate":
		if f.SamplesMatchPlaybackRate, err = parseBool(v); err != nil {
			return fmt.Errorf("error at %s: %s", line, err)
		}
	}
	return nil
}

func (f *Format) setEditorContent(line string) error {
	k, v, err := keyValue(line, `: `)
	if err != nil {
		return err
	}
	// Number-only sections may be trimmed both space.
	v = strings.TrimSpace(v)

	switch k {
	case "Bookmarks":
		for _, s := range strings.Split(v, ",") {
			var bookmark int
			if bookmark, err = parseInt(s); err != nil {
				return fmt.Errorf("error at %s: %s", line, err)
			}
			f.Bookmarks = append(f.Bookmarks, bookmark)
		}
	case "DistanceSpacing":
		if f.DistanceSpacing, err = parseFloat(v); err != nil {
			return fmt.Errorf("error at %s: %s", line, err)
		}
	case "BeatDivisor":
		if f.BeatDivisor, err = parseInt(v); err != nil {
			return fmt.Errorf("error at %s: %s", line, err)
		}
	case "GridSize":
		if f.GridSize, err = parseInt(v); err != nil {
			return fmt.Errorf("error at %s: %s", line, err)
		}
	case "TimelineZoom":
		if f.TimelineZoom, err = parseFloat(v); err != nil {
			return fmt.Errorf("error at %s: %s", line, err)
		}
	}
	return nil
}

func (f *Format) setMetadataContent(line string) error {
	k, v, err := keyValue(line, `: `)
	if err != nil {
		return err
	}

	switch k {
	case "Title":
		f.Title = v
	case "TitleUnicode":
		f.TitleUnicode = v
	case "Artist":
		f.Artist = v
	case "ArtistUnicode":
		f.ArtistUnicode = v
	case "Creator":
		f.Creator = v
	case "Version":
		f.Version = v
	case "Source":
		f.Source = v
	case "Tags":
		f.Tags = strings.Split(v, " ")
	case "BeatmapID":
		if f.BeatmapID, err = parseInt(v); err != nil {
			return fmt.Errorf("error at %s: %s", line, err)
		}
	case "BeatmapSetID":
		if f.BeatmapSetID, err = parseInt(v); err != nil {
			return fmt.Errorf("error at %s: %s", line, err)
		}
	}
	return nil
}

func (f *Format) setDifficultyContent(line string) error {
	k, v, err := keyValue(line, `: `)
	if err != nil {
		return err
	}
	// Number-only sections may be trimmed both space.
	v = strings.TrimSpace(v)

	switch k {
	case "HPDrainRate":
		if f.HPDrainRate, err = parseFloat(v); err != nil {
			return fmt.Errorf("error at %s: %s", line, err)
		}
	case "CircleSize":
		if f.CircleSize, err = parseFloat(v); err != nil {
			return fmt.Errorf("error at %s: %s", line, err)
		}
	case "OverallDifficulty":
		if f.OverallDifficulty, err = parseFloat(v); err != nil {
			return fmt.Errorf("error at %s: %s", line, err)
		}
	case "ApproachRate":
		if f.ApproachRate, err = parseFloat(v); err != nil {
			return fmt.Errorf("error at %s: %s", line, err)
		}
	case "SliderMultiplier":
		if f.SliderMultiplier, err = parseFloat(v); err != nil {
			return fmt.Errorf("error at %s: %s", line, err)
		}
	case "SliderTickRate":
		if f.SliderTickRate, err = parseFloat(v); err != nil {
			return fmt.Errorf("error at %s: %s", line, err)
		}
	}
	return nil
}

func (f *Format) setColoursContent(line string) error {
	k, v, err := keyValue(line, `: `)
	if err != nil {
		return err
	}
	// Number-only sections may be trimmed both space.
	v = strings.TrimSpace(v)

	rgb := newRGB(v)
	if strings.HasPrefix(k, "Combo") {
		i, err := parseInt(k[5:])
		if err != nil {
			return fmt.Errorf("error at %s: %s", line, err)
		}
		f.Combos[i-1] = rgb
		return nil
	}

	switch k {
	case "SliderTrackOverride":
		f.SliderTrackOverride = rgb
	case "SliderBorder":
		f.SliderBorder = rgb
	}
	return nil
}

func newRGB(chunks string) color.RGBA {
	var rgb color.RGBA
	for i, chunk := range strings.Split(chunks, `,`) {
		v, _ := parseInt(chunk)
		switch i {
		case 0:
			rgb.R = uint8(v)
		case 1:
			rgb.G = uint8(v)
		case 2:
			rgb.B = uint8(v)
		}
	}
	rgb.A = 255
	return rgb
}
